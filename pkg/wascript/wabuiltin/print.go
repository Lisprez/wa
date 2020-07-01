// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"bytes"
	"os"
	"sync"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

var (
	PrintOutputMutex sync.Mutex
	PrintOutput      *bytes.Buffer
)

func Print(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	PrintOutputMutex.Lock()
	defer PrintOutputMutex.Unlock()

	ln := fn.Name() == "println"
	var buf bytes.Buffer

	for i, arg := range args {
		if i > 0 && ln {
			buf.WriteRune(' ')
		}
		buf.WriteString(watypes.ToString(arg))
	}
	if ln {
		buf.WriteRune('\n')
	}

	if PrintOutput != nil {
		PrintOutput.Write(buf.Bytes())
	} else {
		os.Stdout.Write(buf.Bytes())
	}
	return nil
}
