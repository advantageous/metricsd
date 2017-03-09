package test

import (
	"fmt"
	"testing"
)

func TestScratch(z *testing.T) {
	var total = `          total        used        free      shared  buff/cache   available
Mem:        7747784      130380     6942216       16624      675188     7355592
Swap:             0           0           0
`
	fmt.Print(total)

}
