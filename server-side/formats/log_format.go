package formats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct {
}

func (formatter *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 自定义日志格式，包括时间、级别和消息
	var message string
	// 提取所有的 key
	keys := make([]string, 0, len(entry.Data))
	for key := range entry.Data {
		keys = append(keys, key)
	}
	// 对 key 进行排序
	sort.Strings(keys)
	// 按排序后的 key 顺序拼接数据
	var dataParts []string
	for _, key := range keys {
		value := entry.Data[key]
		dataParts = append(dataParts, fmt.Sprintf("%s=%v", key, value))
	}
	dataParts = append(dataParts, entry.Message)
	message = strings.Join(dataParts, " ")
	return fmt.Appendf(nil, "[%s] [%s] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Level.String(),
		message,
	), nil
}
