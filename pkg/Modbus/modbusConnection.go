package Modbus

import (
	"fmt"
	"net"
)

var Quantity byte

func ModbusConnect(chOut chan TelegramAnsver, chanbool chan bool, chanStop chan bool, chanStart chan bool) {

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":502")
	chIn := make(chan TelegramAnsver)

	// Открываем порт
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	go handleTCPConnection(conn, chIn)
	go bufferChan(chanbool, chIn, chOut)

	stop := <-chanStop

	if stop {
		conn.Close()
		ln.Close()
	}
}

func ModsbusConnectClient(chOut chan TelegramAnsverSlave, chanbool chan bool, chanStop chan bool, chanStart chan bool) {

	chIn := make(chan TelegramAnsverSlave)
	Quantity = 2
	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:502")
	RegSliceClient[11] = Quantity
	go handleTCPConnectionClient(conn, chIn)
	go bufferChanClient(chanbool, chIn, chOut)

	stop := <-chanStop
	if stop {
		conn.Close()
	}
}
