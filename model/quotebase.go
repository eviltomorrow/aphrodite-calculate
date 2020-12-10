package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SelectQuoteBaseCurrentCodeLimit2 query quote base limit 2
func SelectQuoteBaseCurrentCodeLimit2(db *mongo.Client, code string, date string) ([]*QuoteBase, error) {
	if strings.Count(date, "-") != 2 {
		return nil, fmt.Errorf("Invalid date, date: %s", date)
	}

	var year = strings.Split(date, "-")[0]

	var collection = db.Database(MongodbDatabaseName).Collection(fmt.Sprintf(CollectionNameQuote, year))
	ctx, cancel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cancel()

	var option = options.Find()
	option.SetSort(bson.D{{Key: "date", Value: 1}})
	option.SetLimit(2)
	cur, err := collection.Find(ctx, bson.M{
		"code": code,
		"date": bson.M{
			"$gte": date,
		},
	}, option)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var quotes = make([]*QuoteBase, 0, 2)
	for cur.Next(ctx) {
		var quote = &QuoteBase{}
		if err := cur.Decode(quote); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if len(quotes) == 2 {
		return quotes, nil
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, fmt.Errorf("Parse year to int failure, nest error: %v", err)
	}

	year = fmt.Sprintf("%d", yearInt+1)
	collection = db.Database(MongodbDatabaseName).Collection(fmt.Sprintf(CollectionNameQuote, year))

	option = options.Find()
	option.SetSort(bson.D{{Key: "date", Value: 1}})
	option.SetLimit(1)
	cur, err = collection.Find(ctx, bson.M{
		"code": code,
		"date": bson.M{
			"$gte": date,
		},
	}, option)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var quote = &QuoteBase{}
		if err := cur.Decode(quote); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}

	return quotes, nil
}

// QuoteBase quote base
type QuoteBase struct {
	ObjectID        string  `json:"_id" bson:"_id"`
	Source          string  `json:"source" bson:"source"`                     // 来源
	Code            string  `json:"code" bson:"code"`                         // 代码
	Name            string  `json:"name" bson:"name"`                         // 名称
	Open            float64 `json:"open" bson:"open"`                         // 开盘价格
	YesterdayClosed float64 `json:"yesterday_closed" bson:"yesterday_closed"` // 昨日收盘价格
	High            float64 `json:"high" bson:"high"`                         // 最高价
	Low             float64 `json:"low" bson:"low"`                           // 最低价
	Volume          int64   `json:"volume" bson:"volume"`                     // 成交量
	Account         float64 `json:"account" bson:"account"`                   // 成交额
	Date            string  `json:"date" bson:"date"`                         // 日期
	Suspend         string  `json:"suspend" bson:"suspend"`                   // 停盘状态
	CreateTimestamp int64   `json:"create_timestamp" bson:"create_timestamp"` // 创建时间
	ModifyTimestamp int64   `json:"modify_timestamp" bson:"modify_timestamp"` // 修改时间
}

func (q *QuoteBase) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
