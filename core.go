package imhashdb

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/simon987/fastimagehash-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"syscall"
)

const RedisPrefix = "q."
const UserAgent = "imhashdb/v1.0"
const Concurrency = 4

var ImageSuffixes = []string{
	".jpeg", ".jpg", ".png",
	".jpeg:orig", ".jpg:orig", ".png:orig",
	".bmp", ".webp",
}

var ImageBlackList = []string{}

var Rdb *redis.Client
var Pgdb *pgx.ConnPool
var Logger *zap.Logger

func Init() {
	Logger, _ = zap.NewDevelopment()

	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	Pgdb = DbConnect("localhost", 5432, "imhashdb", "imhashdb", "imhashdb")
	DbInit(Pgdb)
}

func ComputeHash(data []byte) (*fastimagehash.MultiHash, error) {
	h := &fastimagehash.MultiHash{}

	aHash, code := fastimagehash.AHashMem(data, 12)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("aHash error: %d", int(code))
	}
	dHash, code := fastimagehash.DHashMem(data, 12)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("dHash error: %d", int(code))
	}
	mHash, code := fastimagehash.MHashMem(data, 12)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("dHash error: %d", int(code))
	}
	pHash, code := fastimagehash.PHashMem(data, 12, 4)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("pHash error: %d", int(code))
	}
	wHash, code := fastimagehash.WHashMem(data, 8, 0, fastimagehash.Haar)
	if code != fastimagehash.Ok {
		return nil, errors.Errorf("wHash error: %d", int(code))
	}

	h.AHash = *aHash
	h.DHash = *dHash
	h.MHash = *mHash
	h.PHash = *pHash
	h.WHash = *wHash
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

func IsPermanentError(err error) bool {

	if strings.HasPrefix(err.Error(), "HTTP") {
		//TODO: Handle http 429 etc?
		return true
	}

	var opErr *net.OpError

	urlErr, ok := err.(*url.Error)
	if ok {
		opErr, ok = urlErr.Err.(*net.OpError)
		if !ok {
			if urlErr.Err != nil && urlErr.Err.Error() == "Proxy Authentication Required" {
				return true
			}
			return false
		}

		if opErr.Err.Error() == "Internal Privoxy Error" {
			return true
		}

	} else {
		_, ok := err.(net.Error)
		if ok {
			_, ok := err.(*net.DNSError)
			return ok
		}
	}

	if opErr == nil {
		return false
	}

	if opErr.Timeout() {
		// Usually means thalt there is no route to host
		return true
	}

	switch t := opErr.Err.(type) {
	case *net.DNSError:
		return true
	case *os.SyscallError:
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ECONNREFUSED:
				log.Println("connect refused")
				return true
			case syscall.ETIMEDOUT:
				log.Println("timeout")
				return false
			case syscall.ECONNRESET:
				log.Println("connection reset by peer")
				return false
			}
		}
	}

	return false
}
