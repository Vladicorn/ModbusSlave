package Modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PDUAnsverSlave struct {
	FuncCode uint8
	CountReg uint8
	Data     []uint16
}
type TelegramAnsverSlave struct {
	MBAP MBAPHeader
	PDU  PDUAnsverSlave
}

var RegSliceClient = []byte{
	0x00, 0x01,
	0x00, 0x00,
	0x00, 0x06,
	0x01,
	0x03,
	0x00, 0x00,
	0x00, 0x00,
}

func ReadHoldingRegisterPol(data []byte) TelegramAnsverSlave {
	var telegram TelegramAnsverSlave
	conn := bytes.NewBuffer(data)
	//чтение хеадер
	buf := make([]byte, 7)

	if count, err := conn.Read(buf); err != nil {
		fmt.Println(err)
		return telegram
	} else {
		if count != 7 {
			fmt.Println("More 7")
			return telegram
		}
	}

	//обработка хедера
	telegram.MBAP.TranID = binary.BigEndian.Uint16(buf[:2])
	telegram.MBAP.ProtocolID = binary.BigEndian.Uint16(buf[2:4])
	telegram.MBAP.Length = binary.BigEndian.Uint16(buf[4:6])
	telegram.MBAP.UnitID = buf[6]

	//обработка PDU
	buf = make([]byte, telegram.MBAP.Length-1)
	if _, err := conn.Read(buf); err != nil {
		fmt.Println(err)
		return telegram
	}
	telegram.PDU.FuncCode = buf[0]
	//	telegram.PDU.FirstReg = binary.BigEndian.Uint16(buf[1:3])
	//	telegram.PDU.CountReg = binary.BigEndian.Uint16(buf[3:5])
	telegram.PDU.CountReg = buf[1]
	teglen := int(telegram.MBAP.Length - 1)
	slicereg := make([]byte, teglen, 100)
	slicereg = buf[2:teglen]
	slicedata := make([]uint16, (teglen-2)/2, 100)
	for i := 0; i < (teglen-2)/2; i++ {
		slicedata[i] = binary.BigEndian.Uint16(slicereg[i*2 : i*2+2])
	}
	telegram.PDU.Data = slicedata
	return telegram
}
