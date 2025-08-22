package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/molast/crawler-core/app"
	"github.com/molast/crawler-core/logs"
	"github.com/molast/crawler-core/runtime/cache"
	"github.com/molast/crawler-core/runtime/status"
)

var (
	spiderflag *string
)

func main() {
	app.LogicApp.Init(cache.Task.Mode, cache.Task.Port, cache.Task.Master)
	if cache.Task.Mode == status.UNSET {
		return
	}
	run()
}

func loop() {
	for {
		parseInput()
		run()
	}
}

func run() {
	sps := app.LogicApp.GetSpiderLib()
	app.LogicApp.SpiderPrepare(sps).Run()
}

func parseInput() {
	logs.Log.Informational("\n添加任务参数——必填：%v\n添加任务参数——必填可选：%v\n", "-c_spider", []string{
		"-a_keyins",
		"-a_limit",
		"-a_outtype",
		"-a_thread",
		"-a_pause",
		"-a_proxysecond",
		"-a_dockercap",
		"-a_success",
		"-a_failure"})
	logs.Log.Informational("\n添加任务：\n")
retry:
	*spiderflag = ""
	input := [12]string{}
	fmt.Scanln(&input[0], &input[1], &input[2], &input[3], &input[4], &input[5], &input[6], &input[7], &input[8], &input[9])
	if strings.Index(input[0], "=") < 4 {
		logs.Log.Informational("\n添加任务的参数不正确，请重新输入：")
		goto retry
	}
	for _, v := range input {
		i := strings.Index(v, "=")
		if i < 4 {
			continue
		}
		key, value := v[:i], v[i+1:]
		switch key {
		case "-a_keyins":
			cache.Task.Keyins = value
		case "-a_limit":
			limit, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				break
			}
			cache.Task.Limit = limit
		case "-a_outtype":
			cache.Task.OutType = value
		case "-a_thread":
			thread, err := strconv.Atoi(value)
			if err != nil {
				break
			}
			cache.Task.ThreadNum = thread
		case "-a_pause":
			pause, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				break
			}
			cache.Task.Pausetime = pause
		case "-a_proxysecond":
			proxySecond, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				break
			}
			cache.Task.ProxySecond = proxySecond
		case "-a_dockercap":
			dockercap, err := strconv.Atoi(value)
			if err != nil {
				break
			}
			if dockercap < 1 {
				dockercap = 1
			}
			cache.Task.DockerCap = dockercap
		case "-a_success":
			if value == "true" {
				cache.Task.SuccessInherit = true
			} else if value == "false" {
				cache.Task.SuccessInherit = false
			}
		case "-a_failure":
			if value == "true" {
				cache.Task.FailureInherit = true
			} else if value == "false" {
				cache.Task.FailureInherit = false
			}
		case "-c_spider":
			*spiderflag = value
		default:
			logs.Log.Informational("\n不可含有未知参数,必填参数:%v\n可选参数:%v\n", "-c_spider", []string{
				"-a_keyins",
				"-a_limit",
				"-a_outtype",
				"-a_thread",
				"-a_pause",
				"-a_proxysecond",
				"-a_dockercap",
				"-a_success",
				"-a_failure"})
			goto retry
		}
	}
}
