package cmd

import (
	"fmt"
	"github.com/Riften/libp2p-playground/api"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"net/http"
	"os"
)

func peerInfo(copy bool, out string) error {
	res, err := api.Request(http.MethodGet, "/peer/info", nil)
	if err != nil {
		fmt.Println("Error when send request to http: ", err)
		return err
	}
	fmt.Println(string(res))
	if copy {
		clipboard.WriteAll(string(res))
		fmt.Println("Peer info has been written to clipboard.")
	}
	if out!="" {
		err := ioutil.WriteFile(out, res, os.ModePerm)
		if err != nil {
			fmt.Println("Error when write peer info to file ", out, ": ", err)
			return err
		}
	}
	return nil
}

func peerList(copy bool) error {
	res, err := api.Request(http.MethodGet, "/peer/list", nil)
	if err != nil {
		fmt.Println("Error when send request to http: ", err)
		return err
	}
	fmt.Println(string(res))
	if copy {
		clipboard.WriteAll(string(res))
		fmt.Println("Peer info has been written to clipboard.")
	}
	return nil
}

func peerConnect(peerId string, addr string) error {
	params := map[string]string {
		"id": peerId,
		"addr": addr,
	}
	res, err := api.Request(http.MethodPost, "/peer/connect", params)
	if err != nil {
		fmt.Println("Error when send request to http: ", err)
		return err
	}
	fmt.Println(string(res))
	return nil
}
