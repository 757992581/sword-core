package etd

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/moka-mrp/sword-core/config"
	"testing"
)


var client *Client
//-------------------------------- init -----------------------------------------------
func init() {
	//初始化的时候会注入的
	etcdInit(true)
	//从容器中获取资源
	client= GetEtcd()
}

//注入容器以及从容器中快速取出来
func etcdInit(lazyBool bool) {

	etcdConf := config.EtcdConfig{
		Endpoints:  []string{"127.0.0.1:2379"},
		Username:    "",
		Password:    "",
		DialTimeout: 2,
		ReqTimeout:  3,
	}
	//测试容器注入功能(容器本身已经自动在kernel/container/app.go中初始化好了)
	err := Pr.Register(SingletonMain, etcdConf, lazyBool)
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------添加----------------------------------------------
//测试添加一个key的操作
//@author sam@2020-08-21 16:00:13
func TestPut(t *testing.T){
	if putResp, err := client.Put("/test/food", "apple", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision) //整体的版本号
		if putResp.PrevKv != nil {	// 打印修改之前的值
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}


//测试添加一个key的操作(持乐观锁)
//@author sam@2020-08-21 16:37:26
func TestPutWithModRev(t *testing.T){
	if putResp, err := client.PutWithModRev("/test/food", "apple2",1164, clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {	// 打印修改之前的值
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}

//-----------------------------------------查看----------------------------------------------

//测试添加一个key的操作(持乐观锁)
//@author sam@2020-08-21 16:37:26
func TestGet(t *testing.T){
	if getResp, err := client.Get( "/test/food", /*clientv3.WithCountOnly()*/); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
		//[key:"/test/food" create_revision:1140 mod_revision:1163 version:12 value:"apple" ] 1
	}


}

//-----------------------------------------删除----------------------------------------------

func TestDelete(t *testing.T)  {
	// 删除KV   clientv3.WithFromKey()
	if delResp, err := client.Delete("/test/food", clientv3.WithFromKey()); err != nil {
		fmt.Println(err)
		return
	}else{
		fmt.Printf("%+v\r\n",delResp)
		// 被删除之前的value是什么
		if len(delResp.PrevKvs) != 0 {
			for _, kvpair := range delResp.PrevKvs {
				fmt.Println("删除了:", string(kvpair.Key), string(kvpair.Value))
			}
		}
	}

}

//-----------------------------------------监听----------------------------------------------
//-----------------------------------------租约----------------------------------------------
//-----------------------------------------分布式锁----------------------------------------------