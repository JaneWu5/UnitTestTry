package redis_demo

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

// miniredis是一个纯go实现的用于单元测试的redis server。
// 它是一个简单易用的、基于内存的redis替代品，它具有真正的TCP接口，你可以把它当成是redis版本的net/http/httptest。
// 当我们为一些包含Redis操作的代码编写单元测试时就可以使用它来mock Redis操作。
// go get github.com/alicebob/miniredis/v2

func TestDoSomethingWithRedis(t *testing.T) {
	// mock一个redis server
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// 准备数据
	s.Set("q1mi", "wuzheng@qq.com")
	s.SAdd(KeyValidWebsite, "q1mi")

	// 连接mock的redis server
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(), // mock redis server的地址
	})

	// 调用函数
	ok := DoSomethingWithRedis(rdb, "q1mi")
	if !ok {
		t.Fatal("DoSomethingWithRedis returns false")
	}
	// 可以手动检查redis中的值是否复合预期
	if got, err := s.Get("blog"); err != nil || got != "https://wuzheng@qq.com" {
		t.Fatalf("'blog' has the wrong value")
	}
	// 也可以使用帮助工具检查
	s.CheckGet(t, "blog", "https://wuzheng@qq.com")
	// 过期检查
	s.FastForward(5 * time.Second) // 快进5秒
	if s.Exists("blog") {
		t.Fatal("'blog' should not have existed anymore")
	}
}
