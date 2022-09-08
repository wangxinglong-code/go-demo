package cache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	//设置普通缓存
	MemCache().Set("hello", 1, 10*time.Second)
	MemCache().Delete("hello")
	str, ok := MemCache().Get("hello")
	fmt.Println(reflect.TypeOf(str), ok)

	//设置不过期缓存
	MemCache().Set("buguoqi", "1111", NoExpiration)
	MemCache().Delete("buguoqi")
	str2, ok := MemCache().Get("buguoqi")
	fmt.Println(reflect.TypeOf(str2), ok)
}
