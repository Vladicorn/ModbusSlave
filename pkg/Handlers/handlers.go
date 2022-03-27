package Handlers

import (
	"ModbusSlave/pkg/Modbus"
	"net/http"
	"text/template"
)

func (ch *TelStrClin) HomeStartCLient(w http.ResponseWriter, r *http.Request) {
	if !ch.ChFlagStart {
		go Modbus.ModsbusConnectClient(ch.Ch, ch.Chansv, ch.ChStop, ch.ChStart)
		ch.ChFlagStart = true
	}
	ch.Chansv <- true
	teleg := <-ch.Ch
	data := AllowSt{
		Title:   teleg.PDU.Data,
		Message: teleg.MBAP.UnitID,
	}
	tmpl, _ := template.ParseFiles("./html/StartClient.html")
	tmpl.Execute(w, data)
}

func (ch *TelStrClin) HomeStopClient(w http.ResponseWriter, r *http.Request) {
	ch.ChStop <- true
	ch.ChFlagStart = false
	tmpl, _ := template.ParseFiles("./html/Stop.html")
	tmpl.Execute(w, nil)
}

func (ch *TelStr) HomeStart(w http.ResponseWriter, r *http.Request) {
	if !ch.ChFlagStart {
		go Modbus.ModbusConnect(ch.Ch, ch.Chansv, ch.ChStop, ch.ChStart)
		ch.ChFlagStart = true
	}
	ch.Chansv <- true
	teleg := <-ch.Ch
	data := ViewData{
		Title:   teleg.MBAP.TranID,
		Message: teleg.MBAP.UnitID,
	}
	tmpl, _ := template.ParseFiles("./html/Start.html")
	tmpl.Execute(w, data)

}

func (ch *TelStr) HomeStop(w http.ResponseWriter, r *http.Request) {
	ch.ChStop <- true
	ch.ChFlagStart = false
	ch.Chansv <- true
	teleg := <-ch.Ch
	data := ViewData{
		Title:   teleg.MBAP.TranID,
		Message: teleg.MBAP.UnitID,
	}
	tmpl, _ := template.ParseFiles("./html/Stop.html")
	tmpl.Execute(w, data)

}

func (ch *TelStr) HomeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./html/Index.html")
	tmpl.Execute(w, nil)
}
