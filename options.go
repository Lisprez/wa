// 版权 @2019 凹语言 作者。保留所有权利。

package wa

import (
	"go/types"
)

// 可选参数
type Options func(*options)

// 可选参数内部结构
type options struct {
	DebugMode bool // 调试模式(针对Wa语言开发者)

	// vfs?
	Env  map[string]string // 环境变量
	Args []string          // 程序的运行参数

	WaOS      string // 操作系统
	WaArch    string // CPU类型
	WaRoot    string // 包根路径
	WaBackend string // 后端类型

	Optimize  bool        // 是否开启优化
	BuildTags []string    // 额外的构建标志
	Sizes     types.Sizes // 重定义机器字大小
}

func (p *options) applyOptions(opts ...Options) *options {
	for _, fn := range opts {
		fn(p)
	}
	// 初始化默认参数
	return p
}
