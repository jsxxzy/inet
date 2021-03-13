package logging_test

import (
	"testing"

	"github.com/jsxxzy/inet/cmd/inetdaemon/logging"
)

func TestGetDayFormat(t *testing.T) {
	logging.New("runlog")
	t.Log("文件测试成功")
}
