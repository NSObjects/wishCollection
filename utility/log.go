package utility

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	ErrorLog   = "error.log"
	WarningLog = "warning.log"
	FatalLog   = "fatal.log"
	InfoLog    = "info.log"
	DebugLog   = "debug.log"
)

var (
	logs           []Log
	maxnumLogCache int
	eIP            string
	IsDebugModel   bool
)

const appName string = "CrawlerMain"

func init() {
	if ip, err := externalIP(); err == nil {
		eIP = ip
	}

	IsDebugModel = false
	maxnumLogCache = 1
}

func Debugln(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Println(file, line, DebugLog, args)
}

func Infoln(args ...interface{}) {

	_, file, line, _ := runtime.Caller(1)
	Println(file, line, InfoLog, args)
}

func Warningln(args ...interface{}) {

	_, file, line, _ := runtime.Caller(1)
	Println(file, line, WarningLog, args)
}

func Errorln(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Println(file, line, ErrorLog, args)
}

func Fatal(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Println(file, line, FatalLog, args)
	os.Exit(9)
}

func Println(file string, line int, errorType string, args ...interface{}) {

	var etype string
	switch errorType {
	case DebugLog:
		etype = "Debug"
	case ErrorLog:
		etype = "Error"
	case WarningLog:
		etype = "Warning"
	case InfoLog:
		etype = "Info"
	case FatalLog:
		etype = "Fatal"
	}
	msg := lastPath(file, "/")
	msg += ":"
	msg += fmt.Sprintf("%d", line)
	msg += " "
	msg += fmt.Sprintln(args...)

	if IsDebugModel {
		fmt.Println(msg)
	} else {
		log := Log{Message: msg, Ip: eIP}
		t := time.Now()

		log.Time = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d:%02d\n",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
		log.Type = etype
		log.AppName = appName
		logs = append(logs, log)
		if len(logs) >= maxnumLogCache {
			l := logs
			logs = nil
			go sendLog(l)
		}
	}

}

func sendLog(log []Log) {
	if len(log) <= 0 {
		return
	}

	logjson := LogJSON{Data: log}
	body, err := json.Marshal(&logjson)
	if err != nil {
		Errorln(err)
	}

	// Create client
	client := &http.Client{}

	// Create request
	//req, err := http.NewRequest("POST", "http://127.0.0.1:9528/api/log", bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", "http://108.61.162.82:9528/api/log", bytes.NewBuffer(body))
	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		fmt.Println("Failure : ", err)
	}
}

func ClearLog() {

	if len(logs) > 0 {
		sendLog(logs)
	}
	logs = nil
}

func lastPath(s string, sep string) string {
	ts := strings.Split(s, sep)
	return ts[len(ts)-1]
}

func PathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

type Log struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Ip      string `json:"ip"`
	AppName string `json:"app_name"`
}

type LogJSON struct {
	Data []Log `json:"data"`
}

const token string = "7437afe7d4b7db5eb1255d0e9ce75113f0357064c214f062493f8763d9b77862"

func SendDingdingMsg(msg DingDingMsg) {

	jsonMsg, err := json.Marshal(&msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	body := bytes.NewBuffer(jsonMsg)

	// Create client
	client := &http.Client{}

	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token)
	// Create request
	req, err := http.NewRequest("POST", url, body)

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}
	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}

func SendLog(msgs string) {

	t := time.Now()
	m := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d:%02d\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()) + msgs
	_, file, line, _ := runtime.Caller(1)
	content := fmt.Sprintf("%s:%s:%d:%s", appName, file, line, m)
	msg := DingDingMsg{Msgtype: "text", Text: Text{Content: content}}
	SendDingdingMsg(msg)
}

type DingDingMsg struct {
	Msgtype string `json:"msgtype"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}
