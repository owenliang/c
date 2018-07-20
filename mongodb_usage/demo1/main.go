package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

// github.com/mongodb/mongo-go-driver/mongo

// 初始化客户端
func demo1() {
	var (
		client *mongo.Client
		collection *mongo.Collection
		err error
	)

	// 建立连接, 5秒超时
	if client, err = mongo.Connect(context.TODO(), "mongodb://localhost:27017", clientopt.ConnectTimeout(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 选择哪个db的哪个collection
	collection = client.Database("cron").Collection("log")

	collection = collection
}

// 连接mongodb
func main() {
	demo1()
}
