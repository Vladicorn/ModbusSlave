package main

import (
	"ModbusSlave/pkg/Handlers"
	"ModbusSlave/pkg/Modbus"
	"fmt"
	"net/http"
)

func main() {

	var telStr Handlers.TelStr
	var telStrClin Handlers.TelStrClin

	ch := make(chan Modbus.TelegramAnsver)
	chClin := make(chan Modbus.TelegramAnsverSlave)

	chansv := make(chan bool)
	chansvClient := make(chan bool)

	chanCancel := make(chan bool)
	chanStart := make(chan bool)

	telStr.Ch = ch
	telStr.Chansv = chansv
	telStr.ChStop = chanCancel
	telStr.ChStart = chanStart

	telStr.ChFlagStart = false
	telStrClin.ChFlagStart = false

	telStrClin.Chansv = chansvClient
	telStrClin.Ch = chClin
	telStrClin.ChStart = chanStart
	telStrClin.ChStop = chanCancel

	fmt.Println("Server is listening...")
	go http.HandleFunc("/stop", telStr.HomeStop)
	go http.HandleFunc("/start", telStr.HomeStart)
	go http.HandleFunc("/startClient", telStrClin.HomeStartCLient)
	go http.HandleFunc("/stopClient", telStrClin.HomeStopClient)
	go http.HandleFunc("/", telStr.HomeIndex)

	http.ListenAndServe(":8181", nil)

}
