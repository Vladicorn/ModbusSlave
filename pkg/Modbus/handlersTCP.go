package Modbus

import (
	"fmt"
	"io"
	"net"
	"time"
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
