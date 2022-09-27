package utils

import (
	"github.com/asaskevich/govalidator"
	"strings"
)

type uUrl struct{}

var Url = uUrl{}

func (uUrl) IsUrl(str string) bool {
	if !strings.HasPrefix(str, "http") {
		return false
	}
	return govalidator.IsURL(str)
}

func (uUrl) GetIp(str string) string {
	if !strings.Contains(str, ":") {
		return str
	}
	index := strings.Index(str, ":")
	return str[:index]
}
