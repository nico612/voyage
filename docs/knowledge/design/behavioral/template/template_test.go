package template

import (
	"fmt"
	"testing"
)

func Test_doCook(t *testing.T) {
	xihongshi := &XiHongShi{}
	doCook(xihongshi)

	fmt.Println("\n 做另外一道菜")

	chaojidan := &ChaoJiDan{}
	doCook(chaojidan)
}
