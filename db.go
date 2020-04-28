package imhashdb

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/mailru/easyjson"
	"github.com/valyala/gozstd"
	"go.uber.org/zap"
)

const MaxDistance = 100
const MaxLimit = 1000

type Entry struct {
	H      *Hashes
	Size   int
	Sha1   [sha1.Size]byte
	Md5    [md5.Size]byte
	Sha256 [sha256.Size]byte
	Crc32  uint32
	Meta   []Meta
	Url    string
}

type MatchTrigger struct {
	HashType    HashType
	MinDistance int
	Id          int
}

var MatchTriggers = []MatchTrigger{
	{
		HashType:    DHash16,
		MinDistance: 25,
		Id:          1,
	},
	{
		HashType:    PHash16,
		MinDistance: 25,
		Id:          2,
	},
	{
		HashType:    WHash16Haar,
		MinDistance: 6,
		Id:          3,
	},
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
		_, err = Pgdb.Exec("INSERT INTO hash_dhash8 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.DHash8.Bytes)
		if err != nil {
			panic(err)
		}
		_, _ = Pgdb.Exec("INSERT INTO hash_dhash16 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.DHash16.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_dhash32 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.DHash32.Bytes)

		_, _ = Pgdb.Exec("INSERT INTO hash_mhash8 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.MHash8.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_mhash16 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.MHash16.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_mhash32 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.MHash32.Bytes)

		_, _ = Pgdb.Exec("INSERT INTO hash_phash8 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.PHash8.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_phash16 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.PHash16.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_phash32 VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.PHash32.Bytes)

		_, _ = Pgdb.Exec("INSERT INTO hash_whash8haar VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.WHash8.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_whash16haar VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.WHash16.Bytes)
		_, _ = Pgdb.Exec("INSERT INTO hash_whash32haar VALUES ($1, $2) ON CONFLICT DO NOTHING", id, entry.H.WHash32.Bytes)
	}

	var buf []byte
	for _, meta := range entry.Meta {
		compressedMeta := gozstd.CompressDict(buf[:0], meta.Meta, CDict)

		_, err = Pgdb.Exec(
			"INSERT INTO image_meta (id, retrieved_at, meta) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			meta.Id, meta.RetrievedAt, compressedMeta,
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
	case DHash8:
		fallthrough
	case MHash8:
		fallthrough
	case PHash8:
		fallthrough
	case WHash8Haar:
		return len(hash) == 8

	case DHash16:
		fallthrough
	case MHash16:
		fallthrough
	case PHash16:
		fallthrough
	case WHash16Haar:
		return len(hash) == 32

	case DHash32:
		fallthrough
	case MHash32:
		fallthrough
	case PHash32:
		fallthrough
	case WHash32Haar:
		return len(hash) == 128

	default:
		return false
	}
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
	sql = fmt.Sprintf(
		`SELECT image.* FROM image INNER JOIN hash_%s h on image.id = h.image_id 
				WHERE hash_is_within_distance%d(h.hash, $1, $2) 
				ORDER BY image.id LIMIT $3 OFFSET $4`,
		hashType, hashType.HashLength())

	rows, err := tx.Query(sql, hash, distance, limit, offset)
	if err != nil {
		return nil, err
	}

	var images []*Image
	for rows.Next() {
		var im Image
		err := rows.Scan(&im.id, &im.Crc32, &im.Size, &im.Sha1, &im.Md5, &im.Sha256)
		if err != nil {
			Logger.Error("Error while fetching db image", zap.String("err", err.Error()))
			return nil, err
		}

		images = append(images, &im)
	}

	if images == nil {
		b, _ := easyjson.Marshal(ImageList{Images: []*Image{}})
		return b, nil
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

			var compressedMeta []byte
			err := rows.Scan(&ihm.Url, &ihm.Meta.Id, &ihm.Meta.RetrievedAt, &compressedMeta)

			ihm.Meta.Meta, err = gozstd.DecompressDict(nil, compressedMeta, DDict)
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
	id BIGSERIAL PRIMARY KEY NOT NULL,
	crc32 bigint NOT NULL,
	size INT NOT NULL,
	sha1 bytea NOT NULL,
	md5 bytea NOT NULL,
	sha256 bytea NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_image_sha1 ON image(sha1);

CREATE TABLE IF NOT EXISTS image_meta (
	retrieved_at bigint NOT NULL,
	id TEXT PRIMARY KEY,
	meta bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS image_has_meta (
	image_id bigint REFERENCES image(id) NOT NULL,
	url TEXT NOT NULL,
	image_meta_id text REFERENCES image_meta(id) NOT NULL,
	UNIQUE(image_id, image_meta_id)
);

CREATE TABLE IF NOT EXISTS matchlist (
	id smallint,
	distance smallint NOT NULL,
	im1 bigint NOT NULL,
	im2 bigint NOT NULL
);
`
	for _, hashType := range HashTypes {
		sql += fmt.Sprintf(`CREATE TABLE IF NOT EXISTS hash_%s (
							image_id BIGINT REFERENCES image(id) UNIQUE NOT NULL,
							hash bytea NOT NULL);`, hashType)
	}

	for _, trigger := range MatchTriggers {
		sql += fmt.Sprintf(`
CREATE OR REPLACE FUNCTION on_%s_insert() RETURNS TRIGGER AS $$
BEGIN
	INSERT INTO matchlist (id, distance, im1, im2) 
		SELECT %d, hash_distance%d(hash, NEW.hash), NEW.image_id, image_id FROM hash_%s AS h
			WHERE h.image_id != NEW.image_id AND hash_is_within_distance%d(hash, NEW.hash, %d);
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
DROP TRIGGER IF EXISTS on_%s_insert ON hash_%s;
CREATE TRIGGER on_%s_insert AFTER INSERT ON hash_%s
FOR EACH ROW EXECUTE PROCEDURE on_%s_insert();`,
			trigger.HashType, trigger.Id, trigger.HashType.HashLength(), trigger.HashType,
			trigger.HashType.HashLength(), trigger.MinDistance, trigger.HashType,
			trigger.HashType, trigger.HashType, trigger.HashType, trigger.HashType)
	}

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
