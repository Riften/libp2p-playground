package api

import (
	"fmt"
	"net/http"
)

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error when parse request: ", r.RequestURI)
		fmt.Println(err)
		return
	}

	fmt.Println("Get request: "+r.Method)
	//fmt.Printf("Request: %s\n", r.RequestURI)
	fmt.Println("\t",r.URL.Path)
	fmt.Println("query:")
	for k, v := range r.URL.Query() {
		fmt.Println("\t", "key:", k, ", value:", v[0])
	}
	values:= r.PostForm
	fmt.Println("values:")
	for k, v := range values {
		fmt.Println("\t", "key:", k, ", value:", v[0])
	}

	switch r.URL.Path {
	case "/expr":
		responseString(w, "test")
	case "/listPeer":
		a.ListPeers(w)
	}
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("aaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccccccddddddddddddddddddddddddddeeeeeeeeeeeeeeeeeeeeeeeeeffffffffffffffffffffffffffggggggggggggggggggggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiijjjjjjjjjjjjjjjjjjjjjjjjjjjjkkkkkkkkkkkkkkkkkkkkkkkkkk" +
	//	"llllllllllllllllllllllmmmmmmmmmmmmmmmmmmmmnnnnnnnnnnnnnnnnnnnoooooooooooooooooooooppppppppppppppppqqqqqqqqqqqqqqqqqrrrrrrrrrrrrrrrrrsssssssssssssssssssstttttttttttttttttuuuuuuuuuuuuuuuuuvvvvvvvvvvvvvvvvvvvvvvwwwwwwwwwwww"))

	//<-context.Background().Done()
	//fmt.Println("Server end")
}

func responseString(w http.ResponseWriter, str string) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println("Error when write to http response: ", err)
	}
}

func (a *Api) ListPeers(w http.ResponseWriter) {
	peers := a.Node.Peers()
	for _, p := range peers {
		_, err := w.Write([]byte("\t" + p.Pretty()+"\n"))
		if err != nil {
			fmt.Println("Error when write response to")
		}
	}
	w.WriteHeader(http.StatusOK)
}
