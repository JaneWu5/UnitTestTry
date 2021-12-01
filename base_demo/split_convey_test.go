package base_demo

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

// GoConvey是一个非常非常好用的Go测试框架，它直接与go test集成，提供了很多丰富的断言函数，
// 能够在终端输出可读的彩色测试结果，并且还支持全自动的Web UI。

func TestSplitConvey(t *testing.T) {
	convey.Convey("基础用例", t, func() {
		var (
			s      = "a:b:c"
			sep    = ":"
			expect = []string{"a", "b", "c"}
		)
		got := Split(s, sep)
		convey.So(got, convey.ShouldResemble, expect) // 断言
	})
	convey.Convey("不包含分隔符用例", t, func() {
		var (
			s      = "a:b:c"
			sep    = "|"
			expect = []string{"a:b:c"}
		)
		got := Split(s, sep)
		convey.So(got, convey.ShouldResemble, expect) // 断言
	})
}

// goConvey还支持在单元测试中根据需要嵌套调用，比如：
func TestSplitEmbed(t *testing.T) {
	// 只需要在顶层的Convey调用时传入t
	convey.Convey("分隔符在开头或结尾用例", t, func() {
		tt := []struct {
			name   string
			s      string
			sep    string
			expect []string
		}{
			{"分隔符在开头", "*1*2*3", "*", []string{"", "1", "2", "3"}},
			{"分隔符在结尾", "1+2+3+", "+", []string{"1", "2", "3", ""}},
		}
		for _, tc := range tt {
			// 嵌套调用Convey
			convey.Convey(tc.name, func() {
				got := Split(tc.s, tc.sep)
				convey.So(got, convey.ShouldResemble, tc.expect)
			})
		}
	})
}

// 断言方法：

// 1. 一般相等类
// So(thing1, ShouldEqual, thing2)
// So(thing1, ShouldNotEqual, thing2)

// 用于数组、切片、map和结构体相等
// So(thing1, ShouldResemble, thing2)

// So(thing1, ShouldNotResemble, thing2)
// So(thing1, ShouldPointTo, thing2)
// So(thing1, ShouldNotPointTo, thing2)
// So(thing1, ShouldBeNil)
// So(thing1, ShouldNotBeNil)
// So(thing1, ShouldBeTrue)
// So(thing1, ShouldBeFalse)
// So(thing1, ShouldBeZeroValue)

// 2. 数字数量比较类
//So(1, ShouldBeGreaterThan, 0)
//So(1, ShouldBeGreaterThanOrEqualTo, 0)
//So(1, ShouldBeLessThan, 2)
//So(1, ShouldBeLessThanOrEqualTo, 2)
//So(1.1, ShouldBeBetween, .8, 1.2)
//So(1.1, ShouldNotBeBetween, 2, 3)
//So(1.1, ShouldBeBetweenOrEqual, .9, 1.1)
//So(1.1, ShouldNotBeBetweenOrEqual, 1000, 2000)

// tolerance is optional; default 0.0000000001（十位小数）
//So(1.0, ShouldAlmostEqual, 0.99999999, .0001)

//So(1.0, ShouldNotAlmostEqual, 0.9, .0001)

// 3. 包含类
// So([]int{2, 4, 6}, ShouldContain, 4)
// So([]int{2, 4, 6}, ShouldNotContain, 5)
// So(4, ShouldBeIn, ...[]int{2, 4, 6})
// So(4, ShouldNotBeIn, ...[]int{1, 3, 5})
// So([]int{}, ShouldBeEmpty)
// So([]int{1}, ShouldNotBeEmpty)
// So(map[string]string{"a": "b"}, ShouldContainKey, "a")
// So(map[string]string{"a": "b"}, ShouldNotContainKey, "b")
// So(map[string]string{"a": "b"}, ShouldNotBeEmpty)
// So(map[string]string{}, ShouldBeEmpty)

// supports map, slice, chan, and string
// So(map[string]string{"a": "b"}, ShouldHaveLength, 1)

// 4. 字符串类
// So("asdf", ShouldStartWith, "as")
// So("asdf", ShouldNotStartWith, "df")
// So("asdf", ShouldEndWith, "df")
// So("asdf", ShouldNotEndWith, "df")

// optional 'expected occurences' arguments?
// So("asdf", ShouldContainSubstring, "稍等一下")

// So("asdf", ShouldNotContainSubstring, "er")
// So("adsf", ShouldBeBlank)
// So("asdf", ShouldNotBeBlank)

// 5. panic类
// So(func(), ShouldPanic)
// So(func(), ShouldNotPanic)

// or errors.New("something")
// So(func(), ShouldPanicWith, "")

// or errors.New("something")
// So(func(), ShouldNotPanicWith, "")

// 6. 类型检查类
// So(1, ShouldHaveSameTypeAs, 0)
// So(1, ShouldNotHaveSameTypeAs, "asdf")

// 7. 时间和时间间隔类
// So(time.Now(), ShouldHappenBefore, time.Now())
// So(time.Now(), ShouldHappenOnOrBefore, time.Now())
// So(time.Now(), ShouldHappenAfter, time.Now())
// So(time.Now(), ShouldHappenOnOrAfter, time.Now())
// So(time.Now(), ShouldHappenBetween, time.Now(), time.Now())
// So(time.Now(), ShouldHappenOnOrBetween, time.Now(), time.Now())
// So(time.Now(), ShouldNotHappenOnOrBetween, time.Now(), time.Now())
// So(time.Now(), ShouldHappenWithin, duration, time.Now())
// So(time.Now(), ShouldNotHappenWithin, duration, time.Now())

// 8. 自定义断言方法
// 如果上面列出来的断言方法都不能满足你的需要，那么你还可以按照下面的格式自定义一个断言方法。
// 注意：<>中的内容是你需要按照实际需求替换的内容。
// func should<do-something>(actual interface{}, expected ...interface{}) string {
//     if <some-important-condition-is-met(actual, expected)> {
//         return ""   // 返回空字符串表示断言通过
//     }
//     return "<一些描述性消息详细说明断言失败的原因...>"
// }

// goconvey提供全自动的WebUI，只需要在项目目录下执行以下命令:
// goconvey
// 默认就会在本机的8080端口提供WebUI界面，十分清晰地展现当前项目的单元测试数据。
