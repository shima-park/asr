package log

import (
	"fmt"
	//"strconv"
)

var LogOut int = 1

func W(a ...interface{}) {
	if LogOut == 1 {
		fmt.Println(a)
	}
}
