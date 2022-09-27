package utils

import (
	"bufio"
	"github.com/axgle/mahonia"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type uFile struct{}

var File = uFile{}

func (uFile) FileOrPathIsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (uFile) CreateDir(string2 string) error {
	err := os.MkdirAll(string2, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (u uFile) TransContentUtf8(content string) string {
	if content == "" {
		return ""
	}
	var value = strings.TrimSpace(content)
	if !u.IsUtf8([]byte(content)) {
		value = u.ConvertToString(value, "gbk")
	}
	if value != "" {
		return value
	}
	return ""
}

func (u uFile) ReadTxtFileError(filepath string) ([]string, error) {
	var result []string
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	buff := bufio.NewReader(fd)
	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}

		var value = strings.TrimSpace(string(data))
		if !u.IsUtf8(data) {
			value = u.ConvertToString(value, "gbk")
		}

		if value != "" {
			result = append(result, value)
		}
	}
	return result, nil
}

func (uFile) ConvertToString(src string, srcCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	return srcResult
}

func (u uFile) IsUtf8(data []byte) bool {
	return utf8.Valid(data)
}

func (u uFile) WriteFileAppend(filepath string, contents []string) error {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, content := range contents {
		_, _ = f.WriteString(content + "\r\n")
	}
	return nil
}

func (u uFile) DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	return err
}
