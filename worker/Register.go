package worker

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"golang.org/x/net/context"
	"net"
	"github.com/owenliang/c/common"
)

// 注册服务到etcd
type Register struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease

	localIP string // 本机IP
}

var (
	G_register *Register
)

// 工具方法, 获取本机IP
func getLocalIp() (ipv4 string, err error) {
	var (
		addrs []net.Addr
		addr net.Addr
		ipNet *net.IPNet
		isIpNet bool
	)

	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {	// 只接受IPV4
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}

// 自动注册到/cron/workers/
func (register *Register) keepOnline() {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		keepAliveChan <- chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp *clientv3.LeaseKeepAliveResponse
		regKey string
		// putResp *clientv3.PutResponse
		err error
	)

	// 注册路径
	regKey = common.JOB_WORKER_DIR + register.localIP

	// 持续保持在线
	for {
		// 创建租约
		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		// 自动续约
		leaseId = leaseGrantResp.ID
		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseId); err != nil {
			goto RETRY
		}

		// 建立节点
		if _, err = register.kv.Put(context.TODO(), regKey, "", clientv3.WithLease(leaseId)); err != nil {
			goto RETRY
		}

		// 处理续租应答
		for {
			select {
			case keepAliveResp = <- keepAliveChan:
				if keepAliveResp == nil {	// 续租中断
					goto RETRY
				}
			}
		}

		RETRY:
		time.Sleep(1 * time.Second)
	}
}

func InitRegister() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		localIP string
	)

	// 初始化etcd客户端配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,	// 集群列表
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,	// 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 获取本机IP
	if localIP, err = getLocalIp(); err != nil {
		return
	}

	G_register = &Register{
		client: client,
		localIP: localIP,
		kv: kv,
		lease: lease,
	}

	// 保持在线
	go G_register.keepOnline()
	return
}
