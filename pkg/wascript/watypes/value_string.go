// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watypes

import (
	"io"
	"strings"
)

type StringIter struct {
	*strings.Reader
	i int
}

func (it *StringIter) Next() Tuple {
	okv := make(Tuple, 3)
	ch, n, err := it.ReadRune()
	ok := err != io.EOF
	okv[0] = ok
	if ok {
		okv[1] = it.i
		okv[2] = ch
	}
	it.i += n
	return okv
}
