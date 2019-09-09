// 凹语言 版权 @2019 柴树杉 & 丁尔男。保留所有权利。

// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/wa-lang/wa/pkg/wascript"
)

func main() {
	ctx := wascript.DefaultContext()

	ctx.FileSystem["myapp"] = map[string]string{
		"_hello.wa": `
			package main

			func main() {
				println("你好, 凹语言!")

				var sum int
				for i := 1; i <= 100; i++ {
					sum += i
				}

				println(sum)
			}
		`,
	}

	program, err := wascript.LoadProgram(ctx, "myapp")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
