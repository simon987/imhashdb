package imhashdb

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/simon987/fastimagehash-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"strings"
)

const RedisPrefix = "q."
const UserAgent = "imhashdb/v1.0"

var ImageSuffixes = []string{
	".jpeg", ".jpg", ".png",
	".jpeg:orig", ".jpg:orig", ".png:orig",
	".bmp", ".webp",
}

type Config struct {
	PgUser     string
	PgPassword string
	PgDb       string
	PgHost     string
	PgPort     int

	RedisAddr     string
	RedisPassword string
	RedisDb       int

	ApiAddr string

	HasherConcurrency int
	QueryConcurrency  int

	ImgurClientId string
	HasherPattern string
}

var ImageBlackList = []string{}

var Rdb *redis.Client
var Pgdb *pgx.ConnPool
var Logger *zap.Logger
var Conf Config

func Init() {
	Logger, _ = zap.NewDevelopment()

	Rdb = redis.NewClient(&redis.Options{
		Addr:     Conf.RedisAddr,
		Password: Conf.RedisPassword,
		DB:       Conf.RedisDb,
	})
	_, err := Rdb.Ping().Result()
	if err != nil {
		Logger.Fatal("Could not connect to redis server")
	}

	Pgdb = DbConnect(
		Conf.PgHost,
		Conf.PgPort,
		Conf.PgUser,
		Conf.PgPassword,
		Conf.PgDb,
	)
	DbInit(Pgdb)
}

func ComputeHash(data []byte) (*Hashes, error) {
	h := &Hashes{}
	var code fastimagehash.Code

	h.DHash8, code = fastimagehash.DHashMem(data, 8)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("dHash error: %d", int(code))
	}
	h.DHash16, code = fastimagehash.DHashMem(data, 16)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("dHash error: %d", int(code))
	}
	h.DHash32, code = fastimagehash.DHashMem(data, 32)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("dHash error: %d", int(code))
	}

	h.MHash8, code = fastimagehash.MHashMem(data, 8)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("mHash error: %d", int(code))
	}
	h.MHash16, code = fastimagehash.MHashMem(data, 16)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("mHash error: %d", int(code))
	}
	h.MHash32, code = fastimagehash.MHashMem(data, 32)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("mHash error: %d", int(code))
	}

	h.PHash8, code = fastimagehash.PHashMem(data, 8, 4)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("pHash error: %d", int(code))
	}
	h.PHash16, code = fastimagehash.PHashMem(data, 16, 4)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("pHash error: %d", int(code))
	}
	h.PHash32, code = fastimagehash.PHashMem(data, 32, 4)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("pHash error: %d", int(code))
	}

	h.WHash8, code = fastimagehash.WHashMem(data, 8, 0, fastimagehash.Haar)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("wHash error: %d", int(code))
	}
	h.WHash16, code = fastimagehash.WHashMem(data, 16, 0, fastimagehash.Haar)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("wHash error: %d", int(code))
	}
	h.WHash32, code = fastimagehash.WHashMem(data, 32, 0, fastimagehash.Haar)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("wHash error: %d", int(code))
	}

	return h, nil
}

func TransformLink(link string, meta *[]Meta) []string {
	for _, str := range ImageBlackList {
		if strings.Contains(link, str) {
			return nil
		}
	}

	links := handleImgurLink(link, meta)
	if links != nil {
		return links
	}

	return []string{link}
}
func isHttpOk(code int) bool {
	return code >= 200 && code < 300
}

func FetchJson(link string, v interface{}, raw *[]byte, headers ...[]string) error {

	body, err := Fetch(link, headers...)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	*raw = body

	return nil
}

func Fetch(link string, headers ...[]string) ([]byte, error) {
	client := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(link)
	req.Header.Add("User-Agent", UserAgent)
	for _, h := range headers {
		req.Header.Add(h[0], h[1])
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := client.Do(req, resp)

	if err != nil {
		Logger.Warn(
			"Error during HTTP request",
			zap.String("link", link),
			zap.String("err", err.Error()),
		)
		return nil, err
	}

	code := resp.StatusCode()
	if !isHttpOk(code) {
		Logger.Debug(
			"Got HTTP error code",
			zap.String("link", link),
			zap.Int("code", code),
		)
		return nil, errors.Errorf("HTTP %d", code)
	}

	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	Logger.Debug(
		"HTTP Get",
		zap.String("link", link),
		zap.Int("size", len(body)),
	)

	return body, nil
}
