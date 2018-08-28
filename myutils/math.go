package myutils

import (
	"math"
)

//CalcucateMaxPage 分页计算最大页码用方法
func CalcucateMaxPage(totalCount int, pageSize int) int {
	f := math.Ceil(float64(totalCount) / float64(pageSize))
	hv := int(f)
	if f != float64(hv) {
		hv = hv + 1
	}
	return hv
}
