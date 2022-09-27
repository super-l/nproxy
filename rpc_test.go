package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"testing"
	"time"
)

var proxyServiceClient client.XClient // 任务引擎链接
var d *client.Peer2PeerDiscovery

func initXClient() error {
	var err error
	d, err = client.NewPeer2PeerDiscovery("tcp@127.0.0.1:55555", "")
	if err != nil {
		return err
	}

	// 授权密码
	rpcPassword := "tGzv3JOkF0XG5Qx2TlKWIA"

	// 配置选项
	option := client.DefaultOption
	option.IdleTimeout = 30 * time.Second
	option.ConnectTimeout = 60 * time.Second

	// XClient 是对客户端的封装，增加了一些服务发现和服务治理的特性。
	proxyServiceClient = client.NewXClient("Proxy", client.Failtry, client.RandomSelect, d, option)
	if rpcPassword != "" {
		proxyServiceClient.Auth(rpcPassword)
	}
	return nil
}

func CloseProxyClient() {
	_ = proxyServiceClient.Close()
}

type ProxyReply struct {
	Code int
	Data interface{}
	Msg  string
}

type GetProxyArgs struct {
	ProtocolType string // 协议类型
	LineType     int    // 线路类型
	Country      string // 国家
	Count        int    // 个数
}

type AddProxyArgs struct {
	ProtocolType string `json:"protocol_type"` // 协议类型
	LineType     int    `json:"line_type"`     // 线路类型
	Value        string `json:"value"`         // 值
	Source       string `json:"source"`        // 来源
}

func RpcProxyGet() (*ProxyReply, error) {
	if proxyServiceClient == nil {
		return nil, errors.New("rpc server connect faild!")
	}
	reply := &ProxyReply{}
	args := &GetProxyArgs{}
	err := proxyServiceClient.Call(context.Background(), "GetProxy", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func RpcProxyList() (*ProxyReply, error) {
	if proxyServiceClient == nil {
		return nil, errors.New("rpc server connect faild!")
	}
	reply := &ProxyReply{}
	args := &GetProxyArgs{}
	err := proxyServiceClient.Call(context.Background(), "ListProxy", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func RpcProxyAdd() (*ProxyReply, error) {
	if proxyServiceClient == nil {
		return nil, errors.New("rpc server connect faild!")
	}
	reply := &ProxyReply{}
	args := &AddProxyArgs{
		ProtocolType: "http",
		LineType:     2,
		Value:        "103.96.149.195:10809",
		Source:       "manual",
	}
	err := proxyServiceClient.Call(context.Background(), "AddProxy", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func RpcProxyApiList() (*ProxyReply, error) {
	if proxyServiceClient == nil {
		return nil, errors.New("rpc server connect faild!")
	}
	reply := &ProxyReply{}
	args := &GetProxyArgs{}
	err := proxyServiceClient.Call(context.Background(), "ListProxyApi", args, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func TestRpc(t *testing.T) {
	// 链接RPC服务器
	initXClient()

	// 获取一个最佳代理
	reply, err := RpcProxyGet()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(fmt.Sprintf("Code: %d Msg: %s Data: %v", reply.Code, reply.Msg, reply.Data))

	// 添加代理
	reply2, err2 := RpcProxyAdd()
	if err2 != nil {
		fmt.Println(err2.Error())
	}
	fmt.Println(fmt.Sprintf("Code: %d Msg: %s Data: %v", reply2.Code, reply2.Msg, reply2.Data))

	// 获取代理列表
	reply3, err3 := RpcProxyList()
	if err3 != nil {
		fmt.Println(err3.Error())
	}
	fmt.Println(fmt.Sprintf("Code: %d Msg: %s Data: %v", reply3.Code, reply3.Msg, reply3.Data))

	// 获取代理API列表
	reply4, err4 := RpcProxyApiList()
	if err4 != nil {
		fmt.Println(err4.Error())
	}
	fmt.Println(fmt.Sprintf("Code: %d Msg: %s Data: %v", reply4.Code, reply4.Msg, reply4.Data))

	// 关闭链接
	CloseProxyClient()
}
