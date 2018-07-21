package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
	"fmt"
	"time"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
)

// 大写字段是要导出的字段
type TimePoint struct {
	StartTime int64	`bson:"startTime"`
	EndTime int64	`bson:"endTime"`
}

type LogRecord struct {
	JobName string `bson:"jobName"`// 任务名字
	Command string `bson:"command"`// 脚本命令
	Err string `bson:"err"`// 错误原因
	Content string  `bson:"content"` // shell输出内容
	TimePoint TimePoint `bson:"timePoint"` // 执行时间信息
}

// 插入一条记录
func demo2() {
	var (
		client *mongo.Client
		collection *mongo.Collection
		logRecord *LogRecord
		result *mongo.InsertOneResult
		docId objectid.ObjectID
		err error
	)

	if client, err = mongo.Connect(context.TODO(), "mongodb://36.111.184.221:27017", clientopt.ConnectTimeout(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 选择哪个db的哪个collection
	collection = client.Database("cron").Collection("log")

	// 要插入的记录
	logRecord = &LogRecord{
		JobName: "job10",
		Command: "echo 123",
		Err: "",
		Content: "123",
		TimePoint: TimePoint{time.Now().Unix(), time.Now().Unix() + 10},
	}

	// 插入到Mongodb
	if result, err = collection.InsertOne(context.TODO(), logRecord); err != nil {
		fmt.Println(err)
		return
	}

	// 得到自增ID
	docId = result.InsertedID.(objectid.ObjectID)	// 这是一个12字节的二进制ID
	fmt.Println("记录唯一ID:", docId.Hex()) // 十六进制表达
}

// 连接mongodb
func main() {
	demo2()
}
