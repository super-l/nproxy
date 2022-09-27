package task

import (
	"bufio"
	"fmt"
	"github.com/super-l/nproxy/internal"
	"github.com/super-l/nproxy/internal/config"
	"github.com/super-l/nproxy/pkg/suphttp"
	"github.com/super-l/nproxy/services/model"
	"github.com/super-l/nproxy/utils"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

type proxyTaskService struct{}

var ProxyTaskSercice = proxyTaskService{}

func (p proxyTaskService) Start() {
	var ticker = time.NewTicker(time.Duration(config.GetConfig().Proxy.ApiRate) * time.Second)
	defer ticker.Stop()

	internal.SLogger.GetStdoutLogger().Info(fmt.Sprintf("the dynamic proxy api data synchronization thread starts!"))

	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			proxyApiListData, err := model.MProxyApi.List()
			if err == nil && len(proxyApiListData) > 0 {
				for _, data := range proxyApiListData {
					wg.Add(1)
					go p.doItem(data, &wg)
				}
				wg.Wait()
				internal.SLogger.GetStdoutLogger().Info("all api proxy data capture completed!")
			}
		}
	}
}

func (p proxyTaskService) doItem(proxyApiData model.ProxyApi, wg *sync.WaitGroup) {
	defer wg.Done()

	apiUrl := proxyApiData.Value
	if apiUrl == "" || !utils.Url.IsUrl(apiUrl) {
		internal.SLogger.Warn(fmt.Sprintf("api url is error, id:%d data: %s", proxyApiData.ID, apiUrl))
		return
	}

	// 获取内容
	var pageHtmlContent string
	resultData, err := suphttp.New().DiyGet(apiUrl, 30, "nproxy", "", "", 0, nil)
	if err == nil && resultData.StatusCode == 200 && resultData.Body != "" {
		pageHtmlContent = resultData.Body
	} else {
		internal.SLogger.GetStdoutLogger().Errorf("the proxy api service interface is invalid! url: %s  tips: %s", apiUrl, err.Error())
		return
	}

	// 验证内容
	var r io.Reader = strings.NewReader(pageHtmlContent)
	buff := bufio.NewReader(r)
	var count = 0
	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		proxy := string(data)
		if proxy != "" {
			var proxyUri string
			if proxyApiData.ProtocolType != "" {
				proxyUri = proxyApiData.ProtocolType + "://" + proxy
			}

			u, parseErr := url.Parse(proxyUri)
			if parseErr != nil {
				internal.SLogger.StdoutLogger.Warnf("failed to parse url data: %s tips: %s", proxyUri, parseErr.Error())
				continue
			}

			ip := net.ParseIP(u.Host)

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

			// 计算到期时间
			expiredAtUnix := time.Now().UnixNano() + proxyApiData.PeriodValidity*1000
			expiredAt := time.Unix(expiredAtUnix, 0)

			proxyData := model.Proxy{
				ProtocolType: proxyApiData.ProtocolType,
				LineType:     proxyApiData.LineType,
				Value:        proxy,
				Country:      addr,
				Source:       fmt.Sprintf("api-%d", proxyApiData.ID),
				UsedTimes:    0,
				ExpiredAt:    expiredAt,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			}
			model.MProxy.Add(proxyData)
			count++
		}
	}
	internal.SLogger.GetStdoutLogger().Infof("api proxy data(%d) capture completed! results the number of: %d", proxyApiData.ID, count)
}
