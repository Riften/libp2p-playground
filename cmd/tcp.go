package cmd

import (
	"fmt"
	"github.com/Riften/libp2p-playground/api"
	"net/http"
	"strconv"
)

func tcpListen(port int) error {
	portStr := strconv.Itoa(port)
	params := map[string]string {
		"port": portStr,
	}
	res, err := api.Request(http.MethodPost, "/tcp/listen", params)
	if err != nil {
		fmt.Println("Error when request tcp listen: ", err)
		return err
	}
	fmt.Println(string(res))
	return nil
}

func tcpSend(ip string, port int) error {
	portStr := strconv.Itoa(port)
	params := map[string]string {
		"port": portStr,
		"ip": ip,
	}
	res, err := api.Request(http.MethodPost, "/tcp/send", params)
	if err != nil {
		fmt.Println("Error when request tcp send: ", err)
		return err
	}
	fmt.Println(string(res))
	return nil
}
