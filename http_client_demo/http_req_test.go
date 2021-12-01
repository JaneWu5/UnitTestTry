package http_client_demo

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

// 如果不想在测试过程中真正去发送请求或者依赖的外部接口还没有开发完成时，
// 我们可以在单元测试中对依赖的API进行mock。
// 推荐库gock: go get -u gopkg.in/h2non/gock.v1

func TestGetResultByAPI(t *testing.T) {
	defer gock.Off() // 测试执行后刷新挂起的mock
	// mock 请求外部api时传参x=1返回100
	gock.New("http://your-api.com").
		Post("/post").
		MatchType("json").
		JSON(map[string]int{"x": 1}).
		Reply(200).
		JSON(map[string]int{"value": 100})

	// 调用我们的业务函数
	res := GetResultByAPI(1, 1)

	// 校验返回结果是否符合预期
	assert.Equal(t, res, 101)

	// mock 请求外部api时传参x=2返回200
	gock.New("http://your-api.com").
		Post("/post").
		MatchType("json").
		JSON(map[string]int{"x": 2}).
		Reply(200).
		JSON(map[string]int{"value": 200})

	// 调用我们的业务函数
	res = GetResultByAPI(2, 2)
	// 校验返回结果是否符合预期
	assert.Equal(t, res, 202)

	assert.True(t, gock.IsDone()) // 断言mock被触发
}
