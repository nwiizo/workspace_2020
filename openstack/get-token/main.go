package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Payload struct {
	Auth Auth `json:"auth"`
}
type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Auth struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
	TenantID            string              `json:"tenantId"`
}

func execute(resp *http.Response) {
	// response body ioutil.ReadAll
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println(string(b))
	}
}

func main() {
	data := Payload{Auth: Auth{PasswordCredentials: PasswordCredentials{Username: "*******", Password: "************"}, TenantID: "******************"}}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://*************/v2.0/tokens", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	execute(resp)
	defer resp.Body.Close()
}
