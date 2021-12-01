package varys

import (
	"fmt"
	"time"
)

type UserInfo struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
}

// 如果我们为GetInfo编写单元测试的时候CalcAge方法的功能还未完成，这个时候我们可以使用monkey进行打桩。
// CalcAge 计算用户年龄
func (u *UserInfo) CalcAge() int {
	t, err := time.Parse("2006-01-02", u.Birthday)
	if err != nil {
		return -1
	}
	return int(time.Now().Sub(t).Hours()/24.0) / 365
}

// GetInfo 获取用户相关信息
func (u *UserInfo) GetInfo() string {
	age := u.CalcAge()
	if age <= 0 {
		return fmt.Sprintf("%s很神秘，我们还不了解ta。", u.Name)
	}
	return fmt.Sprintf("%s今年%d岁了，ta是我们的朋友。", u.Name, age)
}

func GetInfoByUID(uid int64) (*UserInfo, error) {
	return &UserInfo{
		Uid:  uid,
		Name: "wz",
	}, nil
}
