package Modbus

import (
	"fmt"
	"io"
	"net"
	"time"
)

var Quantity byte

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

func ModsbusConnectClient() {
	ticker := time.NewTicker(time.Second)
	Quantity = 4
	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:502")
	RegSlice[11] = Quantity
	for {
		select {
		case <-ticker.C:
			conn.Write(RegSlice)
			RegSlice[1] = RegSlice[1] + 1
			if RegSlice[1] > 244 {
				RegSlice[0] = RegSlice[0] + 1
				RegSlice[1] = 0
			}
		}
		// Прослушиваем ответ
		lenresp := (Quantity*2 + 9)
		buf := make([]byte, lenresp)
		if _, err := conn.Read(buf); err != nil {
			if err == io.EOF {
				fmt.Println("Connection lost.")
				return
			} else {
				fmt.Println(err)
				return
			}
		}
		teleg := ReadHoldingRegisterPol(buf)
		fmt.Println(teleg)
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
