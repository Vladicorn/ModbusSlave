package Modbus

//буфферная функция
func bufferChan(chanbool chan bool, chIn, chOut chan TelegramAnsver) {
	var bufval TelegramAnsver
	for {
		select {
		case bufval = <-chIn:
			//fmt.Println(bufval)
		case <-chanbool:
			//fmt.Println("Handle request")
			chOut <- bufval
		}
	}
}

//буфферная функция
func bufferChanClient(chanbool chan bool, chIn, chOut chan TelegramAnsverSlave) {
	var bufval TelegramAnsverSlave
	for {
		select {
		case bufval = <-chIn:
			//fmt.Println(bufval)
		case <-chanbool:
			chOut <- bufval
		}
	}
}
