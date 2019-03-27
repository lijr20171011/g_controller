package g_controller_test

import (
	"fmt"
	"time"

	"../g_controller"
)

func Example() {
	// 新建gqc 最多同时执行4个goroutine
	gqc := g_controller.NewGQC(4)
	// defer gqc.Close()

	// 执行函数
	f := func(p ...interface{}) interface{} {
		fmt.Println(p...)
		time.Sleep(1 * time.Second)
		panic("===") // 测试捕捉panic
		return nil
	}

	// 添加10个任务
	for i := 0; i < 10; i++ {
		if i == 7 {
			gqc.Close() // 测试关闭控制器
		}
		// 添加任务
		fmt.Println(gqc.AddGoFunc(
			true, // 是否捕捉panic
			f,
			"aaa", i, nil, // 参数
		), i)
	}

	// 防止测试时死锁
	go func() {
		for {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(10 * time.Second)
		}
	}()
	select {}
}
