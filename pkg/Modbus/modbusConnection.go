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

func handleTCPConnectionClient(conn net.Conn, ch chan TelegramAnsverSlave) {
	defer conn.Close()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:

			conn.Write(RegSliceClient)
			RegSliceClient[1] = RegSliceClient[1] + 1
			if RegSliceClient[1] > 244 {
				RegSliceClient[0] = RegSliceClient[0] + 1
				RegSliceClient[1] = 0
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
		ch <- teleg

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

//буфферная функция
func bufferChan(chanbool chan bool, chIn, chOut chan TelegramAnsver) {
	var bufval TelegramAnsver
	for {
		select {
		case bufval = <-chIn:
			//fmt.Println(bufval)
		case <-chanbool:
			fmt.Println("Handle request")
			chOut <- bufval
		}
	}
}

//буфферная функция
//func bufferChanClient(chanbool chan bool, chIn, chOut chan TelegramAnsverSlave) {
func bufferChanClient(chanbool chan bool, chIn, chOut chan TelegramAnsverSlave) {

	var bufval TelegramAnsverSlave
	for {
		fmt.Println("Handle request")
		select {
		case bufval = <-chIn:
			fmt.Println(bufval)
		case <-chanbool:
			fmt.Println("123123Handle request")
			chOut <- bufval

		}
	}
}
