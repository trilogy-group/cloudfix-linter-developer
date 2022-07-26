package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type CloudfixAuth struct {
	//no data fields required
}

type Payload struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Structure for unmarshalling the login data
type LoginData struct {
	AccessToken    string    `json:"accessToken"`
	RefreshToken   string    `json:"refreshToken"`
	ExpirationDate time.Time `json:"expirationDate"`
	TokenType      string    `json:"tokenType"`
}

//Member functions follow

func (ca *CloudfixAuth) loginAndSaveCreds() (string, *customError) {
	token, errL := ca.handleLogin()
	if errL != nil {
		return "", errL
	}
	loginParsed, errP := ca.parseLoginCreds(token)
	if errP != nil {
		return "", errP
	}
	errS := ca.storeCreds(loginParsed)
	if errS != nil {
		return loginParsed.AccessToken, errS
	}
	return loginParsed.AccessToken, nil
}

func (ca *CloudfixAuth) parseLoginCreds(loginData []byte) (*LoginData, *customError) {
	var loginParsed LoginData
	errJ := json.Unmarshal(loginData, &loginParsed)
	if errJ != nil {
		return &loginParsed, &customError{GENERIC_ERROR, "Internal Error"}
	}
	return &loginParsed, nil
}

func (ca *CloudfixAuth) getToken() (string, *customError) {
	homedir, errU := os.UserHomeDir()
	if errU != nil {
		return "", &customError{GENERIC_ERROR, "Can't access creds"}
	}
	homedir += "/.cloudfix-creds"
	accessTokenBytes, errO := ioutil.ReadFile(homedir + "/access-token")
	if errO != nil {
		// no file called access-token
		//need to login and save the access token
		return ca.loginAndSaveCreds()

	} else {
		accessToken := string(accessTokenBytes[:])
		statusCode, errV := ca.validate(accessToken)
		if errV != nil {
			return "", errV
		}
		switch statusCode {
		case 200:
			return accessToken, nil
		case 401:
			return ca.loginAndSaveCreds()
		default:
			return "", &customError{GENERIC_ERROR, "Internal Error"}
		}
	}
}

func (ca *CloudfixAuth) handleLogin() ([]byte, *customError) {
	username, present := os.LookupEnv("CLOUDFIX_USERNAME")
	if !present {
		return []byte{}, &customError{CRED_ERROR, "Error retreiving username. Setup environment correctly"}
	}
	password, present := os.LookupEnv("CLOUDFIX_PASSWORD")
	if !present {
		return []byte{}, &customError{CRED_ERROR, "Error retreiving password. Setup environment correctly"}
	}
	requestBody := Payload{password, username}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return []byte{}, &customError{GENERIC_ERROR, "Error logging user in"}
	}
	body := bytes.NewReader(requestBodyBytes)
	//response, err := http.Post("https://w9lnd111rl.execute-api.us-east-1.amazonaws.com/default/api/v1/auth/login", "application/json", body)
	requestHTTP, err := http.NewRequest("POST", "https://w9lnd111rl.execute-api.us-east-1.amazonaws.com/default/api/v1/auth/login", body)
	if err != nil {
		return []byte{}, &customError{GENERIC_ERROR, "Could not make login request"}
	}
	requestHTTP.Header.Set("Accept", "application/json")
	requestHTTP.Header.Set("Content-Type", "application/json")
	//fmt.Println(formatRequest(requestHTTP))
	response, err := http.DefaultClient.Do(requestHTTP)
	if err != nil {
		return []byte{}, &customError{GENERIC_ERROR, "Error logging user in"}
	}
	//defer response.Body.Close()
	if response.StatusCode == http.StatusBadRequest { //Wrong Username and Passwords revert back a 400 response
		return []byte{}, &customError{3, "Unauthorized login credentials. Please check username and password and try again."}
	} else if response.StatusCode != http.StatusOK {
		return []byte{}, &customError{GENERIC_ERROR, "Error logging in"}
	}
	responeData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, &customError{GENERIC_ERROR, "Internal Error"}
	}
	return responeData, nil
}

func (ca *CloudfixAuth) storeCreds(loginParsed *LoginData) *customError {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &customError{GENERIC_ERROR, "Internal Error"}
	}
	homeDir += "/.cloudfix-creds"
	fileA, errC := os.Create(homeDir + "/access-token")
	if errC != nil {
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	_, errF := fileA.WriteString(loginParsed.AccessToken)
	if errF != nil {
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	fileR, errRT := os.Create(homeDir + "/refresh-token")
	if errRT != nil {
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	_, errWRT := fileR.WriteString(loginParsed.RefreshToken)
	if errWRT != nil {
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	return nil
}

func (ca *CloudfixAuth) validate(token string) (int, *customError) {
	req, err := http.NewRequest("GET", "https://w9lnd111rl.execute-api.us-east-1.amazonaws.com/default/api/v1/financial/accounts", nil)
	if err != nil {
		return 0, &customError{GENERIC_ERROR, "Internal Error"}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, &customError{GENERIC_ERROR, "Internal Error"}
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
