// 版权 @2019 凹语言 作者。保留所有权利。

package wa_test

import (
	"fmt"
	"log"

	"github.com/wa-lang/wa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func Example_script() {
	app := wa.NewScript().MustLoad("hello.wa", `
		package main

		func main() {
			println("你好, 凹语言!")

			var sum int
			for i := 1; i <= 100; i++ {
				sum += i
			}

			println(sum)
		}
	`)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// 你好, 凹语言!
	// 5050
}

func Example_userFunc() {
	app := wa.NewScript().MustLoad("hello.wa", `
		package main

		func sayHello(s string)

		func main() {
			sayHello("你好, 凹语言!")
		}
	`)

	app.DefineFunc("hello.wa.sayHello", func(fr *wa.Frame, args ...wa.Value) wa.Value {
		for _, a := range args {
			switch a := a.(type) {
			case []watypes.Value:
				for _, a := range a {
					fmt.Print(a)
				}
			default:
				fmt.Print(a)
			}
		}
		return nil
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// 你好, 凹语言!
}

func _ExampleScript_call() {
	w := wa.NewScript()

	w.DefineFunc("println", func(fr *wa.Frame, a ...wa.Value) wa.Value {
		fmt.Println(a)
		return nil
	})

	// 加载单个凹脚本文件
	if err := w.Load("hello.wa", nil); err != nil {
		log.Fatal(err)
	}

	// 调用其中的 main 函数
	v, err := w.Call("hello.wa.main")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)
}
