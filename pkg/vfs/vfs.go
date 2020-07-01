// 版权 @2019 凹语言 作者。保留所有权利。

package vfs

import (
	"io/ioutil"
	"path/filepath"
	"sort"
)

// 目录
type Files map[string]string

// 文件系统
type FileSystem map[string]Files

// 从磁盘加载目录下的全部凹代码(不递归)
func LoadFiles(path string, filter func(name string) bool) (Files, error) {
	m := make(map[string]string)

	wa_matches, err := filepath.Glob(path + "/*.wa")
	if err != nil {
		return m, err
	}

	wago_matches, err := filepath.Glob(path + "/*.wa.go")
	if err != nil {
		return m, err
	}

	for _, name := range append(wa_matches, wago_matches...) {
		if filter != nil && !filter(name) {
			continue
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			return m, err
		}

		m[filepath.Base(name)] = string(data)
	}

	return m, nil
}

// 加载整个包文件系统
// TODO: 支持 vendor
// TODO: 支持 wa.mod 文件指定的当前包路径
func LoadFileSystem(root string, filter func(name string) bool, pkgs ...string) (FileSystem, error) {
	if len(pkgs) == 0 {
		pkgs = walkWaPkgList(root)
	}
	prog := make(map[string]Files)
	for _, importPath := range pkgs {
		pkgpath := filepath.Join(root, importPath)
		pkg, err := LoadFiles(pkgpath, filter)
		if err != nil {
			return prog, err
		}
		if len(pkg) != 0 {
			prog[importPath] = pkg
		}
	}
	return prog, nil
}

func (d Files) Files() []string {
	var ss []string
	for s, _ := range d {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	return ss
}

func (fs FileSystem) Pkgs() []string {
	var ss []string
	for s, _ := range fs {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	return ss
}

func (fs FileSystem) Merge(m FileSystem) FileSystem {
	for k, v := range m {
		fs[k] = v
	}
	return fs
}
