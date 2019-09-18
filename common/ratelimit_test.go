package common

import (
	"testing"
	"time"
)

func TestAccess(t *testing.T) {
	lb := &LeakyBucket{
		capcity:  10,
		interval: time.Second,
	}
	// 做10次
	for i := 0; i < 10; i++ {
		if !lb.Access() {
			t.Errorf("test fail: %v", i)
		}
	}
	time.Sleep(time.Second * 1)
	// 第11的操作
	if !lb.Access() {
		t.Error("11 test fail")
	}
}
