package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/WprideA/operationtime"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

//设置 log  等级
const (
	LDebug = iota
	LInfo
	LWarning
	LError
)
var isToConsole  bool  //是否输出到控制台


var logRWMutex sync.RWMutex
var logs *log.Logger   //创建 logger 对象
var logsc * log.Logger  //创建logger 对象（输出到控制台和 文件 s）


//Log 配置信息
type logging struct {
	level int
	file  *os.File
}

var logFile logging
var filePath string
var ctx  context.Context
var cf   context.CancelFunc
//获取当前程序所在的路径
func getwd() (string, error) {
	path, err := os.Getwd()     //获取当前程序的路径
	if err != nil {
		return "", errors.New(fmt.Sprintln("Getwd Error:", err))
	}
	ctx ,cf = context.WithCancel(context.Background())
	filepathstr := filepath.Join(path, "logs")
	err = os.MkdirAll(filepathstr, 0755)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Create Dirs Error:", err))
	}

	return filepathstr, nil
}

//Init function
func init (){
	filePath, _ = getwd()
	if filePath == ""{
		log.Fatal("Init    logs   error !")
	}
	time.Now()
	//获取当前时间
	nowTime := time.Now()
	AssignLog(nowTime)
	Debug("Logs init over ! ")
}

func AssignLog(t time.Time){
	logRWMutex.Lock()
	defer  logRWMutex.Unlock()
	logPath := filePath+"\\"+operationtime.ParseTimeToString3(t)+".log"
	fmt.Println("LogPath :",logPath)
	file,err := os.OpenFile(logPath,os.O_WRONLY|os.O_CREATE|os.O_APPEND,0666)
	if err != nil {
		fmt.Println("err:",err)
	}
	logFile.file = file
	logFile.level = LDebug   //默认为Debug，即全部输出
	logs = log.New(logFile.file,"",log.Ltime|log.Ldate)
	isToConsole =true    //默认开始输出到控制台
	// 定义多个写入器
	writers := []io.Writer{
		logFile.file,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logsc = log.New(fileAndStdoutWriter,"",log.Ltime|log.Ldate)

	go UpdateLogName()
}

//按照时间更新 log 名称
func UpdateLogName (){
	t := time.NewTicker(2*time.Second)
	pd := false
	for {
		select {
		case <-t.C:
			if operationtime.JudgeTime(0) && pd ==true {
			    fmt.Println("要更新啦~")
				AssignLog(time.Now())
			    pd = false
			}else {
				if !operationtime.JudgeTime(23){
					if !pd {
						fmt.Println("赋值哦")
						pd =true
					}
				}else {
					fmt.Println("不用更新哦")
				}
			}
		case <-ctx.Done():
			fmt.Println("Logs 监听时间，更新Log文件名称 关闭")
			t.Stop()   //停止定时器
			runtime.Goexit()
		}
	}
}


func Close(){
	if cf !=nil{
		cf()
	}
}




//设置等级
func SetLever (level  int ){
	logRWMutex.Lock()
	defer  logRWMutex.Unlock()
	logFile.level = level  // equal print all
}

//设置是否输出到控制台
func SetIsToConsole  (b bool){
	logRWMutex.Lock()
	defer  logRWMutex.Unlock()
	isToConsole = b
}

//获取get 获取 logger 对象
func getLog() *log.Logger{
	logRWMutex.RLock()
	defer  logRWMutex.RUnlock()
	if isToConsole{
		return  logsc
	}else {
		return  logs
	}
}

//以下方法为只将 log信息 输出到文件中
func Debug(v ...interface{}){
	if logFile.level <= LDebug {
		debug :=color.New(color.FgHiBlue).Sprint("[DEBUG]--",v)
		getLog().Printf(debug)
	}
}
func Info(v ...interface{}){
	if logFile.level <= LInfo {
		info :=color.New(color.FgHiGreen).Sprint("[INFO]--", v)
		getLog().Printf(info)
	}
}
func Warning(format string ,v ...interface{}){
	if logFile.level <= LWarning {
		warning :=color.New(color.FgHiMagenta).Sprint("[WARNING]--", v)
		getLog().Printf(warning)
	}
}
func Error(format string ,v ...interface{}){
	if logFile.level <= LError {
		error :=color.New(color.FgRed).Sprint("[ERROR]--", v)
		getLog().Printf(error)
	}
}

func Debugf (format string ,v ...interface{}){
	if logFile.level <= LDebug {
		debug :=color.New(color.FgHiBlue).Sprintf("[DEBUG]--"+format, v...)
		getLog().Printf(debug)
	}
}
func Infof (format string ,v ...interface{}){
	if logFile.level <= LInfo {
		info :=color.New(color.FgHiGreen).Sprintf("[INFO]--"+format, v...)
		getLog().Printf(info)
	}
}
func Warningf (format string ,v ...interface{}){
	if logFile.level <= LWarning {
		warning :=color.New(color.FgHiMagenta).Sprintf("[WARNING]--"+format, v...)
		getLog().Printf(warning)
	}
}
func Errorf (format string ,v ...interface{}){
	if logFile.level <= LError {
		error :=color.New(color.FgRed).Sprintf("[ERROR]--"+format, v...)
		getLog().Printf(error)
	}
}

func main (){
	wg := sync.WaitGroup{}
	for i := 0; i<10000;i++{
		wg.Add(1)
		func (){

			defer wg.Done()
			SetIsToConsole(true)
			Debugf("TestDebugf :::: %v" , 1 )
			Infof("TestInfof  ;:::%v ", 1)
			Warningf("TestWarningf  :::::%v",1)
			Errorf("TestErrorf ::::::%v", 1)
			fmt.Println("-----------------------------------------------------")
			Debug("TestDebug",1,"XXXX",2)
			Info("TestInfo",1,"XXXX",2)
			Warning("TestWarning",1,"XXXX",2)
			Error("TestError",1,"XXXX",2)
		}()
	}
	wg.Wait()
	Close()
}