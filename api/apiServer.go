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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("aaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbccccccccccccccccccccccccccccddddddddddddddddddddddddddeeeeeeeeeeeeeeeeeeeeeeeeeffffffffffffffffffffffffffggggggggggggggggggggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiijjjjjjjjjjjjjjjjjjjjjjjjjjjjkkkkkkkkkkkkkkkkkkkkkkkkkk" +
		"llllllllllllllllllllllmmmmmmmmmmmmmmmmmmmmnnnnnnnnnnnnnnnnnnnoooooooooooooooooooooppppppppppppppppqqqqqqqqqqqqqqqqqrrrrrrrrrrrrrrrrrsssssssssssssssssssstttttttttttttttttuuuuuuuuuuuuuuuuuvvvvvvvvvvvvvvvvvvvvvvwwwwwwwwwwww"))
}
