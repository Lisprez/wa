// 版权 @2019 凹语言 作者。保留所有权利。

// +build ignore

package main

import (
	"log"
	"os"

	"github.com/wa-lang/wa"
)

func main() {
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

	if err := app.Run(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
