package goee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name/:age", nil)
	r.addRouter("GET", "/hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "/asserts/*filepath", nil)
	return r
}

// 测试解析
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = reflect.DeepEqual(parsePattern("/p/*name"), []string{"p", "*name"})

	if !ok {
		t.Fatal("test ParsePattern failed!")
	}
}

// 测试路由解析
func TestGetRouter(t *testing.T) {
	r := newTestRouter()
	node, params := r.getRouter("GET", "/hello/zhangsan/1c")

	if node == nil {
		t.Fatal("nil should not be returned")
	}

	fmt.Println(node.pattern)
	if node.pattern != "/hello/:name/:age" {
		t.Fatal("should match /hello/:name/:age")
	}

	if params["name"] != "zhangsan" {
		t.Fatal("name should be equal to 'zhangsan'")
	}
	fmt.Printf("matched pattern:%s,params['name']:%s\n", node.pattern, params["name"])
}
