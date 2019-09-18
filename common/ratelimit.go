package common

import (
	"time"
)

// 漏桶算法
type LeakyBucket struct {
	capcity      int           // bucket 的总容量 10
	interval     time.Duration // 漏出水滴的时间  1s
	NumDrops     int           // 当前水滴的数量
	lastLeakTime time.Time     // 上次漏出的时间
}

// Access 主函数
func (l *LeakyBucket) Access() bool {
	now := time.Now()
	since := now.Sub(l.lastLeakTime)

	//漏出了水量
	leaks := int(float64(since) / float64(l.interval))

	if leaks > 0 {
		// 滑动窗口 // 归零  什么不归零
		// 为什么
		// leaks next pre leaks
		if l.NumDrops <= leaks {
			// 重置漏桶中的数量
			l.NumDrops = 0
		} else {
			// 减少漏桶中的数量
			l.NumDrops -= leaks
		}
		l.lastLeakTime = now
	}

	// 漏桶没有满的情况下
	if l.NumDrops < l.capcity {
		l.NumDrops++
		return true
	}
	// 已经满了呢？
	return false

}
