package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/super-l/nproxy/internal"
	"github.com/super-l/nproxy/internal/config"
	"net"
	"os"
)

var clientConn net.Conn
var connected = false

func InitRpcServer() {
	s := server.NewServer()
	s.RegisterName("Proxy", new(ProxyService), "")

	// 如果密码不为空
	if config.GetConfig().Rpc.Password != "" {
		s.AuthFunc = auth
	}
	address := fmt.Sprintf("%s:%s", config.GetConfig().Rpc.Ip, config.GetConfig().Rpc.Port)
	err := s.Serve("tcp", address)
	if err != nil {
		internal.SLogger.StdoutLogger.Error(err.Error())
		os.Exit(1)
	}
}

func auth(ctx context.Context, req *protocol.Message, token string) error {
	if token == config.GetConfig().Rpc.Password {
		return nil
	}
	return errors.New("invalid token")
}
