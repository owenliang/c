package master

import (
	"net/http"
	"net"
	"strconv"
	"time"
	"github.com/owenliang/c/common"
	"encoding/json"
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
	var (
		job common.Job
		oldJob *common.Job
		postJob string
		bytes []byte
		rawMsg json.RawMessage
		data *json.RawMessage
		err error
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要保存的job json串
	postJob = req.PostForm.Get("job")

	// json反序列化
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	// 保存job
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	// 如果是更新, 返回旧任务信息
	if oldJob != nil {
		if bytes, err = json.Marshal(oldJob); err == nil {
			rawMsg = json.RawMessage(bytes)
			data = &rawMsg
		}
	}

	// 返回成功应答
	if bytes, err = common.BuildResponse(0, "success", data); err == nil {
		resp.Write(bytes)
	}
	return

	// 返回异常应答
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), data); err == nil {
		resp.Write(bytes)
	}

	// 从命令行展示一下任务被成功保存
	// ETCDCTL_API=3 ./etcdctl get "/cron/jobs/" --prefix
}

//


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