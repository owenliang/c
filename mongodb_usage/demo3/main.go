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

// 插入N条记录
func demo3() {
	var (
		client *mongo.Client
		collection *mongo.Collection
		logRecord *LogRecord
		manyResult *mongo.InsertManyResult
		insertId interface{}
		docId objectid.ObjectID
		logArr []interface{}
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

	// 要写入的3条日志
	logArr = []interface{}{logRecord, logRecord, logRecord}

	// 批量插入mongodb
	if manyResult, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
		return
	}

	// 遍历打印自增ID
	for _, insertId = range manyResult.InsertedIDs {
		docId = insertId.(objectid.ObjectID)
		fmt.Println("插入ID:", docId)
	}
}

// 连接mongodb
func main() {
	demo3()
}
