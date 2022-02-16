package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	addr = "0.0.0.0:42069"
)

var topicdata map[string]interface{}

func main() {
	topicdata = dataloader()
	fmt.Println(topicdata)
	h := mainrequestHandler

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func dataloader() map[string]interface{} {
	// Open our jsonFile
	jsonFile, err := os.Open("data.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	err2 := json.Unmarshal([]byte(byteValue), &result)
	if err2 != nil {
		fmt.Println(err2)
	}
	return result
}

func mainrequestHandler(rqu *fasthttp.RequestCtx) {
	path := string(rqu.Path()[:])
	pathformated := strings.Split(string(rqu.Path()[:]), "/")
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
	} else if cap(pathformated) > 1 && pathformated[1] == "topic" {
		if _, ok := topicdata["topic index"]; ok {
			rqu.SetContentType("text/plain")
			rqu.SetStatusCode(200)
			fmt.Println(topicdata["topic index"])
			rqu.SetBodyString(pathformated[2])
		}
	} else if path == "/login" {
		// redirect to microsoft api call
	} else {
		// make false tries...
		rqu.SetContentType("text/plain")
		rqu.SetStatusCode(404)
		rqu.SetBodyString("404")
	}
	rqu.Response.Header.Set("Access-Control-Allow-Origin", "*")

}
