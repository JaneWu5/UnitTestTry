package interface_mock_demo

import "io/ioutil"

// GoStub也是一个单元测试中的打桩工具，它支持为全局变量、函数等打桩。
// 不过我个人感觉它为函数打桩不太方便，我一般在单元测试中只会使用它来为全局变量打桩。
// go get github.com/prashantv/gostub

// 官方为全局变量打桩
var (
	configFile = "config.json"
	maxNum     = 10
)

func GetConfig() ([]byte, error) {
	return ioutil.ReadFile(configFile)
}

func ShowNumber() int {
	// ...
	return maxNum
}
