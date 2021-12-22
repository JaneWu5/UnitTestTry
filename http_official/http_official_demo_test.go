package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTitle(t *testing.T) {
	// mock一个HTTP请求
	r := httptest.NewRequest(
		// 请求方法
		"GET",
		// 请求URI
		"/view/TestPage",
		// 请求参数
		strings.NewReader(""),
	)

	// mock一个响应记录器
	w := httptest.NewRecorder()
	title, err := getTitle(w, r)
	if err != nil {
		t.Fatalf("err happend: %v", err)
	}
	t.Logf("title: %s", title)
}
