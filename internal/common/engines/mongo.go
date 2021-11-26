package engines

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"time"
)

var (
	defaultMaxAggregateTime = 15 * time.Second
	defaultConnTimeout      = 10 * time.Second
)

type MongoConfig struct {
	User        string
	Password    string
	Addr        string
	DBName      string
	AuthSource  string
	PoolMaxSize uint64
	IdleTime    time.Duration
}

func NewMongoClient(conf MongoConfig) (*mongo.Database, error) {
	uri := getMongoUrl(conf)
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(conf.PoolMaxSize).
		SetMinPoolSize(1).
		SetMaxConnIdleTime(conf.IdleTime))
	if err != nil {
		return nil, err
	}

	// check the connection
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client.Database(conf.DBName), err
}

func getMongoUrl(conf MongoConfig) string {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=%s",
		conf.User, url.QueryEscape(conf.Password), conf.Addr, conf.DBName, conf.AuthSource)
	return uri
}
