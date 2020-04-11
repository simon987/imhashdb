package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/mailru/easyjson"
	. "github.com/simon987/imhashdb"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func submitQuery(value string) bool {
	if Rdb.SIsMember(wipQueue, value).Val() {
		return false
	}

	if Rdb.Exists(outQueue+value).Val() == 1 {
		return false
	}

	Rdb.ZAdd(inQueue, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: value,
	})
	return true
}

func pollQuery(ctx context.Context, reqStr string) ([]byte, error) {
	key := outQueue + reqStr

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		default:
		}

		value, _ := Rdb.Get(key).Bytes()
		if value != nil {
			return value, nil
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func query(c *gin.Context) {
	var req QueryReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"err": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	reqJson, _ := easyjson.Marshal(req)
	value := string(reqJson)

	submitQuery(value)
	b, err := pollQuery(ctx, value)
	if err != nil {
		b, _ = easyjson.Marshal(QueryResp{
			Err: err.Error(),
		})
	}
	c.Data(200, gin.MIMEJSON, b)
}

func hash(c *gin.Context) {
	var req HashReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"err": "Invalid request"})
		return
	}

	h, err := ComputeHash(req.Data)
	if err != nil {
		c.JSON(500, gin.H{"err": "Couldn't compute image hash"})
		return
	}

	b, _ := easyjson.Marshal(HashResp{
		AHash: h.AHash.Bytes, DHash: h.DHash.Bytes,
		MHash: h.MHash.Bytes, PHash: h.PHash.Bytes,
		WHash: h.WHash.Bytes,
	})
	c.Data(200, gin.MIMEJSON, b)
}

func main() {
	Init()

	f, err := os.Create("prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	go func() {
		time.Sleep(time.Second * 15)
		pprof.StopCPUProfile()
		fmt.Println("!!!!!!!!!!!!!!!")
		f.Close()
	}()

	r := gin.Default()
	r.POST("/api/hash", hash)
	r.POST("/api/query", query)

	//TODO: concurrency
	go queryWorker()

	r.Run()
}
