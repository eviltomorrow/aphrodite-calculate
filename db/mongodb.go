package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
)

//
var (
	MongodbDSN     string
	MongodbMinOpen uint64 = 5
	MongodbMaxOpen uint64 = 10
	MongoDB        *mongo.Client
)

//
var (
	DefaultTimeout = 10 * time.Second
)

// BuildMongoDB build mongodb
func BuildMongoDB() {
	pool, err := build(MongodbDSN)
	if err != nil {
		zlog.Fatal("Build mongodb connection failure", zap.Error(err))
	}
	MongoDB = pool
}

// CloseMongoDB close mongodb
func CloseMongoDB() error {
	zlog.Info("Close mongodb connection", zap.String("dsn", MongodbDSN))
	if MongoDB == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	err := MongoDB.Disconnect(ctx)
	if err != nil {
		zlog.Error("Close mongodb connection failure", zap.Error(err))
	}
	return err
}

func build(dsn string) (*mongo.Client, error) {
	pool, err := mongo.NewClient(
		options.Client().ApplyURI(dsn).SetMaxPoolSize(MongodbMaxOpen).SetMinPoolSize(MongodbMinOpen),
	)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	if err := pool.Connect(ctx); err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	if err := pool.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return pool, nil
}
