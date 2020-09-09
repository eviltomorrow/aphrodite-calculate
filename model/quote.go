package model

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
)

// QueryQuoteOne query quote one
func QueryQuoteOne(db *mongo.Client, where map[string]interface{}) (*Quote, error) {
	var collection = db.Database(MongodbDatabaseName).Collection(CollectionNameQuote)
	ctx, cancel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cancel()

	result := collection.FindOne(ctx, where)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var quote = &Quote{}
	if err := result.Decode(quote); err != nil {
		return nil, err
	}
	return quote, nil
}

// Quote quote
type Quote struct {
	ObjectID        string  `json:"_id" bson:"_id"`
	Source          string  `json:"source" bson:"source"`                     // 来源
	Code            string  `json:"code" bson:"code"`                         // 代码
	Name            string  `json:"name" bson:"name"`                         // 名称
	Open            float64 `json:"open" bson:"open"`                         // 开盘价格
	YesterdayClosed float64 `json:"yesterday_closed" bson:"yesterday_closed"` // 昨日收盘价格
	High            float64 `json:"high" bson:"high"`                         // 最高价
	Low             float64 `json:"low" bson:"low"`                           // 最低价
	Volume          uint64  `json:"volume" bson:"volume"`                     // 成交量
	Account         float64 `json:"account" bson:"account"`                   // 成交额
	Date            string  `json:"date" bson:"date"`                         // 日期
	Suspend         string  `json:"suspend" bson:"suspend"`                   // 停盘状态
	CreateTimestamp int64   `json:"create_timestamp" bson:"create_timestamp"` // 创建时间
	ModifyTimestamp int64   `json:"modify_timestamp" bson:"modify_timestamp"` // 修改时间
}

func (q *Quote) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
