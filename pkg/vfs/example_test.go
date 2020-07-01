// 版权 @2019 凹语言 作者。保留所有权利。

package vfs_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/wa-lang/wa/pkg/vfs"
)

func ExampleLoadFiles() {
	m, err := vfs.LoadFiles("./testdata/hello", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range m.Files() {
		fmt.Printf("%s: %s\n", s, strings.TrimSpace(m[s]))
	}

	// Output:
	// a.wa: // a.wa
	// b.wa: // b.wa
}

func ExampleLoadFileSystem() {
	prog, err := vfs.LoadFileSystem("./testdata", nil, "hello", "fmt")
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range prog.Pkgs() {
		fmt.Println(pkg)
		for _, s := range prog[pkg].Files() {
			fmt.Printf("%s: %s\n", s, strings.TrimSpace(prog[pkg][s]))
		}
	}

	// Output:
	// fmt
	// fmt.wa: // fmt.wa
	// string.wa: // string.wa
	// hello
	// a.wa: // a.wa
	// b.wa: // b.wa
}
