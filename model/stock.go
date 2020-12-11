package model

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/eviltomorrow/aphrodite-calculate/db"
)

// SelectStockManyByCodesForMySQL select stock list with code for mysql
func SelectStockManyByCodesForMySQL(db db.ExecMySQL, codes []string) ([]*Stock, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var fields = make([]string, 0, len(codes))
	var args = make([]interface{}, 0, len(codes))
	for _, code := range codes {
		fields = append(fields, "?")
		args = append(args, code)
	}
	var _sql = fmt.Sprintf(`select code, name, source, create_timestamp, modify_timestamp from stock where code in (%s)`, strings.Join(fields, ","))
	rows, err := db.QueryContext(ctx, _sql, args...)
	if err != nil {
		return nil, err
	}

	var stocks = make([]*Stock, 0, len(codes))
	for rows.Next() {
		var stock = Stock{}
		if err := rows.Scan(&stock.Code, &stock.Name, &stock.Source, &stock.CreateTimestamp, &stock.ModifyTimestamp); err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return stocks, nil
}

// UpdateStockByCodeForMySQL update stock by code for mysql
func UpdateStockByCodeForMySQL(db db.ExecMySQL, code string, stock *Stock) (int64, error) {
	ctx, cannel := context.WithTimeout(context.Background(), UpdateTimeout)
	defer cannel()

	var _sql = `update stock set name = ?, source = ?, modify_timestamp = now() where code = ?`
	result, err := db.ExecContext(ctx, _sql, stock.Name, stock.Source, code)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SelectStockManyForMySQL select stock list for mysql
func SelectStockManyForMySQL(db *sql.DB, offset, limit int64) ([]*Stock, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = `select code, name, source from stock limit ?, ?`
	rows, err := db.QueryContext(ctx, _sql, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks = make([]*Stock, 0, limit)
	for rows.Next() {
		var stock = &Stock{}
		if err := rows.Scan(&stock.Code, &stock.Name, &stock.Source); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

// SelectStockManyForMongoDB select stock list for mongodb
func SelectStockManyForMongoDB(db *mongo.Client, offset, limit int64, lastID string) ([]*Stock, error) {
	if limit <= 0 {
		return []*Stock{}, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var collection = db.Database(MongodbDatabaseName).Collection(CollectionNameStock)

	var opt = &options.FindOptions{}
	opt.SetLimit(limit)

	var filter = bson.M{}
	if lastID != "" {
		objectID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": bson.M{"$gt": objectID}}
	} else {
		opt.SetSkip(offset)
	}

	// fmt.Println(filter)
	cur, err := collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var stocks = make([]*Stock, 0, limit)
	for cur.Next(context.Background()) {
		var stock = &Stock{}
		if err := cur.Decode(stock); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}
	return stocks, nil
}

// InsertStockManyForMySQL insert stock many for mysql
func InsertStockManyForMySQL(db db.ExecMySQL, stocks []*Stock) (int64, error) {
	if len(stocks) == 0 {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), InsertTimeout)
	defer cannel()

	var fields = make([]string, 0, len(stocks))
	var args = make([]interface{}, 0, 3*len(stocks))
	for _, stock := range stocks {
		fields = append(fields, "(?, ?, ?, now(), null)")
		args = append(args, stock.Code)
		args = append(args, stock.Name)
		args = append(args, stock.Source)
	}

	var _sql = fmt.Sprintf("insert into stock (%s) values %s", strings.Join(stockFields, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//
const (
	StockFieldCode            = "code"
	StockFieldName            = "name"
	StockFieldSource          = "source"
	StockFieldCreateTimestamp = "create_timestamp"
	StockFieldModifyTimestamp = "modify_timestamp"
)

var stockFields = []string{
	StockFieldCode,
	StockFieldName,
	StockFieldSource,
	StockFieldCreateTimestamp,
	StockFieldModifyTimestamp,
}

// Time time
type Time time.Time

// UnmarshalJSON unmarshal json
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	num, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*t = Time(time.Unix(int64(num), 0))
	return nil
}

// MarshalJSON marshal json
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalBSON unmarshal bson
func (t *Time) UnmarshalBSON(data []byte) error {
	var num int64
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.LittleEndian, &num); err != nil {
		return err
	}
	if num == 0 {
		return nil
	}
	*t = Time(time.Unix(int64(num), 0))
	return nil
}

// MarshalBSON marshal bson
func (t Time) MarshalBSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t Time) String() string {
	if time.Time(t).IsZero() {
		return `""`
	}
	return time.Time(t).Format("2006-01-02 15:04:05")
}

// Scan scan
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch value.(type) {
	case time.Time:
		*t = Time(value.(time.Time))
	default:
	}
	return nil
}

// Value value
func (t Time) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return "", nil
	}
	return t.String(), nil
}

// Stock stock
type Stock struct {
	ObjectID        string `json:"_id" bson:"_id"`
	Code            string `json:"code" bson:"code"`
	Name            string `json:"name" bson:"name"`
	Source          string `json:"source" bson:"source"`
	Valid           bool   `json:"valid,omitempty" bson:"valid"`
	CreateTimestamp Time   `json:"create_timestamp,omitempty" bson:"create_timestamp"`
	ModifyTimestamp Time   `json:"modify_timestamp,omitempty" bson:"modify_timestamp"`
}

func (s *Stock) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s)
	return string(buf)
}
