package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const ApiUrl = "http://localhost:8080"
const ApiPort = ":8080"

func Request(meth string, path string, values map[string]string) ([]byte, error) {
	data := url.Values{}
	if values != nil {
		for k, v := range values {
			data.Set(k, v)
		}
	}
	u, _ := url.ParseRequestURI(ApiUrl)
	u.Path = path
	urlStr := u.String()
	client := &http.Client{}
	r, _ := http.NewRequest(meth, urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("Error when send http request: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error when read from response body: ", err)
		return nil, err
	}

	return body, nil
}

