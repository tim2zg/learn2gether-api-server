package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type data struct {
	OdataContext      string        `json:"@odata.context"`
	BusinessPhones    []interface{} `json:"businessPhones"`
	DisplayName       string        `json:"displayName"`
	GivenName         string        `json:"givenName"`
	JobTitle          string        `json:"jobTitle"`
	Mail              string        `json:"mail"`
	MobilePhone       interface{}   `json:"mobilePhone"`
	OfficeLocation    interface{}   `json:"officeLocation"`
	PreferredLanguage interface{}   `json:"preferredLanguage"`
	Surname           string        `json:"surname"`
	UserPrincipalName string        `json:"userPrincipalName"`
	ID                string        `json:"id"`
}

func getredicturl() string {
	cred, err := confidential.NewCredFromSecret(os.Getenv("secret"))
	if err != nil {
		fmt.Println("could not create a cred from a secret: %w", err)
	}
	confidentialClientApp, _ := confidential.New(os.Getenv("clientid"), cred, confidential.WithAuthority("https://login.microsoftonline.com/common"))
	url, err := confidentialClientApp.AuthCodeURL(context.Background(), os.Getenv("clientid"), "http://localhost:42069/getthetocken", []string{"User.Read"})
	if err != nil {
		return ""
	}
	// Redirecting to the URL we have received
	return url
}

func gettoken(code string) string {
	// Initializing the client credential
	cred, err := confidential.NewCredFromSecret(os.Getenv("secret"))
	if err != nil {
		fmt.Println("could not create a cred from a secret: %w", err)
	}
	confidentialClientApp, err := confidential.New(os.Getenv("clientid"), cred, confidential.WithAuthority("https://login.microsoftonline.com/common"))
	result, err := confidentialClientApp.AcquireTokenByAuthCode(context.Background(), code, "http://localhost:42069/getthetocken", []string{"User.Read"})
	if err != nil {
		log.Fatal(err)
	}
	// Prints the access token on the webpage
	//fmt.Println("Access token is " + result.AccessToken)
	return result.AccessToken
}

func request(tocken string) string {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		//Handle Error
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer " + tocken},
	}

	resp, err := client.Do(req)
	if err != nil {
		//Handle Error
	}
	var datafromrequest data

	if err != nil {
		log.Fatal(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	err2 := json.Unmarshal(body, &datafromrequest)
	if err2 != nil {
		log.Fatal(err2)
	}

	//fmt.Println(datafromrequest.UserPrincipalName)
	return datafromrequest.UserPrincipalName
}
