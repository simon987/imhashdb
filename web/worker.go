package api

import (
	"context"
	"github.com/mailru/easyjson"
	. "github.com/simon987/imhashdb"
	"go.uber.org/zap"
	"time"
)

const inQueue = "qq:in"
const outQueue = "qq:out:"
const wipQueue = "qq:wip"

const CacheLength = time.Second * 5

func queryWorker() {
	Logger.Info("Query worker started")
	for {
		value := Rdb.BZPopMin(time.Second*30, inQueue).Val()
		if value == nil {
			continue
		}
		Logger.Info("worker query start")
		member := value.Member.(string)
		var req QueryReq
		_ = easyjson.Unmarshal([]byte(member), &req)

		resp, err := dbQuery(req, member)

		var b []byte
		if err != nil {
			Logger.Warn("worker query error", zap.Error(err))
			b, _ = easyjson.Marshal(QueryResp{
				Err: err.Error(),
			})
		} else {
			Logger.Info("worker query done")
			b = resp
		}
		Rdb.Set(outQueue+member, b, CacheLength)
	}
}

func dbQuery(req QueryReq, value string) ([]byte, error) {
	Rdb.SAdd(wipQueue, value)
	Rdb.Expire(wipQueue, time.Minute*10)

	defer Rdb.SRem(wipQueue, value)

	resp, err := FindImagesByHash(context.Background(), req.Hash, req.HashType, req.Distance, req.Limit, req.Offset)
	if err != nil {
		Logger.Error("Couldn't perform query")
		return nil, err
	}

	return resp, nil
}
