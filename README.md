# 凹语言

凹语言（凹读音“Wa”）是[柴树杉](https://github.com/chai2010)、[丁尔男](https://github.com/3dgen)和[史斌](https://github.com/benshi001)设计的实验性编程语言。目前推出的凹脚本语言是一种可以嵌入Go语言环境的脚本语言。

凹语言主页：https://wa-lang.org

```
+---+    +---+
| o |    | o |
|   +----+   |
|            |
|    \/\/    |
|            |
+------------+
```

## 安装环境

1. `go get github.com/wa-lang/wa/cmd/wa`

## 运行例子(命令行)

[_hello.wa](_hello.wa) 程序:

```go
package main

func main() {
	println("你好, 凹语言!")

	var sum int
	for i := 1; i <= 100; i++ {
		sum += i
	}

	println(sum)
}
```

编译和运行:

```
$ wa run-script _hello.wa
你好, 凹语言!
5050
```

## 运行例子(嵌入Go语言)

```go
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
```

运行脚本:

```
$ go run hello.go 
你好, 凹语言!
5050
```

## 版权

版权 @2019 凹语言 作者。保留所有权利。
