package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type PDU struct {
	FuncCode uint8
	FirstReg uint16
	//	Data     []uint16
	CountReg uint16
}

type PDUAnsver struct {
	FuncCode uint8
	CountReg uint8
	Data     []byte
}

type Telegram struct {
	MBAP MBAPHeader
	PDU  PDU
}
type TelegramAnsver struct {
	MBAP MBAPHeader
	PDU  PDUAnsver
}

type MBAPHeader struct {
	TranID     uint16
	ProtocolID uint16
	Length     uint16
	UnitID     uint8
}

var regSlice = []byte{
	0x00, 0x0A,
	0x00, 0xFF,
	0x00, 0x01,
	0x00, 0x03,
	0x00, 0x04,
	0x00, 0x00,
	0x00, 0x05,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
}
var slaveID uint8 = 2

func main() {

	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":502")

	// Открываем порт
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Запускаем цикл
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
		fmt.Println("Message Received:", buf)

		conn.Write(ReadHoldingRegister(buf))
	}
}

//парсинг модбас
func ReadHoldingRegister(data []byte) []byte {
	var telegram Telegram
	conn := bytes.NewBuffer(data)
	//чтение хеадер
	buf := make([]byte, 7)

	if count, err := conn.Read(buf); err != nil {
		fmt.Println(err)
		return nil
	} else {
		if count != 7 {
			fmt.Println("More 7")
			return nil
		}
	}

	//обработка хедера
	telegram.MBAP.TranID = binary.BigEndian.Uint16(buf[:2])
	telegram.MBAP.ProtocolID = binary.BigEndian.Uint16(buf[2:4])
	telegram.MBAP.Length = binary.BigEndian.Uint16(buf[4:6])
	telegram.MBAP.UnitID = buf[6]
	if telegram.MBAP.UnitID == slaveID {
		//обработка PDU
		buf = make([]byte, telegram.MBAP.Length-1)
		if _, err := conn.Read(buf); err != nil {
			fmt.Println(err)
			return nil
		}
		telegram.PDU.FuncCode = buf[0]
		telegram.PDU.FirstReg = binary.BigEndian.Uint16(buf[1:3])
		telegram.PDU.CountReg = binary.BigEndian.Uint16(buf[3:5])

		var telegDataMB TelegramAnsver
		switch telegram.PDU.FuncCode {
		case 3:
			telegDataMB = ReadDataMBHolding(&telegram)
		case 2:
			fmt.Println("In process")
		default:
			fmt.Println("Unknow FuncCode")
		}

		var str []byte
		str1 := make([]byte, 2)
		binary.BigEndian.PutUint16(str1, telegDataMB.MBAP.TranID)
		str2 := make([]byte, 2)
		binary.BigEndian.PutUint16(str2, telegDataMB.MBAP.ProtocolID)
		str3 := make([]byte, 2)
		binary.BigEndian.PutUint16(str3, telegDataMB.MBAP.Length)
		str4 := telegDataMB.MBAP.UnitID
		str5 := telegDataMB.PDU.FuncCode
		str6 := telegDataMB.PDU.CountReg

		str = append(str, str1...)
		str = append(str, str2...)
		str = append(str, str3...)
		str = append(str, str4)
		str = append(str, str5)
		str = append(str, str6)
		str = append(str, telegDataMB.PDU.Data...)

		return str
	} else {
		fmt.Println("Another SlaveID")
	}

	return nil
}

//Разбор дата поинт
func ReadDataMBHolding(teleg *Telegram) TelegramAnsver {

	slicereg := make([]byte, teleg.PDU.CountReg*2, 100)
	var telega TelegramAnsver

	j := 0
	//проверка на количество считываний не сделано

	for i := teleg.PDU.FirstReg; i < teleg.PDU.CountReg*2; i++ {
		slicereg[j] = regSlice[i]
		j++
	}

	telega.MBAP = teleg.MBAP
	telega.MBAP.Length = telega.MBAP.Length + 3
	telega.PDU.FuncCode = teleg.PDU.FuncCode
	telega.PDU.CountReg = uint8(j)
	telega.PDU.Data = slicereg

	return telega

}
