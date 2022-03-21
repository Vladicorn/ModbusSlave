package main

import (
	"ModbusSlave/pkg/Modbus"
	"fmt"
	"html/template"
	"net/http"
)

type ViewData struct {
	Title   uint16
	Message uint8
}

func (ch *TelStr) HomeStart(w http.ResponseWriter, r *http.Request) {
	if !ch.chFlagStart {
		go Modbus.ModbusConnect(ch.ch, ch.chansv, ch.chStop, ch.chStart)
		ch.chFlagStart = true
	}

	//ch.chStart <- true
	fmt.Println("енене")
	ch.chansv <- true

	teleg := <-ch.ch
	fmt.Println(teleg)
	data := ViewData{
		Title:   teleg.MBAP.TranID,
		Message: teleg.MBAP.UnitID,
	}
	/*	r.ParseForm()
		// logic part of log in
		a0 := r.Form["sliceID0"]

		if len(a0) > 0 {
			a01 := []byte(a0[0])
			fmt.Println(a01)
			Modbus.RegSlice[0] = a01[0]

		}*/

	tmpl, _ := template.ParseFiles("./html/Start.html")
	tmpl.Execute(w, data)

}

func (ch *TelStr) HomeStop(w http.ResponseWriter, r *http.Request) {

	ch.chStop <- true
	ch.chFlagStart = false
	ch.chansv <- true
	teleg := <-ch.ch
	data := ViewData{
		Title:   teleg.MBAP.TranID,
		Message: teleg.MBAP.UnitID,
	}
	r.ParseForm()
	// logic part of log in
	a0 := r.Form["sliceID0"]

	if len(a0) > 0 {
		a01 := []byte(a0[0])
		fmt.Println(a01)
		Modbus.RegSlice[0] = a01[0]

	}

	tmpl, _ := template.ParseFiles("./html/Stop.html")

	tmpl.Execute(w, data)

}

func (ch *TelStr) HomeIndex(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("./html/Index.html")
	data := ViewData{}
	tmpl.Execute(w, data)

}

type TelStr struct {
	ch          chan Modbus.TelegramAnsver
	chansv      chan bool
	chStop      chan bool
	chStart     chan bool
	chFlagStart bool
	chFlagStop  bool
}

func main() {

	//var telegram Modbus.TelegramAnsver
	var telStr TelStr

	ch := make(chan Modbus.TelegramAnsver)
	chansv := make(chan bool)
	chanCancel := make(chan bool)
	chanStart := make(chan bool)
	telStr.ch = ch
	telStr.chansv = chansv
	telStr.chStop = chanCancel
	telStr.chStart = chanStart
	telStr.chFlagStart = false
	fmt.Println("Server is listening...")
	go http.HandleFunc("/stop", telStr.HomeStop)
	go http.HandleFunc("/start", telStr.HomeStart)
	go http.HandleFunc("/", telStr.HomeIndex)

	http.ListenAndServe(":8181", nil)

	// Запускаем цикл

}
