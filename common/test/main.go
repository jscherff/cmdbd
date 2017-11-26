package main

import (
	`fmt`
	`github.com/jscherff/cmdbd/common`
)

func main() {

	x := common.CallerInfo()

	fmt.Printf("Base:\t%s\nPath:\t%s\nFile:\t%s\nFunc:\t%s\n",
		x.Base, x.Path, x.File, x.Func,
	)
}
