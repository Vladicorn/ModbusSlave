package Handlers

import "ModbusSlave/pkg/Modbus"

type ViewData struct {
	Title   uint16
	Message uint8
}

type AllowSt struct {
	Title   []uint16
	Message uint8
}

type TelStr struct {
	Ch          chan Modbus.TelegramAnsver
	Chansv      chan bool
	ChStop      chan bool
	ChStart     chan bool
	ChFlagStart bool
}

type TelStrClin struct {
	Ch          chan Modbus.TelegramAnsverSlave
	Chansv      chan bool
	ChStop      chan bool
	ChStart     chan bool
	ChFlagStart bool
}
