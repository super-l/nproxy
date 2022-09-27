package suphttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type supClient struct {
	url         string
	header      http.Header
	client      *http.Client
	transport   *http.Transport
	allowType   []string
	autoUtf8    bool
	maxBodySize int64
	postData    map[string]interface{}
}

type ResponseData struct {
	Body       string
	Status     string
	StatusCode int
	Scheme     string
	Server     string
	XPoweredBy string
}

func New() *supClient {
	c := &supClient{
		header: http.Header{},
		client: &http.Client{Timeout: 30 * time.Second},
		transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			ResponseHeaderTimeout: 35 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		},
		allowType:   []string{},
		autoUtf8:    false,
		maxBodySize: 0,
	}
	return c
}

func (c *supClient) SetUrl(url string) *supClient {
	c.url = url
	return c
}

func (c *supClient) SetTimeOut(timeout uint) *supClient {
	c.client.Timeout = time.Duration(timeout) * time.Second
	return c
}

func (c *supClient) SetMaxBodySize(size int64) *supClient {
	c.maxBodySize = size
	return c
}

func (c *supClient) SetHeader(param string, value string) *supClient {
	c.header.Set(param, value)
	return c
}

func (c *supClient) SetTLSClientConfig(config *tls.Config) *supClient {
	c.transport.TLSClientConfig = config
	return c
}

func (c *supClient) AddHeader(param string, value string) *supClient {
	c.header.Add(param, value)
	return c
}

func (c *supClient) SetProxy(proxyUrl string) (*supClient, error) {
	parsedProxyUrl, err := url.Parse(proxyUrl)
	if err == nil {
		if proxyUrl == "" {
			c.transport.Proxy = nil
		} else {
			c.transport.Proxy = http.ProxyURL(parsedProxyUrl)
		}
	} else {
		return c, err
	}
	return c, nil
}

func (c *supClient) AddAllowedType(contentType string) *supClient {
	c.allowType = append(c.allowType, contentType)
	return c
}

func (c *supClient) checkAllowedType(contentType string) bool {
	if len(c.allowType) == 0 {
		return true
	}
	for _, value := range c.allowType {
		if strings.Contains(contentType, value) {
			return true
		}
	}
	return false
}

func (c *supClient) SetAutoUtf8(enable bool) *supClient {
	c.autoUtf8 = enable
	return c
}

func (c *supClient) makeRequest(method string, url string) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)
	if req, err = http.NewRequest(strings.ToUpper(method), url, nil); err != nil {
		return nil, err
	}
	for k, values := range c.header {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}
	return req, err
}

func (c *supClient) DiyGet(url string, timeOut uint, agent string, referer string, proxy string, maxSize int64, contentType []string) (ResponseData, error) {
	if url == "" {
		return ResponseData{}, errors.New("url is error")
	}
	c.SetUrl(url)

	if timeOut != 0 {
		c.SetTimeOut(timeOut)
	} else {
		c.SetTimeOut(30)
	}

	if agent != "" {
		c.SetHeader("User-Agent", agent)
	}

	if referer != "" {
		c.SetHeader("referer", referer)
	}

	if proxy != "" {
		c.SetProxy(proxy)
	}

	if maxSize != 0 {
		c.SetMaxBodySize(maxSize)
	}

	if contentType != nil && len(contentType) != 0 {
		for _, value := range contentType {
			c.AddAllowedType(value)
		}
	}
	return c.sysGet()
}

func (c *supClient) Get(url string) (string, error) {
	if url == "" {
		return "", errors.New("url is error")
	}
	c.SetUrl(url)
	result, err := c.sysGet()
	if err != nil {
		return "", err
	}
	return result.Body, nil
}

func (c *supClient) sysGet() (ResponseData, error) {
	req, err := c.makeRequest(http.MethodGet, c.url)
	if err != nil {
		return ResponseData{}, err
	}
	c.client.Transport = c.transport
	resp, err := c.client.Do(req)
	if err != nil {
		return ResponseData{}, err
	}
	defer resp.Body.Close()

	var bodyStr string
	buf := new(bytes.Buffer)

	if c.maxBodySize == 0 {
		_, err = io.Copy(buf, resp.Body)
	} else {
		_, err = io.Copy(buf, io.LimitReader(resp.Body, c.maxBodySize))
		if err != nil {
			return ResponseData{}, err
		}
	}
	bodyStr = buf.String()
	result := ResponseData{
		Body:       bodyStr,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Scheme:     resp.Request.URL.Scheme,
		Server:     resp.Header.Get("server"),
		XPoweredBy: resp.Header.Get("X-Powered-By"),
	}
	return result, nil
}

func (c *supClient) PostJson(data map[string]interface{}) (string, error) {
	jsonStr, _ := json.Marshal(data)
	resp, err := c.client.Post(c.url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

func (c *supClient) PostForm(data url.Values) (string, error) {
	resp, err := c.client.PostForm(c.url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
