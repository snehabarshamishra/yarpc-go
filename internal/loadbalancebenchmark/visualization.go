package loadbalancingbenchmark

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
