// 版权 @2019 凹语言 作者。保留所有权利。

package wa

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/wa-lang/wa/pkg/3rdparty/cli"
)

// 命令行主入口函数
func Main() {
	var app = cli.NewApp()

	app.Name = "wa"
	app.Usage = "Wa Programming Language (凹语言) - https://github.com/wa-lang"
	app.Version = "v0.0.1"

	app.Authors = []cli.Author{
		cli.Author{Name: "柴树杉", Email: "chaishushan@gmail.com"},
		cli.Author{Name: "丁尔男"},
		cli.Author{Name: "史斌"},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:   "waarch",
			EnvVar: "WAARCH",
			Value:  runtime.GOARCH,
		},
		&cli.StringFlag{
			Name:   "waos",
			EnvVar: "WAOS",
			Value:  runtime.GOOS,
		},
		&cli.StringFlag{
			Name:   "waroot",
			EnvVar: "WAROOT",
			Value: func() string {
				if s, _ := os.UserHomeDir(); s != "" {
					return filepath.Join(s, "wa")
				}
				return ""
			}(),
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "run-script",
			Usage:     "run wa script",
			ArgsUsage: "hello.wa",

			Action: func(c *cli.Context) error {
				if len(c.Args()) == 0 {
					fmt.Println("no file")
					return nil
				}

				err := NewScript().MustLoad(c.Args().First()).Run()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				return nil
			},
		},

		{
			Name:   "run",
			Hidden: true,

			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:   "build",
			Hidden: true,

			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}

// ${WaRoot} or ${HOME}/wa
func getWaRoot() string {
	if s := os.Getenv("WaRoot"); s != "" {
		return s
	}
	if s, _ := os.UserHomeDir(); s != "" {
		return filepath.Join(s, "wa")
	}
	return ""
}
