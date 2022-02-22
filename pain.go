package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"io/ioutil"
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

func getredicturl() (string, error) {
	cred, err1 := confidential.NewCredFromSecret(os.Getenv("secret"))
	if err1 != nil {
		fmt.Println("could not create a cred from a secret: %w", err1)
		return "", err1
	}
	confidentialClientApp, err2 := confidential.New(os.Getenv("clientid"), cred, confidential.WithAuthority("https://login.microsoftonline.com/common"))
	if err2 != nil {
		return "", err2
	}
	url, err3 := confidentialClientApp.AuthCodeURL(context.Background(), os.Getenv("clientid"), "http://localhost:42069/getthetocken", []string{"User.Read"})
	if err3 != nil {
		return "", err3
	}
	// Redirecting to the URL we have received
	return url, nil
}

func gettoken(code string) (string, error) {
	// Initializing the client credential
	cred, err := confidential.NewCredFromSecret(os.Getenv("secret"))
	if err != nil {
		fmt.Println("could not create a cred from a secret: %w", err)
		return "", err
	}
	confidentialClientApp, err2 := confidential.New(os.Getenv("clientid"), cred, confidential.WithAuthority("https://login.microsoftonline.com/common"))
	if err2 != nil {
		return "", err2
	}
	result, err3 := confidentialClientApp.AcquireTokenByAuthCode(context.Background(), code, "http://localhost:42069/getthetocken", []string{"User.Read"})
	if err3 != nil {
		return "", err3
	}
	// Prints the access token on the webpage
	//fmt.Println("Access token is " + result.AccessToken)
	return result.AccessToken, nil
}

func request(tocken string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer " + tocken},
	}

	resp, err2 := client.Do(req)
	if err2 != nil {
		return "", err2
	}
	var datafromrequest data

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}

	err3 := json.Unmarshal(body, &datafromrequest)
	if err3 != nil {
		return "", err3
	}

	//fmt.Println(datafromrequest.UserPrincipalName)
	return datafromrequest.UserPrincipalName, nil
}
