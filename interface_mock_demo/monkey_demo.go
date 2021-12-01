package interface_mock_demo

import (
	"UnitTestTry/varys"
	"fmt"
)

// monkey是一个Go单元测试中十分常用的打桩工具，它在运行时通过汇编语言重写可执行文件，将目标函数或方法的实现跳转到桩实现，其原理类似于热补丁。
// monkey库很强大，但是使用时需注意以下事项：
// 1. monkey不支持内联函数，在测试的时候需要通过命令行参数-gcflags=-l关闭Go语言的内联优化。
// 2. monkey不是线程安全的，所以不要把它用到并发的单元测试中。

// 假设你们公司中台提供了一个用户中心的库varys，使用这个库可以很方便的根据uid获取用户相关信息。
// 但是当你编写代码的时候这个库还没实现，或者这个库要经过内网请求但你现在没这能力，
// 这个时候要为MyFunc编写单元测试，就需要做一些mock工作。
func MyFunc(uid int64) string {
	u, err := varys.GetInfoByUID(uid)
	if err != nil {
		return "welcome"
	}
	// 这里是一些逻辑代码...
	return fmt.Sprintf("hello %s\n", u.Name)
}
