package master

import (
	"io/ioutil"
	"encoding/json"
)

// 程序配置
type Config struct {
	ApiPort int `json:"apiPort"`	// API监听端口
	ApiReadTimeout int `json:"apiReadTimeout"`	// API读超时
	ApiWriteTimeout int `json:"apiWriteTimeout"`	// API写超时
	EtcdEndpoints []string `json:etcdEndpoints` // etcd集群列表
	EtcdDialTimeout int `json:etcdDialTimeout` // etcd连接超时
	Webroot string `json:"webroot"`	// 静态页面根目录
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