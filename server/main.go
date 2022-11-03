package main

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)
var mySigningKey = []byte("mysupersecretphrase")
func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret information")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return  nil, fmt.Errorf("there was an error")
				}
				return mySigningKey, nil
			})

			if err!=nil{
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid{
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func handleRequests() {
	// before creating isAuthorized function
	// http.HandleFunc("/", homepage)

	http.Handle("/", isAuthorized(homepage))

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	fmt.Println("my simple server")
	handleRequests()
}
