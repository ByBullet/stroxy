package test

import (
	"path"
	"runtime"
	"testing"

	"github.com/ByBullet/stroxy/util"
)

func TestGetCurrentAbPath(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	contextPath := path.Dir(path.Dir(file))
	if contextPath != util.GetCurrentAbPath() {
		t.Errorf("GetCurrentAbPath测试未通过,GetCurrentAbPath返回%s", util.GetCurrentAbPath())
	}
}
