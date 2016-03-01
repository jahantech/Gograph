package main

import (
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var Queue chan string
var TimerVal int

type ValuePost struct {
	Value string `json:"value"`
}

func init() {
	Queue = make(chan string)
	TimerVal = 0
}

func main() {

	dc := gg.NewContext(1000, 1000)

	go GraphUpdater(dc)

	http.HandleFunc("/Values", ValueInsert)
	if err := http.ListenAndServe(":8085", nil); err != nil {
		fmt.Println(err.Error())
	}

	return
}

func ValueInsert(rw http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(rw, http.StatusText(405), 405)
		return
	}

	var NewValue ValuePost
	err = json.Unmarshal(body, &NewValue)
	if err != nil {
		http.Error(rw, http.StatusText(405), 405)
		return
	}

	fmt.Println("Adding:" + NewValue.Value)
	Queue <- NewValue.Value

	return
}

func GraphUpdater(dc *gg.Context) {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case job := <-Queue:
			intVal, _ := strconv.Atoi(job)
			dc.DrawCircle(float64(TimerVal), float64(1000-(intVal)), 5)

			dc.Fill()
			dc.SavePNG("out.png")
			if TimerVal >= 1000 {
				TimerVal = 0
				dc.Clear()
			}
		case <-ticker.C:
			TimerVal = TimerVal + 1
		}
	}
}
