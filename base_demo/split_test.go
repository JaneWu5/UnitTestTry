package base_demo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
func TestSplit(t *testing.T) {
	// 程序输出的结果
	got := Split("a:b:c", ":")
	// 期望的结果
	want := []string{"a", "b", "c"}
	// 因为slice不能比较直接，借助反射包中的方法比较
	//if !reflect.DeepEqual(want, got) {
	//	// 测试失败输出错误提示
	//	t.Errorf("expected:%v, got:%v", want, got)
	//}
	// 使用testify断言
	assert.Equal(t, got, want)
}

func TestSplitWithComplexSep(t *testing.T) {
	got := Split("abcd", "bc")
	want := []string{"a", "d"}
	//if !reflect.DeepEqual(want, got) {
	//	t.Errorf("expected:%v, got:%v", want, got)
	//}
	// 使用testify断言
	assert.Equal(t, got, want)
}

// 我们修改了代码之后仅仅执行那些失败的测试用例或新引入的测试用例是错误且危险的，
// 正确的做法应该是完整运行所有的测试用例，保证不会因为修改代码而引入新的问题。
// go test -v 运行全部单元测试
// 我们可以看到，有了单元测试就能够在代码改动后快速进行回归测试，极大地提高开发效率并保证代码的质量。

// 当执行"go test -short"时就不会执行TestTimeConsuming
func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("short模式下会跳过该测试用例")
	} else {
		TestSplit(t)
	}
}

// 上面我们为每一个测试数据编写了一个测试函数，
// 而通常单元测试中需要多组测试数据保证测试的效果。
// Go1.7+中新增了子测试，支持在测试函数中使用t.Run执行一组测试用例，
// 这样就不需要为不同的测试数据定义多个测试函数了。
func TestAGroup(t *testing.T) {
	t.Run("case1", TestSplit)
	t.Run("case2", TestSplitWithComplexSep)
	t.Run("case3", TestTimeConsuming)
}

func TestSplitAll(t *testing.T) {
	// 定义测试表格
	// 这里使用匿名结构体定义了若干个测试用例
	// 并且为每个测试用例设置了一个名称
	tests := []struct {
		name  string
		input string
		sep   string
		want  []string
	}{
		{"base case", "a:b:c", ":", []string{"a", "b", "c"}},
		{"wrong sep", "a:b:c", ",", []string{"a:b:c"}},
		{"more sep", "abcd", "bc", []string{"a", "d"}},
		{"leading sep", "沙河有沙又有河", "沙", []string{"", "河有", "又有河"}},
	}
	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // 使用t.Run()执行子测试
			got := Split(tt.input, tt.sep)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("expected:%#v, got:%#v", tt.want, got)
			//}
			// 使用testify断言
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestParallelSplitAll(t *testing.T) {
	// 将 TLog 标记为能够与其他测试并行运行
	t.Parallel()
	// 定义测试表格
	// 这里使用匿名结构体定义了若干个测试用例
	// 并且为每个测试用例设置了一个名称
	tests := []struct {
		name  string
		input string
		sep   string
		want  []string
	}{
		{"base case", "a:b:c", ":", []string{"a", "b", "c"}},
		{"wrong sep", "a:b:c", ",", []string{"a:b:c"}},
		{"more sep", "abcd", "bc", []string{"a", "d"}},
		{"leading sep", "沙河有沙又有河", "沙", []string{"", "河有", "又有河"}},
	}
	// 遍历测试用例
	for _, tt := range tests {
		// 注意这里重新声明tt变量（避免多个goroutine中使用了相同的变量）
		tt := tt
		t.Run(tt.name, func(t *testing.T) { // 使用t.Run()执行子测试
			t.Parallel() // 将每个测试用例标记为能够彼此并行运行
			got := Split(tt.input, tt.sep)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("expected:%#v, got:%#v", tt.want, got)
			//}
			// 使用testify断言
			assert.Equal(t, got, tt.want)
		})
	}
}

// go test -cover -coverprofile=c.out 将覆盖率报告输出到一个叫做c.out的文件中
// 执行"go tool cover -html=c.out"，使用cover工具来处理生成的记录信息，该命令会打开本地的浏览器窗口生成一个HTML报告

// testify是一个社区非常流行的Go单元测试工具包，其中使用最多的功能就是它提供的断言工具——testify/assert或testify/require
// go get github.com/stretchr/testify

// go test -run=TestSomething -v
// 当我们有多个断言语句时，还可以使用assert := assert.New(t) 创建一个assert对象，
// 它拥有前面所有的断言方法，只是不需要再传入Testing.T参数了。
func TestSomething(t *testing.T) {
	assert := assert.New(t)
	// assert equality
	assert.Equal(123, 123, "they should be equal")
	// assert inequality
	assert.NotEqual(123, 456, "they should not be equal")
	obj := struct {
		Value string
	}{
		Value: "Something",
	}
	// assert for nil (good for errors)
	assert.Nil(nil)
	// assert for not nil (good when you expect something)
	if assert.NotNil(obj) {
		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.Equal("Something", obj.Value)
	}
}
