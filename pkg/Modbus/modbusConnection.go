package Modbus

import (
	"fmt"
	"io"
	"net"
)

//Accept после тсп конекта
func handleTCPConnection(conn net.Conn, ch chan TelegramAnsver) {
	defer conn.Close()
	for {

		//считываем 12 байт в буффер
		buf := make([]byte, 12)
		if _, err := conn.Read(buf); err != nil {
			if err == io.EOF {
				fmt.Println("Connection lost.")
				return
			} else {
				fmt.Println(err)
				return
			}
		}

		//	fmt.Println("Message Received:", buf)
		telegW, teg := ReadHoldingRegister(buf)
		conn.Write(telegW)
		ch <- teg
	}
}

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

//буфферная функция
func bufferChan(chanbool chan bool, chIn, chOut chan TelegramAnsver) {
	var bufval TelegramAnsver
	for {
		select {
		case bufval = <-chIn:
		//	fmt.Println("Blank")
		case <-chanbool:
			fmt.Println("Handle request")
			chOut <- bufval
		}
	}
}
