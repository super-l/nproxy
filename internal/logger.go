package internal

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/super-l/nproxy/internal/config"
	"github.com/super-l/nproxy/internal/consts"
	"os"
	"strings"
	"time"
)

type sLogger struct {
	Level        logrus.Level
	StdoutLogger *logrus.Logger
	FileLogger   *logrus.Logger
}

var SLogger = sLogger{}

func (s *sLogger) InitLogger() error {
	level, err := logrus.ParseLevel(config.GetConfig().System.LogLevel)
	if err != nil {
		return err
	}
	s.Level = level

	s.initStdoutLogger()
	s.initFileLogger()

	return nil
}

//Custom log format definition
type MyFormatter struct{}

func (s *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func (s *sLogger) initStdoutLogger() {
	s.StdoutLogger = logrus.New()

	s.StdoutLogger.Out = os.Stdout

	s.StdoutLogger.SetLevel(s.Level)

	//自定writer就行， hook 交给 lfshook
	s.StdoutLogger.SetFormatter(new(MyFormatter))
}

func (s *sLogger) initFileLogger() {
	s.FileLogger = logrus.New()
	s.FileLogger.SetLevel(s.Level)
	s.FileLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "01-02 15:04:05",
	})

	writer, _ := rotatelogs.New(
		consts.LogDirPath+"%Y%m%d%H"+"-"+consts.LogFileName,
		rotatelogs.WithLinkName(consts.LogDirPath+consts.LogFileName),
		rotatelogs.WithMaxAge(30*24*time.Hour),                     // 一个月
		rotatelogs.WithRotationTime(time.Duration(60)*time.Minute), // 按分钟
	)
	s.FileLogger.SetOutput(writer)
	//writeMap := lfshook.WriterMap{
	//	logrus.InfoLevel:  writer,
	//	logrus.FatalLevel: writer,
	//	logrus.DebugLevel: writer,
	//	logrus.WarnLevel: writer,
	//	logrus.ErrorLevel: writer,
	//	logrus.PanicLevel: writer,
	//}
	//lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
	//	TimestampFormat:"2006-01-02 15:04:05",
	//})
	//fileLogger.AddHook(lfHook)

	//如果你希望将调用的函数名添加为字段，请通过以下方式设置： 开启这个模式会增加性能开销。
	//fileLogger.SetReportCaller(true)
}

func (s *sLogger) GetLogger() *logrus.Logger {
	return s.GetFileLogger()
}

func (s *sLogger) GetStdoutLogger() *logrus.Logger {
	return s.StdoutLogger
}

func (s *sLogger) GetFileLogger() *logrus.Logger {
	return s.FileLogger
}

func (s *sLogger) Warn(content string) {
	s.GetStdoutLogger().Warn(content)
	s.GetFileLogger().Warn(content)
}

func (s *sLogger) Error(content string) {
	s.GetStdoutLogger().Error(content)
	s.GetFileLogger().Error(content)
}
