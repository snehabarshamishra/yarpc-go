// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package loadbalancebenchmark

import (
	"fmt"
	"strings"
)

func PrintToTerminal(ids []int, cnts []int, types []MachineType, maxCnt int, size int) {
	maxStarCnt := 60
	base := float64(maxCnt) / float64(maxStarCnt)

	for i := 0; i < size; i++ {
		var sType string
		id, cnt, tp := ids[i], cnts[i], types[i]
		if tp == SlowMachine {
			sType = "-"
		} else if tp == NormalMachine {
			sType = " "
		} else if tp == FastMachine {
			sType = "+"
		}
		starCnt := int(float64(cnt) / base)
		fmt.Println(fmt.Sprintf("%s\t%d\t%s", sType, id, strings.Repeat("*", starCnt)))
	}
}
