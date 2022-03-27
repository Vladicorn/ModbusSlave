package Modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var slaveID uint8 = 1

//парсинг модбас
func ReadHoldingRegister(data []byte) ([]byte, TelegramAnsver) {
	var telegram Telegram
	var dummy TelegramAnsver
	conn := bytes.NewBuffer(data)
	//чтение хеадер
	buf := make([]byte, 7)

	if count, err := conn.Read(buf); err != nil {
		fmt.Println(err)
		return nil, dummy
	} else {
		if count != 7 {
			fmt.Println("More 7")
			return nil, dummy
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
			return nil, dummy
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

		return str, telegDataMB
	} else {
		fmt.Println("Another SlaveID")
	}

	return nil, dummy
}

//Разбор дата поинт
func ReadDataMBHolding(teleg *Telegram) TelegramAnsver {
	slicereg := make([]byte, teleg.PDU.CountReg*2, 100)
	var telega TelegramAnsver
	j := 0
	//проверка на количество считываний не сделано
	for i := teleg.PDU.FirstReg; i < teleg.PDU.CountReg*2; i++ {
		slicereg[j] = RegSlice[i]
		j++
	}
	telega.MBAP = teleg.MBAP
	telega.MBAP.Length = telega.MBAP.Length + 3
	telega.PDU.FuncCode = teleg.PDU.FuncCode
	telega.PDU.CountReg = uint8(j)
	telega.PDU.Data = slicereg
	return telega

}
