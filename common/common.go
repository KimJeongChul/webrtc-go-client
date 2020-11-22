package common

import (
	"fmt"
	"log"
)

func LogE(funcName string, logMsg ...interface{}) {
	logStr := "[E] < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
func LogI(funcName string, trxId string, logMsg ...interface{}) {
	logStr := "[I] < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
func LogD(funcName string, trxId string, logMsg ...interface{}) {
	//todo :: Debug Log가 켜져 있는 경우만 출력하도록
	logStr := "[D] < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
