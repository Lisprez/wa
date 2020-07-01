// 版权 @2019 凹语言 作者。保留所有权利。

package wa

import (
	"os"
)

func isFile(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode().IsRegular()
	}
	return false
}
