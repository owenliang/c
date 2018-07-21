package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

// 小于某时间
// $lt: timestamp
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// 开始时间小于某时间
// timePoint.startTime: {$lt: timestamp}
type DeleteCond struct {
	StartBeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

// 查找记录
func demo5() {
	var (
		client *mongo.Client
		collection *mongo.Collection
		deleteCond *DeleteCond
		delResult *mongo.DeleteResult
		err error
	)

	// 建立连接, 5秒超时
	if client, err = mongo.Connect(context.TODO(), "mongodb://36.111.184.221:27017", clientopt.ConnectTimeout(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 选择哪个db的哪个collection
	collection = client.Database("cron").Collection("log")

	// 筛选出当前时间之前的记录
	deleteCond = &DeleteCond{StartBeforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	// 查看一下查询请求是否符合期望
	// bson, err := mongo.TransformDocument(deleteCond)
	// fmt.Println(bson.ToExtJSON(false))

	// 执行删除
	if delResult, err = collection.DeleteMany(context.TODO(), deleteCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除行数:", delResult.DeletedCount)
}

// 连接mongodb
func main() {
	demo5()
}
