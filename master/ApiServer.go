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

// 删除任务
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		name string
		err error
		oldJob *common.Job
		bytes []byte
		rawMsg json.RawMessage
		data *json.RawMessage
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要删除的任务名
	name = req.PostForm.Get("name")

	// 删除任务
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	// 如果删除成功, 返回被删除的任务信息
	if oldJob != nil {
		// 忽略这种错误
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
}

// 列举所有任务
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*common.Job
		bytes []byte
		rawMsg json.RawMessage
		data *json.RawMessage
		err error
	)

	// 获取任务列表
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	// 序列化data
	if bytes, err = json.Marshal(jobList); err != nil {
		goto ERR
	}

	rawMsg = json.RawMessage(bytes)
	data = &rawMsg

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
}

// 强制杀死任务
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		name string
		bytes []byte
		err error
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要杀死的任务名
	name = req.PostForm.Get("name")

	// 删除任务
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}

	// 返回成功应答
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}
	return

	// 返回异常应答
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

/** 对外接口 **/

// 初始化API服务
func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
		httpServer *http.Server
		listener net.Listener
		staticDir http.Dir	// 静态文件目录
		staticHandler http.Handler // 静态文件响应器
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	// 静态文件路由
	staticDir = http.Dir(G_config.Webroot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

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