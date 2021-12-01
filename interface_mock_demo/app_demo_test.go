package interface_mock_demo

import (
	mocks "UnitTestTry/mocks/interface_mock"
	"github.com/golang/mock/gomock"
	"testing"
)

// 执行 mockgen -source=app_demo.go -destination=../mocks/interface_mock/app_mock.go -package=mocks
// Go测试库testify目前也提供类似的mock工具—testify/mock和mockery

func TestGetFromDB(t *testing.T) {
	// 创建gomock控制器，用来记录后续的操作信息
	ctrl := gomock.NewController(t)
	// 断言期望的方法都被执行
	// Go1.14+的单测中不再需要手动调用该方法
	defer ctrl.Finish()
	// 调用mockgen生成代码中的NewMockDB方法
	// 这里mocks是我们生成代码时指定的package名称
	db := mocks.NewMockDB(ctrl)
	// 【打桩（stub）】
	// 软件测试中的打桩是指用一些代码（桩stub）代替目标代码，通常用来【屏蔽】或【补齐】业务逻辑中的关键代码方便进行单元测试。
	// 屏蔽：不想在单元测试用引入数据库连接等重资源
	// 补齐：依赖的上下游函数或方法还未实现
	// gomock支持针对参数、返回值、调用次数、调用顺序等进行打桩操作。
	// 当传入Get函数的参数为wuzheng.com时返回1和nil —— 【屏蔽】操作，不去查询数据库，直接返回制定值(1,nil)
	db.EXPECT().
		Get(gomock.Eq("wuzheng.com")). // 参数
		Return(1, nil).                // 返回值
		Times(1)
	// 调用次数
	// Times() 断言 Mock 方法被调用的次数。
	// MaxTimes() 最大次数。
	// MinTimes() 最小次数。
	// AnyTimes() 任意次数（包括 0 次）。

	// db.EXPECT().Get(gomock.Not("q1mi")).Return(10, nil).Times(1)
	// db.EXPECT().Get(gomock.Any()).Return(20, nil).Times(1)
	// db.EXPECT().Get(gomock.Nil()).Return(-1, nil).Times(1)

	// 调用GetFromDB函数时传入上面的mock对象m
	if v := GetFromDB(db, "wuzheng.com"); v != 1 {
		t.Fatal("Not expected return value")
	}
	// 再次调用上方mock的Get方法时不满足调用次数为1的期望
	//if v := GetFromDB(db, "wuzheng.com"); v != 1 {
	//	t.Fatal()
	//}

	// 对返回值的处理
	db.EXPECT().Get(gomock.Any()).Return(20, nil)
	if v := GetFromDB(db, "wuzheng1.com"); v != 20 {
		t.Fatalf("GetFromDB value: %d, not expected!", v)
	}
	db.EXPECT().Get(gomock.Any()).Do(func(key string) {
		t.Logf("input key is %v\n", key)
	})
	GetFromDB(db, "wuzheng2.com")
	db.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
		t.Logf("input key is %v\n", key)
		return 10, nil
	})
	if v := GetFromDB(db, "wuzheng3.com"); v != 10 {
		t.Fatalf("GetFromDB value: %d", v)
	}

	// 指定顺序
	gomock.InOrder(
		db.EXPECT().Get("1"),
		db.EXPECT().Get("2"),
		db.EXPECT().Get("3"),
	)
	// 按顺序调用
	GetFromDB(db, "1")
	GetFromDB(db, "2")
	GetFromDB(db, "3")
}
