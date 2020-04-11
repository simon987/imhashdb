package imhashdb

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/mailru/easyjson"
	"github.com/simon987/fastimagehash-go"
	"go.uber.org/zap"
)

const MaxDistance = 30
const MaxLimit = 1000

type Entry struct {
	AHash  *fastimagehash.Hash
	DHash  *fastimagehash.Hash
	MHash  *fastimagehash.Hash
	PHash  *fastimagehash.Hash
	WHash  *fastimagehash.Hash
	Size   int
	Sha1   [sha1.Size]byte
	Md5    [md5.Size]byte
	Sha256 [sha256.Size]byte
	Crc32  uint32
	Meta   []Meta
	Url    string
}

func Store(entry *Entry) {
	row := Pgdb.QueryRow(
		`INSERT INTO image (size, sha1, md5, sha256, crc32) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		entry.Size, entry.Sha1[:], entry.Md5[:], entry.Sha256[:], entry.Crc32,
	)

	var id int64
	imageExists := false
	err := row.Scan(&id)
	if err != nil {
		imageExists = true
		row = Pgdb.QueryRow(`SELECT id FROM image WHERE sha1=$1`, entry.Sha1[:])
		err := row.Scan(&id)

		if err != nil {
			Logger.Error("FIXME: Could not insert image", zap.Error(err))
			return
		}
	}

	if !imageExists {
		_, err = Pgdb.Exec("INSERT INTO hash_ahash VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.AHash.Bytes)
		if err != nil {
			Logger.Error("Could not insert ahash", zap.Error(err))
		}
		_, err = Pgdb.Exec("INSERT INTO hash_dhash VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.DHash.Bytes)
		if err != nil {
			Logger.Error("Could not insert dhash", zap.Error(err))
		}
		_, err = Pgdb.Exec("INSERT INTO hash_mhash VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.MHash.Bytes)
		if err != nil {
			Logger.Error("Could not insert mhash", zap.Error(err))
		}
		_, err = Pgdb.Exec("INSERT INTO hash_phash VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.PHash.Bytes)
		if err != nil {
			Logger.Error("Could not insert phash", zap.Error(err))
		}
		_, err = Pgdb.Exec("INSERT INTO hash_whash VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.WHash.Bytes)
		if err != nil {
			Logger.Error("Could not insert whash", zap.Error(err))
		}
	}

	for _, meta := range entry.Meta {
		_, err = Pgdb.Exec(
			"INSERT INTO image_meta VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			meta.Id, meta.RetrievedAt, meta.Meta,
		)
		if err != nil {
			Logger.Error("Could not insert meta", zap.Error(err))
			return
		}
		_, err = Pgdb.Exec(
			"INSERT INTO image_has_meta VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			id, entry.Url, meta.Id,
		)
		if err != nil {
			Logger.Error("Could not insert ihm", zap.Error(err))
			return
		}
	}
}

func isHashValid(hash []byte, hashType HashType) bool {
	switch hashType {
	case AHash12:
		if len(hash) != 18 {
			return false
		}
	case DHash12:
		if len(hash) != 18 {
			return false
		}
	case MHash12:
		if len(hash) != 18 {
			return false
		}
	case PHash12:
		if len(hash) != 18 {
			return false
		}
	case WHash8Haar:
		if len(hash) != 8 {
			return false
		}
	default:
		return false
	}
	return true
}

func FindImagesByHash(ctx context.Context, hash []byte, hashType HashType, distance, limit, offset uint) ([]byte, error) {

	if !isHashValid(hash, hashType) {
		return nil, errors.New("invalid hash")
	}

	if distance > MaxDistance {
		return nil, errors.New("Invalid distance")
	}

	if limit > MaxLimit {
		return nil, errors.New("Invalid distance")
	}

	tx, err := Pgdb.BeginEx(ctx, &pgx.TxOptions{IsoLevel: pgx.ReadUncommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	var sql string
	switch hashType {
	case AHash12:
		sql = `SELECT image.* FROM image INNER JOIN hash_ahash h on image.id = h.image_id 
				WHERE hash_is_within_distance18(h.hash, $1, $2) ORDER BY image.id LIMIT $3 OFFSET $4`
	case DHash12:
		sql = `SELECT image.* FROM image INNER JOIN hash_dhash h on image.id = h.image_id 
				WHERE hash_is_within_distance18(h.hash, $1, $2) ORDER BY image.id LIMIT $3 OFFSET $4`
	case MHash12:
		sql = `SELECT image.* FROM image INNER JOIN hash_mhash h on image.id = h.image_id 
				WHERE hash_is_within_distance18(h.hash, $1, $2) ORDER BY image.id LIMIT $3 OFFSET $4`
	case PHash12:
		sql = `SELECT image.* FROM image INNER JOIN hash_phash h on image.id = h.image_id 
				WHERE hash_is_within_distance18(h.hash, $1, $2) ORDER BY image.id LIMIT $3 OFFSET $4`
	case WHash8Haar:
		sql = `SELECT image.* FROM image INNER JOIN hash_whash h on image.id = h.image_id 
				WHERE hash_is_within_distance8(h.hash, $1, $2) ORDER BY image.id LIMIT $3 OFFSET $4`
	}

	rows, err := tx.Query(sql, hash, distance, limit, offset)
	if err != nil {
		return nil, err
	}

	var images []*Image
	for rows.Next() {
		var im Image
		err := rows.Scan(&im.id, &im.Size, &im.Sha1, &im.Md5, &im.Sha256, &im.Crc32)
		if err != nil {
			Logger.Error("Error while fetching db image", zap.String("err", err.Error()))
			return nil, err
		}

		images = append(images, &im)
	}

	batch := tx.BeginBatch()
	defer batch.Close()
	for _, im := range images {
		batch.Queue(
			`SELECT ihm.url, meta.id, meta.retrieved_at, meta.meta FROM image_has_meta ihm
			INNER JOIN image_meta meta on ihm.image_meta_id = meta.id
			WHERE image_id=$1`,
			[]interface{}{im.id},
			[]pgtype.OID{pgtype.Int4OID},
			nil,
		)
	}

	err = batch.Send(ctx, nil)
	if err != nil {
		Logger.Error("Error while fetching db meta", zap.String("err", err.Error()))
		return nil, err
	}

	for _, im := range images {
		rows, err := batch.QueryResults()
		if err != nil {
			Logger.Error("Error while fetching db meta", zap.String("err", err.Error()))
			return nil, err
		}

		for rows.Next() {
			var ihm ImageHasMeta
			err := rows.Scan(&ihm.Url, &ihm.Meta.Id, &ihm.Meta.RetrievedAt, &ihm.Meta.Meta)
			if err != nil {
				return nil, err
			}
			im.Meta = append(im.Meta, ihm)
		}
	}

	b, _ := easyjson.Marshal(ImageList{Images: images})
	return b, nil
}

func DbInit(pool *pgx.ConnPool) {

	sql := `
CREATE TABLE IF NOT EXISTS image (
	id BIGSERIAL PRIMARY KEY,
	size INT,
	sha1 bytea,
	md5 bytea,
	sha256 bytea,
	crc32 bigint
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_image_sha1 ON image(sha1);
CREATE INDEX IF NOT EXISTS idx_image_md5 ON image(md5);
CREATE INDEX IF NOT EXISTS idx_image_sha256 ON image(sha256);
CREATE INDEX IF NOT EXISTS idx_image_crc32 ON image(crc32);

CREATE TABLE IF NOT EXISTS image_meta (
	id TEXT UNIQUE,
	retrieved_at bigint,
	meta bytea
);

CREATE TABLE IF NOT EXISTS image_has_meta (
	image_id bigint REFERENCES image(id),
	url TEXT,
	image_meta_id text REFERENCES image_meta(id),
	UNIQUE(image_id, image_meta_id)
);

CREATE TABLE IF NOT EXISTS hash_phash (
	image_id BIGINT REFERENCES image(id) UNIQUE,
    hash bytea
);

CREATE TABLE IF NOT EXISTS hash_ahash (
	image_id BIGINT REFERENCES image(id) UNIQUE,
    hash bytea
);

CREATE TABLE IF NOT EXISTS hash_dhash (
	image_id BIGINT REFERENCES image(id) UNIQUE,
    hash bytea
);

CREATE TABLE IF NOT EXISTS hash_mhash (
	image_id BIGINT REFERENCES image(id) UNIQUE,
    hash bytea
);

CREATE TABLE IF NOT EXISTS hash_whash (
	image_id BIGINT REFERENCES image(id) UNIQUE,
    hash bytea
);
`

	_, err := pool.Exec(sql)
	if err != nil {
		Logger.Fatal("Could not initialize database", zap.String("err", err.Error()))
	}
}

func DbConnect(host string, port int, user, password, database string) *pgx.ConnPool {
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     host,
			User:     user,
			Port:     uint16(port),
			Password: password,
			Database: database,
		},
		MaxConnections: 10,
	}

	var err error
	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		panic(err)
	}

	return pool
}
