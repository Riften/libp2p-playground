package cmd

import (
	"fmt"
	"github.com/Riften/libp2p-playground/api"
	"log"
	"net/http"
)

func speedSend(peer string) error{
	params := map[string]string {
		"peer": peer,
	}
	res, err := api.Request(http.MethodPost, "/speed/send", params)
	if err != nil {
		log.Println("Error when send request to http: ", err)
		return err
	}
	fmt.Println(string(res))
	return nil
}
