package master

import (
	"net/http"
	"net"
	"strconv"
	"time"
)

// 对Web提供的HTTP API接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	// 单例
	G_apiServer *ApiServer
)

/** 内部实现 **/

// 保存任务
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	
}


/** 对外接口 **/

// 初始化API服务
func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
		httpServer *http.Server
		listener net.Listener
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 监听端口
	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建服务
	httpServer = &http.Server{
		ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler: mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 拉起服务
	go httpServer.Serve(listener)

	return
}