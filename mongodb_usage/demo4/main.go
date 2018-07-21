package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

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

// 按jobName字段过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

// 查找记录
func demo4() {
	var (
		client *mongo.Client
		collection *mongo.Collection
		findByJobName *FindByJobName
		cursor mongo.Cursor
		logRecord *LogRecord
		err error
	)

	// 建立连接, 5秒超时
	if client, err = mongo.Connect(context.TODO(), "36.111.184.221:27017", clientopt.ConnectTimeout(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 选择哪个db的哪个collection
	collection = client.Database("cron").Collection("log")

	// 准备过滤条件
	findByJobName = &FindByJobName{JobName: "job10"}

	// 发起查询(过滤+翻页)
	if cursor, err = collection.Find(context.TODO(), findByJobName, findopt.Skip(0), findopt.Limit(5)); err != nil {
		fmt.Println(err)
		return
	}

	// 最后释放游标资源
	defer cursor.Close(context.TODO())

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		logRecord = &LogRecord{}

		// 反序列化到对象
		if err = cursor.Decode(logRecord); err != nil {
			fmt.Println(err)
			return
		}

		// 打印读到的记录
		fmt.Println(*logRecord)
	}
}

// 连接mongodb
func main() {
	demo4()

	// 需要为查询字段建立索引
	// db.log.createIndex({"jobName": 1})

	// 查看建立的索引
	// db.log.getIndexes()
}
