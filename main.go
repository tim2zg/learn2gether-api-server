package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	addr = "0.0.0.0:42069"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var topicdata []string
var authtokens []string
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var users []string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	users, err = readLines("user.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	topicdata = dataloader()

	//fmt.Println(topicdata)
	h := mainrequestHandler

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func readLines(path string) ([]string, error) { // Tanks https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
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

func randomstring(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func mainrequestHandler(rqu *fasthttp.RequestCtx) {
	path := string(rqu.Path()[:])
	//fmt.Println(path)
	if path[len(path)-1:] == "/" {
		path = strings.TrimSuffix(path, "/")
	}
	pathformated := strings.Split(path, "/")
	auth := string(rqu.Request.Header.Cookie("auth"))

	//fmt.Println("Curent Cookie: " + auth)
	//fmt.Println("All Cookies: ")
	//fmt.Println(authtokens)

	if stringInSlice(auth, authtokens) {
		//fmt.Printf("True")
		if path == "/topics" && rqu.IsGet() { // Got it
			j, err := json.Marshal(topicdata)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				//fmt.Println(string(j))
			}
			_, err2 := fmt.Fprintf(rqu, string(j))
			if err2 != nil {
				return
			}
			rqu.SetBodyString(string(j))
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
	} else {
		if path == "/login" {
			rqu.Redirect(getredicturl(), 200)
		} else if path == "/getthetocken" {
			code := string(rqu.QueryArgs().Peek("code"))
			name := request(gettoken(code))
			if stringInSlice(name, users) {
				var c fasthttp.Cookie
				rqu.SetContentType("text/plain")
				rqu.SetBodyString("OK")
				rqu.SetStatusCode(200)
				randomstringlol := randomstring(420)
				authtokens = append(authtokens, randomstringlol)
				c.SetMaxAge(3600000)
				c.SetSecure(true)
				c.SetKey("auth")
				c.SetValue(randomstringlol)
				rqu.Response.Header.SetCookie(&c)
			} else {
				rqu.SetBodyString("Not allowed!")
				rqu.SetStatusCode(430)
			}
		} else {
			// mark false tries...
			rqu.SetContentType("text/plain")
			rqu.SetStatusCode(404)
			rqu.SetBodyString("404")
		}
	}
	rqu.Response.Header.Set("Access-Control-Allow-Origin", "*")

}
