package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertStockManyForMySQL insert stock many for mysql
func InsertStockManyForMySQL(db *sql.DB, stocks []*Stock) (int64, error) {
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

	var _sql = fmt.Sprintf("insert into stock (%s) values %s", strings.Join(stockFeilds, ","), strings.Join(fields, ","))
	result, err := db.ExecContext(ctx, _sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateStockByCodeForMySQL update stock by code for mysql
func UpdateStockByCodeForMySQL(db *sql.DB, code string, stock *Stock) (int64, error) {
	ctx, cannel := context.WithTimeout(context.Background(), UpdateTimeout)
	defer cannel()

	var _sql = `update stock set name = ?, source = ?, modify_timestamp = now() where code = ?`
	stmt, err := db.PrepareContext(ctx, _sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(stock.Name, stock.Source, stock.Code)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// QueryStockOneForMySQL query stock one for mysql
func QueryStockOneForMySQL(db *sql.DB, code string) (*Stock, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = `select code, name, source, create_timestamp, modify_timestamp from stock where code = ?`
	row := db.QueryRowContext(ctx, _sql, code)
	// if err := row.Err(); err != nil {
	// 	return nil, err
	// }

	var stock = &Stock{}
	if err := row.Scan(&stock.Code, &stock.Name, &stock.Source, &stock.CreateTimestamp, &stock.ModifyTimestamp); err != nil {
		return nil, err
	}
	return stock, nil
}

// QueryStockListForMySQL query stock list for mysql
func QueryStockListForMySQL(db *sql.DB, offset, limit int64) ([]*Stock, error) {
	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var _sql = `select code, name, source from stock limit ?, ?`
	rows, err := db.QueryContext(ctx, _sql, offset, limit)
	if err != nil {
		return nil, err
	}

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

// QueryStockListForMongoDB query stock list for mongodb
func QueryStockListForMongoDB(db *mongo.Client, offset, limit int64) ([]*Stock, error) {
	if limit == 0 {
		return []*Stock{}, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), SelectTimeout)
	defer cannel()

	var collection = db.Database(MongodbDatabaseName).Collection(CollectionNameStock)

	var opt = &options.FindOptions{}
	opt.SetSkip(offset)
	opt.SetLimit(limit)
	cur, err := collection.Find(ctx, bson.M{}, opt)
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

//
const (
	StockFeildCode            = "code"
	StockFeildName            = "name"
	StockFeildSource          = "source"
	StockFeildCreateTimestamp = "create_timestamp"
	StockFeildModifyTimestamp = "modify_timestamp"
)

var stockFeilds = []string{
	StockFeildCode,
	StockFeildName,
	StockFeildSource,
	StockFeildCreateTimestamp,
	StockFeildModifyTimestamp,
}

// Stock stock
type Stock struct {
	ObjectID        string `json:"_id" bson:"_id"`
	Code            string `json:"code" bson:"code"`
	Name            string `json:"name" bson:"name"`
	Source          string `json:"source" bson:"source"`
	Valid           bool   `json:"valid" bson:"valid"`
	CreateTimestamp int64  `json:"create_timestamp" bson:"create_timestamp"`
	ModifyTimestamp int64  `json:"modify_timestamp" bson:"modify_timestamp"`
}

func (s *Stock) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(s)
	return string(buf)
}
