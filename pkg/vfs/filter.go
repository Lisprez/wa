// 版权 @2019 凹语言 作者。保留所有权利。

package vfs

import (
	"os"
	"runtime"
	"strings"
)

func Filter(envWaOS, envWaArch string) func(name string) bool {
	if envWaOS == "" {
		if s := os.Getenv("WaOS"); s != "" {
			envWaOS = s
		} else {
			envWaOS = runtime.GOOS
		}
	}
	if envWaArch == "" {
		if s := os.Getenv("WaArch"); s != "" {
			envWaArch = s
		} else {
			envWaArch = runtime.GOARCH
		}
	}

	return func(name string) bool {
		return goodOSArchFile(name, envWaOS, envWaArch)
	}
}

func goodOSArchFile(name, envWaOS, envWaArch string) bool {
	if dot := strings.Index(name, "."); dot != -1 {
		name = name[:dot]
	}

	i := strings.Index(name, "_")
	if i < 0 {
		return true
	}
	name = name[i:] // ignore everything before first _

	l := strings.Split(name, "_")
	if n := len(l); n > 0 && l[n-1] == "test" {
		l = l[:n-1]
	}

	var (
		knownOS = map[string]bool{
			"windows": true,
			"linux":   true,
			"darwin":  true,
		}
		knownArch = map[string]bool{
			"386":   true,
			"amd64": true,
			"arm64": true,
			"wasm":  true,
		}
	)

	n := len(l)
	if n >= 2 && knownOS[l[n-2]] && knownArch[l[n-1]] {
		return l[n-2] == envWaOS && l[n-1] == envWaArch
	}
	if n >= 1 && (knownOS[l[n-1]] || knownArch[l[n-1]]) {
		return l[n-1] == envWaOS || l[n-1] == envWaArch
	}

	return true
}
