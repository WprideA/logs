package logs

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
	lDebug = iota
	lInfo
	lWarning
	lError
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
	assignLog(nowTime)
	Debug("Logs init over ! ")
}

func assignLog(t time.Time){
	logRWMutex.Lock()
	defer  logRWMutex.Unlock()
	logPath := filePath+"\\"+operationtime.ParseTimeToString3(t)+".log"
	fmt.Println("LogPath :",logPath)
	file,err := os.OpenFile(logPath,os.O_WRONLY|os.O_CREATE|os.O_APPEND,0666)
	if err != nil {
		fmt.Println("err:",err)
	}
	logFile.file = file
	logFile.level = lDebug //默认为Debug，即全部输出
	logs = log.New(logFile.file,"",log.Ltime|log.Ldate)
	isToConsole =true    //默认开始输出到控制台
	// 定义多个写入器
	writers := []io.Writer{
		logFile.file,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logsc = log.New(fileAndStdoutWriter,"",log.Ltime|log.Ldate)

	go updateLogName()
}

//按照时间更新 log 名称
func updateLogName(){
	t := time.NewTicker(2*time.Second)
	pd := false
	for {
		select {
		case <-t.C:
			if operationtime.JudgeTime(0) && pd ==true {
				assignLog(time.Now())
			    pd = false
			}else {
				if !operationtime.JudgeTime(23){
					if !pd {
						pd =true
					}
				}
			}
		case <-ctx.Done():
			fmt.Println("Logs 监听时间，更新Log文件名称 关闭")
			t.Stop()   //停止定时器
			runtime.Goexit()
		}
	}
}

//关闭 logs
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
	if logFile.level <= lDebug {
		debug :=color.New(color.FgHiBlue).Sprint("[DEBUG]--",v)
		getLog().Printf(debug)
	}
}
func Info(v ...interface{}){
	if logFile.level <= lInfo {
		info :=color.New(color.FgHiGreen).Sprint("[INFO]--", v)
		getLog().Printf(info)
	}
}
func Warning(format string ,v ...interface{}){
	if logFile.level <= lWarning {
		warning :=color.New(color.FgHiMagenta).Sprint("[WARNING]--", v)
		getLog().Printf(warning)
	}
}
func Error(format string ,v ...interface{}){
	if logFile.level <= lError {
		error :=color.New(color.FgRed).Sprint("[ERROR]--", v)
		getLog().Printf(error)
	}
}

func Debugf (format string ,v ...interface{}){
	if logFile.level <= lDebug {
		debug :=color.New(color.FgHiBlue).Sprintf("[DEBUG]--"+format, v...)
		getLog().Printf(debug)
	}
}
func Infof (format string ,v ...interface{}){
	if logFile.level <= lInfo {
		info :=color.New(color.FgHiGreen).Sprintf("[INFO]--"+format, v...)
		getLog().Printf(info)
	}
}
func Warningf (format string ,v ...interface{}){
	if logFile.level <= lWarning {
		warning :=color.New(color.FgHiMagenta).Sprintf("[WARNING]--"+format, v...)
		getLog().Printf(warning)
	}
}
func Errorf (format string ,v ...interface{}){
	if logFile.level <= lError {
		error :=color.New(color.FgRed).Sprintf("[ERROR]--"+format, v...)
		getLog().Printf(error)
	}
}

