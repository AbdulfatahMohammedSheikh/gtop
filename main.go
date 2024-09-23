package main

import (
	"github.com/gin-gonic/gin"
	progressparser "sys-monitor/progress_parser"
	"sys-monitor/util"
)

func main() {

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {

		// TODO: get initial inforation like os name - kernal - system info

		info := progressparser.Sysinfo

		ctx.HTML(200, "home.html", gin.H{
			"info": info,
		})

	})

	r.GET("/body", func(ctx *gin.Context) {

		type Info struct {
			Processes *[]progressparser.Process
		}

		info := Info{Processes: progressparser.NewProcessInfo()}

		progressparser.NewProcessInfo()
		ctx.HTML(200, "table", gin.H{
			"info": info,
		})
	})

	r.GET("/time", func(ctx *gin.Context) {

		numberOfProcesses, err := progressparser.GetTotalNumberOfProcesses()

		if nil != err {
			// TODO: log error
		}

		numberOfRunningProcesses, err := progressparser.GetNumberOfRunningProcesses()

		if nil != err {
			// TODO: log error
		}

		numbrOfCores, err := progressparser.GetNumberOfCores()

		if nil != err {
			// TODO: log error
		}

		temp, err := progressparser.GetSysUpTime()

		if nil != err {
			// TODO: log error
		}
		uptime := util.FormateTime(*temp)

		info := progressparser.SysInfo{
			NumberOfProcesses:        *numberOfProcesses,
			NumberOfRunningProcesses: *numberOfRunningProcesses,
			NumberOfCores:            *numbrOfCores,
			UpTime:                   uptime,
		}

		ctx.HTML(200, "update-time", gin.H{
			"info": info,
		})

	})

	// TODO: create sample page

	r.LoadHTMLGlob("template/*")
	r.Static("/assets", "./static")
	r.Run(":8080")
}
