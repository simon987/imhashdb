package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"hash/crc32"
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

func dispatchFromQueue(pattern string, queue chan []string) {

	for {
		keys, err := Rdb.Keys(pattern).Result()
		if err != nil {
			Logger.Error("Could not get keys for pattern", zap.String("pattern", pattern))
			continue
		}

		rawTask, err := Rdb.BLPop(time.Second*30, keys...).Result()
		if err != nil {
			continue
		}

		queue <- rawTask
	}
}

func worker(queue chan []string) {
	for rawTask := range queue {
		computeAndStore(rawTask)
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
		Id:          rawTask[0][len(RedisPrefix):] + "." + strconv.FormatInt(task.Id, 10),
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
				if !IsPermanentError(err) {
					// Retry later
					Logger.Debug("Will retry task later", zap.String("link", link))
					Rdb.RPush(rawTask[0], rawTask[1])
				}
				continue
			}

			if len(data) == 0 {
				continue
			}

			h, err := ComputeHash(data)
			if err != nil {
				return
			}

			Store(&Entry{
				H:      h,
				Size:   len(data),
				Sha256: sha256.Sum256(data),
				Sha1:   sha1.Sum(data),
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

func main() {
	Init()

	_, err := Rdb.Ping().Result()
	if err != nil {
		Logger.Fatal("Could not connect to redis server")
	}

	queue := make(chan []string, 100)

	for i := 0; i < Concurrency; i++ {
		go worker(queue)
	}

	dispatchFromQueue("q.reddit.*", queue)
}
