package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// прием лога
type Logg struct {
	Who     string `json:"who"`
	Where   string `json:"where"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func (lg *Logg) String() string {
	var lev string
	if lg.Level == "0" {
		lev = "info"
	} else {
		lev = "error"
	}
	return "@SERVICE@: " + lg.Who + ". @LEVEL@: " + lev + ". @DATE@: " + time.Now().Format(time.RFC850) + ". @MESSAGE@: " + lg.Message + ".\n"
}

func logger(w http.ResponseWriter, r *http.Request) {
	var lg Logg
	err := json.NewDecoder(r.Body).Decode(&lg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lg.save_log()
	w.Write([]byte("0"))
}

func (lg *Logg) save_log() {
	arr_path := strings.Split(lg.Where, string(filepath.Separator))
	file_name := arr_path[len(arr_path)-1]
	folder_path := strings.Split(lg.Where, string(filepath.Separator)+file_name)[0]

	os.MkdirAll(folder_path, os.ModePerm)
	file, err := os.OpenFile(lg.Where, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("невозможно открыть файл. err: ", err)
	}
	defer file.Close()

	_, err = file.WriteString(lg.String())
	if err != nil {
		log.Fatal("невозможно записать в файл. err: ", err)
	}
}

// останов сервера
var server_is_stop bool = false

func is_stop(timeout int) {
	for {
		if server_is_stop {
			os.Exit(0)
		}
		server_is_stop = true
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

func live(w http.ResponseWriter, r *http.Request) {
	server_is_stop = false
	w.Write([]byte("0"))
}

// точка входа
func main() {
	timeout := flag.Int("timeout", 30, "время которе будет ожидать сервис без подключения сторонних программ")
	port := flag.Int("port", 5500, "порт, на котором будет подкниматься сервис")
	flag.Parse()

	go is_stop(*timeout)

	log.Println("Прослушиваение http://127.0.0.1:" + strconv.Itoa(*port) + "/; timeout = " + strconv.Itoa(*timeout))
	http.HandleFunc("/log", logger)
	http.HandleFunc("/live", live)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*port), nil))
}
