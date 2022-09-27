package consts

import (
	"path/filepath"
)

const (
	LogDirPath  = "data" + string(filepath.Separator) + "logs" + string(filepath.Separator)
	LogFileName = "nproxy.log"
	Version     = "1.0.0" // 当前版本
	Name        = "NProxy"
)
