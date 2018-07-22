package worker

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"golang.org/x/net/context"
	"github.com/owenliang/c/common"
	"fmt"
)

// mongodb保存日志
type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan *common.JobLog
	autoCommitChan chan *LogBatch
}

// 打包日志, 提升吞吐
type LogBatch struct {
	logs []interface{}	// 多行日志
}

var (
	G_logSink *LogSink
)

// 保存日志
func (logSink *LogSink) saveLogs(batch *LogBatch) {
	fmt.Println("保存日志:", len(batch.logs))
	logSink.logCollection.InsertMany(context.TODO(), batch.logs)
}

// 日志存储线程
func (logSink *LogSink) writeLoop() {
	var (
		log* common.JobLog
		commitTimer *time.Timer
		logBatch *LogBatch	// 当前批次
		timeoutBatch *LogBatch // 过期批次
	)

	// 延迟1秒批量提交
	for {
		select {
		// 新到来的日志, 加入到批次
		case log = <- logSink.logChan:
			if logBatch == nil {
				logBatch = &LogBatch{}
				// 启动定时器
				commitTimer = time.AfterFunc(time.Duration(G_config.JogLogCommitTimeout) * time.Millisecond, func(batch *LogBatch) func() {
					return func() {
						logSink.autoCommitChan <- logBatch
					}
				}(logBatch))
			}

			// 日志加入批次
			logBatch.logs = append(logBatch.logs, log)

			// 如果批次满了, 立即发送
			if len(logBatch.logs) >= G_config.JobLogBatchSize {
				//发送日志
				logSink.saveLogs(logBatch)
				// 清空批次
				logBatch = nil
				// 停止自动提交
				commitTimer.Stop()
			}

		case timeoutBatch = <- logSink.autoCommitChan:	// 到期的批次, 判断是否自动提交
			if timeoutBatch != logBatch {	// 一个没来得及取消的旧定时器, 忽略
				continue
			}
			// 批量写入到mongo
			logSink.saveLogs(logBatch)
			// 下次启用新批次
			logBatch = nil
		}
	}
}

// 初始化mongo日志存储
func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	// 建立mongodb连接
	if client, err = mongo.Connect(
			context.TODO(),
			G_config.MongodbUri,
			clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout) * time.Second)); err != nil {
		return
	}

	G_logSink = &LogSink{
		client: client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan: make(chan *common.JobLog, G_config.JobLogChanSize),
		autoCommitChan: make(chan*LogBatch, 1000),
	}

	// 启动mongodb保存协程
	go G_logSink.writeLoop()
	return
}

// 发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	fmt.Println("发送日志", *jobLog)
	select {
	case logSink.logChan <- jobLog:
	default:
		// 队列满了就丢掉
	}
}