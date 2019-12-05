package etcd

import(
	etcd_client "go.etcd.io/etcd/clientv3"
	"github.com/astaxie/beego/logs"
	//"go.etcd.io/etcd/mvcc/mvccpb"
	//"github.com/coreos/etcd/clientv3"
	"fmt"
	"time"
)
type EtcdClient struct{
	client *etcd_client.Client
	logKeys []string
}
var(
	etcdClient *EtcdClient
)
func ConnOpenCloseEtcd(){
	cli, err :=etcd_client.New(etcd_client.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	defer cli.Close()
	if err !=nil {
		logs.Error("connect failed, err:", err)
		return
	}
	fmt.Println("etcd connect succ +close()")

	etcdClient = &EtcdClient{
		client: cli,
	} 
}
func connOpenEtcd(){
	cli, err :=etcd_client.New(etcd_client.Config{
		Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	//defer cli.Close()
	if err !=nil {
		logs.Error("connect failed, err:", err)
		return
	}
	fmt.Println("etcd connect succ -no close()")

	etcdClient = &EtcdClient{
		client: cli,
	} 
}
