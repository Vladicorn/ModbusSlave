package Modbus

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

var RegSlice = []byte{
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
