// 版权 @2019 凹语言 作者。保留所有权利。

package vfs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isWaFile(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".wa")
}

func walkWaPkgList(root string) (pkgs []string) {
	m := make(map[string]string)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isWaFile(info.Name()) {
			s := filepath.ToSlash(filepath.Dir(path))
			m[s] = s
		}
		return nil
	})
	for s := range m {
		pkgs = append(pkgs, s)
	}
	sort.Strings(pkgs)
	return pkgs
}
