package main

import (
	"fmt"
	"sync"
	"testing"
)


//Test方法 可能打印不出来颜色
func TestLoggerPrint(t *testing.T) {
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


