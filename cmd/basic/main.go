package main

import (
	"fmt"
	progressparser "sys-monitor/progress_parser"
)

func main() {

	info := progressparser.Sysinfo
	fmt.Println(info.Processes[0])
}
