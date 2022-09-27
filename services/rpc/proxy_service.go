package rpc

import (
	"context"
	"fmt"
	"github.com/super-l/nproxy/internal"
	"github.com/super-l/nproxy/services"
	"github.com/super-l/nproxy/services/model"
	"github.com/super-l/nproxy/services/rpc/bean"
	"github.com/super-l/nproxy/utils"
	"net"
	"time"
)

type ProxyService int

func (ProxyService) GetProxy(ctx context.Context, args bean.GetProxyArgs, reply *bean.ProxyReply) error {
	result, err := model.MProxy.Get(args)
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	var values []string
	for _, proxy := range result {
		proxyData := fmt.Sprintf("%s://%s", proxy.ProtocolType, proxy.Value)
		values = append(values, proxyData)
	}

	if len(values) == 0 {
		reply.Error(services.NoDataMsg.Code, services.NoDataMsg.Msg)
		return nil
	}
	reply.Success(200, values, "success")
	return nil
}

func (ProxyService) AddProxy(ctx context.Context, args bean.AddProxyArgs, reply *bean.ProxyReply) error {
	data, err := model.MProxy.GetByValue(args.Value)
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}

	// 已经存在数据
	if !model.MProxy.IsEmpty(data) {
		internal.SLogger.StdoutLogger.Warnf("this is a duplicate data: %s", data.Value)
		reply.Error(services.RepeatMsg.Code, services.RepeatMsg.Msg)
		return nil
	}

	var proxyData model.Proxy
	proxyData.CreatedAt = time.Now()
	proxyData.UsedTimes = 0
	proxyData.Source = args.Source
	proxyData.LineType = args.LineType
	proxyData.ProtocolType = args.ProtocolType
	proxyData.Value = args.Value

	ip := net.ParseIP(utils.Url.GetIp(proxyData.Value))

	// Country of the computing server
	var addr string
	record, errReadCountry := internal.IpDb.GetIpDbInstance().Country(ip)
	if errReadCountry != nil {
		addr = ""
	} else {
		// Adhere to the one China principle
		addr = record.Country.Names["zh-CN"]
		if addr == "中华民国" {
			addr = "中国台湾"
		}
	}
	proxyData.Country = addr

	result, addErr := model.MProxy.Add(proxyData)
	if addErr != nil {
		internal.SLogger.StdoutLogger.Warn(addErr.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, result, "success")
	return nil
}

func (ProxyService) ListProxy(ctx context.Context, args bean.ListProxyArgs, reply *bean.ProxyReply) error {
	listData, err := model.MProxy.List()
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	if len(listData) == 0 {
		reply.Error(services.NoDataMsg.Code, services.NoDataMsg.Msg)
		return nil
	}
	reply.Success(200, listData, "success")
	return nil
}

func (ProxyService) UpdateProxy(ctx context.Context, args bean.UpdateProxyArgs, reply *bean.ProxyReply) error {
	proxyData, err := model.MProxy.GetById(args.Id)
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	if model.MProxy.IsEmpty(proxyData) {
		internal.SLogger.StdoutLogger.Warnf("non-existent data: %s", proxyData.Value)
		reply.Error(services.NoDataMsg.Code, services.NoDataMsg.Msg)
		return nil
	}

	proxyData.UpdatedAt = time.Now()
	proxyData.Value = args.Value
	proxyData.LineType = args.LineType
	proxyData.ProtocolType = args.ProtocolType
	proxyData.Source = args.Source
	updateErr := model.MProxy.Update(proxyData)
	if err != nil {
		internal.SLogger.StdoutLogger.Warn(updateErr.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, proxyData, "success")
	return nil
}

func (ProxyService) DeleteProxy(ctx context.Context, args bean.DeleteProxyArgs, reply *bean.ProxyReply) error {
	affect := model.MProxy.DeleteMore(args.IdList)
	if affect == 0 {
		internal.SLogger.Warn("delete proxy data failed")
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, nil, "success")
	return nil
}

func (ProxyService) AddProxyApi(ctx context.Context, args bean.AddProxyApiArgs, reply *bean.ProxyReply) error {
	// Duplicate data judgment
	data, err := model.MProxyApi.GetByValue(args.Value)
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	if !model.MProxyApi.IsEmpty(data) {
		internal.SLogger.StdoutLogger.Warnf("this is a duplicate data: %s", data.Value)
		reply.Error(services.RepeatMsg.Code, services.RepeatMsg.Msg)
		return nil
	}

	var proxyApiData model.ProxyApi
	proxyApiData.CreatedAt = time.Now()
	proxyApiData.GetTimes = 0
	proxyApiData.LineType = args.LineType
	proxyApiData.ProtocolType = args.ProtocolType
	proxyApiData.Value = args.Value

	result, addErr := model.MProxyApi.Add(proxyApiData)

	if addErr != nil {
		internal.SLogger.StdoutLogger.Warn(addErr.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, result, "success")
	return nil
}

func (ProxyService) ListProxyApi(ctx context.Context, args bean.ListProxyApiArgs, reply *bean.ProxyReply) error {
	listData, err := model.MProxyApi.List()
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, listData, "success")
	return nil
}

func (ProxyService) UpdateProxyApi(ctx context.Context, args bean.UpdateProxyApiArgs, reply *bean.ProxyReply) error {
	proxyApiData, err := model.MProxyApi.GetById(args.Id)
	if err != nil {
		internal.SLogger.Warn(err.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	if model.MProxyApi.IsEmpty(proxyApiData) {
		internal.SLogger.StdoutLogger.Warnf("non-existent data: %s", proxyApiData.Value)
		reply.Error(services.NoDataMsg.Code, services.NoDataMsg.Msg)
		return nil
	}

	proxyApiData.UpdatedAt = time.Now()
	proxyApiData.Value = args.Value
	proxyApiData.LineType = args.LineType
	proxyApiData.ProtocolType = args.ProtocolType
	updateErr := model.MProxyApi.Update(proxyApiData)
	if err != nil {
		internal.SLogger.StdoutLogger.Warn(updateErr.Error())
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, proxyApiData, "success")
	return nil
}

func (ProxyService) DeleteProxyApi(ctx context.Context, args bean.DeleteProxyApiArgs, reply *bean.ProxyReply) error {
	affect := model.MProxyApi.DeleteMore(args.IdList)
	if affect == 0 {
		internal.SLogger.Warn("delete proxy api data failed")
		reply.Error(services.CommonErrMsg.Code, services.CommonErrMsg.Msg)
		return nil
	}
	reply.Success(200, nil, "success")
	return nil
}
