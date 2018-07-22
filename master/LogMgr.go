package master

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"golang.org/x/net/context"
	"github.com/owenliang/c/common"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// mongodb日志管理
type LogMgr struct {
	client  *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

// 初始化
func InitLogMgr() (err error) {
	var (
		client *mongo.Client
	)

	// 建立mongodb连接
	if client, err = mongo.Connect(
		context.TODO(),
		G_config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Second)); err != nil {
		return
	}

	G_logMgr = &LogMgr{
		client: client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

// 查询任务执行日志
func (logMgr *LogMgr) ListLog(filter *common.JobLogFilter, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		cursor mongo.Cursor
		jobLog *common.JobLog
		logSort *common.SortLogByStartTime
	)

	// 初始化日志数组
	logArr = make([]*common.JobLog, 0)

	// 按开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder: -1}

	// mongo查询
	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, findopt.Skip(int64(skip)), findopt.Limit(int64(limit)), findopt.Sort(logSort)); err != nil {
		return
	}

	// 延迟释放游标
	defer cursor.Close(context.TODO())

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		// 反序列化到对象
		if err = cursor.Decode(jobLog); err != nil {
			return
		}

		// 记录返回结果
		logArr = append(logArr, jobLog)
	}
	return
}