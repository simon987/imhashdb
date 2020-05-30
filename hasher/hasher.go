package hasher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "github.com/simon987/imhashdb"
	"go.uber.org/zap"
)

type Task struct {
	Urls []string `json:"_urls"`
	Id   int64    `json:"_id"`
}

func dispatchFromQueue(pattern string, queue chan []string) error {

	for {
		keys, err := Rdb.Keys(pattern).Result()
		if err != nil {
			Logger.Error("Could not get keys for pattern", zap.String("pattern", pattern))
			continue
		}

		if len(keys) == 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		rawTask, err := Rdb.BLPop(time.Second*30, keys...).Result()
		if err != nil {
			continue
		}

		//TODO: put in WIP list, resume on crash
		queue <- rawTask
	}

	return nil
}

func worker(queue chan []string) {
	for rawTask := range queue {
		computeAndStore(rawTask)
	}
}

func storeData(data []byte, sha1 [20]byte, link string) {

	sha1Str := hex.EncodeToString(sha1[:])

	filename := fmt.Sprintf("%s/%c/%s/",
		DataPath,
		sha1Str[0],
		sha1Str[1:3],
	)
	err := os.MkdirAll(filename, 0755)
	if err != nil {
		panic(err)
	}
	filename += sha1Str + filepath.Ext(link)

	Logger.Debug("Storing image data to file", zap.String("path", filename))

	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		panic(err)
	}
}

func computeAndStore(rawTask []string) {
	var task Task
	err := json.Unmarshal([]byte(rawTask[1]), &task)
	if err != nil {
		Logger.Error("Corrupt task body", zap.String("body", rawTask[1]))
		return
	}

	meta := []Meta{{
		RetrievedAt: time.Now().Unix(),
		Id:          rawTask[0][len(Pattern)-1:] + "." + strconv.FormatInt(task.Id, 10),
		Meta:        []byte(rawTask[1]),
	}}

	for _, link := range task.Urls {
		for _, turl := range TransformLink(link, &meta) {
			if !IsImageLink(turl) {
				Logger.Debug("Ignoring non-image URL", zap.String("link", link))
				continue
			}

			data, err := Fetch(turl)
			if err != nil {
				continue
			}

			if len(data) == 0 {
				continue
			}

			h, err := ComputeHash(data)
			if err != nil {
				return
			}

			sha1sum := sha1.Sum(data)

			if StoreData {
				storeData(data, sha1sum, link)
			}

			Store(&Entry{
				H:      h,
				Size:   len(data),
				Sha256: sha256.Sum256(data),
				Sha1:   sha1sum,
				Md5:    md5.Sum(data),
				Crc32:  crc32.ChecksumIEEE(data),
				Meta:   meta,
				Url:    trimUrl(turl),
			})
		}
	}
}

func trimUrl(link string) string {
	if strings.HasPrefix(link, "https://") {
		return link[len("https://"):]
	} else if strings.HasPrefix(link, "http://") {
		return link[len("http://"):]
	}

	return link
}

var StoreData bool
var DataPath string
var Pattern = "imhash.*"

func Main() error {
	StoreData = Conf.Store != ""
	DataPath = Conf.Store

	queue := make(chan []string)

	for i := 0; i < Conf.HasherConcurrency; i++ {
		go worker(queue)
	}

	return dispatchFromQueue(Pattern, queue)
}
