package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

var (
	addr = "0.0.0.0:42069"
)

var topicdata []string
var authtokens []string

func main() {
	topicdata = dataloader()
	fmt.Println(topicdata)
	h := mainrequestHandler

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func dataloader() []string {
	return []string{"1, 2, 3", "asdf"}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func mainrequestHandler(rqu *fasthttp.RequestCtx) {
	path := string(rqu.Path()[:])
	fmt.Println(path)
	if path[len(path)-1:] == "/" {
		path = strings.TrimSuffix(path, "/")
	}
	pathformated := strings.Split(path, "/")
	auth := string(rqu.Request.Header.Cookie("auth"))

	if stringInSlice(auth, authtokens) {
		if path == "/topics" && rqu.IsGet() { // Got it
			j, err := json.Marshal(topicdata)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				fmt.Println(string(j))
			}
			_, err2 := fmt.Fprintf(rqu, string(j))
			if err2 != nil {
				return
			}
			rqu.SetContentType("application/json; charset=utf8")
		} else if cap(pathformated) > 2 && pathformated[1] == "topic" {
			if stringInSlice(pathformated[2], topicdata) {
				rqu.SetContentType("text/plain")
				rqu.SetStatusCode(200)
				rqu.SetBodyString("true")
			}
		} else {
			// mark false tries...
			rqu.SetContentType("text/plain")
			rqu.SetStatusCode(404)
			rqu.SetBodyString("404")
		}
		rqu.Redirect("/topics", 200)
	} else {
		if path == "/login" {
		}
		//call microsoft graph api
		// mark false tries...
		rqu.SetContentType("text/plain")
		rqu.SetStatusCode(404)
		rqu.SetBodyString("404")
	}
	rqu.Response.Header.Set("Access-Control-Allow-Origin", "*")

}
