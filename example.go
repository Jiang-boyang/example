package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./example <IP_ADDRESS> <PORT>")
		os.Exit(1)
	}
	ip := os.Args[1]
	portStr := os.Args[2]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid port number:", portStr)
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		SendAlert(fmt.Sprintf("Port %d on %s is not reachable", port, ip))
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Port %d on %s is reachable\n", port, ip)
}

func SendAlert(msg string) {
	b, _ := json.Marshal(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
	})

	resp := PublicHttpRequest(
		"POST",
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=73de9b90-87d2-4160-a541-091ec270b5e4",
		b,
		map[string]string{"Content-Type": "application/json"},
	)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func PublicHttpRequest(method, url string, values []byte, header map[string]string) *http.Response {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(values))) // URL-encoded payload
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for headKey, headValue := range header {
		req.Header.Add(headKey, headValue)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	return resp
}
