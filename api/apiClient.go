package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func SendRequest(path string, values map[string]string, port int) error{
	//testReq, err := http.NewRequest("POST", "/cmd", nil)
	apiUrl := fmt.Sprintf("%s:%d", localhost, port)
	//resource := "test"
	data := url.Values{}
	if values != nil {
		for k, v := range values {
			data.Set(k, v)
		}
	}

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = path

	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	buf := make([]byte, 1000)
	for {
		_, err := resp.Body.Read(buf)
		fmt.Print(buf)
		if err == io.EOF {
			fmt.Println("== eof ==")
			break
		} else if err != nil {
			fmt.Println("Error when read from response body: ", err)
		}
	}
	fmt.Println("")

	defer resp.Body.Close()

	return nil
}
