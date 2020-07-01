// 版权 @2019 凹语言 作者。保留所有权利。

// 凹模块信息(wa.json文件).
package mod

import (
	"encoding/json"
	"os"
	"strings"
)

// 模块信息
type Module struct {
	Path        string   `json:"path"`        // 包路径
	Version     string   `json:"version"`     // 版本信息
	Description string   `json:"description"` // 描述信息
	Keywords    []string `json:"keywords"`    // 关键字
	Author      []string `json:"author"`      // 作者
	License     string   `json:"license"`     // 版权
}

func Load(path string) (*Module, error) {
	if !strings.HasSuffix(path, "wa.json") {
		path = path + "/wa.json"
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info := new(Module)
	err = json.NewDecoder(f).Decode(info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
