package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/mailru/easyjson"
	. "github.com/simon987/imhashdb"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

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

	if req.Data == nil {
		c.JSON(400, gin.H{"err": "Invalid request"})
		return
	}

	h, err := ComputeHash(req.Data)
	if err != nil {
		c.JSON(500, gin.H{"err": "Couldn't compute image hash"})
		return
	}

	b, _ := easyjson.Marshal(h)
	c.Data(200, gin.MIMEJSON, b)
}

func Main() error {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/api/hash", hash)
	r.POST("/api/query", query)

	for i := 0; i < Conf.QueryConcurrency; i++ {
		go queryWorker()
	}

	return r.Run(Conf.ApiAddr)
}
