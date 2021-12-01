package interface_mock_demo

import (
	"UnitTestTry/varys"
	"bou.ke/monkey"
	"reflect"
	"strings"
	"testing"
)

// 我们使用monkey库对varys.GetInfoByUID进行打桩。

func TestMyFunc(t *testing.T) {
	// 无论传入的uid是多少，都返回 &varys.UserInfo{Name: "liwenzhou"}, nil
	monkey.Patch(varys.GetInfoByUID, func(int64) (*varys.UserInfo, error) {
		return &varys.UserInfo{Name: "wz"}, nil
	})
	ret := MyFunc(123)
	if !strings.Contains(ret, "wz") {
		t.Fatal("返回姓名不正确")
	}
}

// 这里为防止内联优化添加了-gcflags=-l参数:
// go test -run=TestMyFunc -v -gcflags=-l

func TestUser_GetInfo(t *testing.T) {
	var u = &varys.UserInfo{
		Name:     "wz",
		Birthday: "1989-09-22",
	}
	// 为对象方法打桩
	monkey.PatchInstanceMethod(reflect.TypeOf(u), "CalcAge", func(info *varys.UserInfo) int {
		return 18
	})
	ret := u.GetInfo() // 内部调用u.CalcAge方法时会返回18
	t.Log(ret)
	if !strings.Contains(ret, "朋友") {
		t.Fatal()
	}
}

// 社区中还有一个参考monkey库实现的gomonkey库，原理和使用过程基本相似，这里就不再啰嗦了。
// 除此之外社区里还有一些其他打桩工具如GoStub（上一篇介绍过为全局变量打桩）等。
// 熟练使用各种打桩工具能够让我们更快速地编写合格的单元测试，为我们的软件保驾护航。
