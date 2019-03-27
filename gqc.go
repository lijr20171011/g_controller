package g_controller

import (
	"fmt"
	"my/printl"
)

// goroutine数量控制器
type GQC struct {
	maxGNum int           // 最大goroutine数量
	addCh   chan struct{} // go任务数量控制队列
	runCh   chan *goTask  // 控制器执行go任务队列
	closeCh chan struct{} // 控制器停止队列
	isClose bool          // 控制器是否关闭
}

// 执行函数参数包装结构
type goTask struct {
	f         goFunc        // 执行函数
	p         []interface{} // 执行参数
	isRecover bool          // 是否捕捉panic
}

// go任务函数类型
type goFunc func(...interface{}) interface{}

// 新建goroutine数量控制器
func NewGQC(n int) *GQC {
	// 新建控制器
	gqc := &GQC{
		maxGNum: n,
		addCh:   make(chan struct{}, n),
		runCh:   make(chan *goTask),
		closeCh: make(chan struct{}),
	}

	// 执行go任务
	go gqc.run()

	return gqc
}

// 添加goFunc
func (gqc *GQC) AddGoFunc(isRecover bool, f goFunc, p ...interface{}) bool {
	if gqc.isClose {
		return false
	}
	gt := goTask{
		f:         f,
		p:         p,
		isRecover: isRecover,
	}
	gqc.addCh <- struct{}{}
	gqc.runCh <- &gt
	return true
}

// 执行go任务
func (gqc *GQC) run() {
	for {
		select {
		case gt := <-gqc.runCh: // 从队列中取出任务
			// 执行go任务
			go func(*goTask) {
				defer func() {
					// 执行结束,数量减1
					<-gqc.addCh
				}()
				// panic控制
				if gt.isRecover {
					// 捕捉panic
					defer PanicRecover()
				}
				// 执行go任务
				gt.f(gt.p...)
			}(gt)
		case <-gqc.closeCh: // 控制器已关闭
			printl.Debug("控制器已关闭")
			return
		}
	}
}

// 关闭控制器
func (gqc *GQC) Close() {
	gqc.isClose = true
	close(gqc.closeCh)
}

// 捕捉panic
func PanicRecover() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}
