package chooserbenchmark

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	s := time.Now()
	l := NewLogNormalLatency(&LatencyConfig{Median: time.Millisecond * 100, Pin: true})
	minVal := math.MaxFloat64
	maxVal := math.SmallestNonzeroFloat64
	for i := 0; i < 10000000; i++ {
		l.Random()
		//r := l.Random()
		//if r < minVal {
		//	minVal = r
		//}
		//if r > maxVal {
		//	maxVal = r
		//}
	}
	e := time.Now()
	fmt.Println(fmt.Sprintf("max: %v", time.Duration(maxVal)))
	fmt.Println(fmt.Sprintf("min: %v", time.Duration(minVal)))
	fmt.Println(fmt.Sprintf("time: %v", time.Duration(e.Sub(s))))
}
