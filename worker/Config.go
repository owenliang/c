package worker

import (
	"io/ioutil"
	"encoding/json"
)

// 程序配置
type Config struct {
	EtcdEndpoints []string `json:etcdEndpoints` // etcd集群列表
	EtcdDialTimeout int `json:etcdDialTimeout` // etcd连接超时
	JobEventChanSize int `json:"jobEventChanSize"`	// 任务事件队列长度
	MongodbUri string `json:"mongodbUri"` // mongodb地址
	MongodbConnectTimeout int `json:"mongodbConnectTimeout"` // mongodb连接超时
	JobLogChanSize int `json:"jobLogChanSize"` // 执行日志队列长度
	JobLogBatchSize int `json:"jobLogBatchSize"`	// 日志写入批次大小
	JogLogCommitTimeout int `json:"jogLogCommitTimeout"` // 在未达到批次大小前, 超时自动提交(毫秒)
}

var (
	// 单例
	G_config *Config
)

// 加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf Config
	)

	// 读配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// JSON反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf
	return
}