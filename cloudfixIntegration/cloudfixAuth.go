package cloudfixIntegration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/trilogy-group/cloudfix-linter/logger"
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
	dlog := logger.DevLogger()
	username, present := os.LookupEnv("CLOUDFIX_USERNAME")
	if !present {
		dlog.Errorln("Couldn't find environment variable CLOUDFIX_USERNAME")
		return []byte{}, &customError{CRED_ERROR, "Could not login. Environment Variable CLOUDFIX_USERNAME not present."}
	}
	password, present := os.LookupEnv("CLOUDFIX_PASSWORD")
	if !present {
		dlog.Errorln("Couldn't find environment variable CLOUDFIX_PASSWORD")
		return []byte{}, &customError{CRED_ERROR, "Could not login. Environment Variable CLOUDFIX_PASSWORD not present."}
	}
	requestBody := Payload{password, username}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		dlog.Error(err)
		return []byte{}, &customError{GENERIC_ERROR, "Error logging user in"}
	}
	body := bytes.NewReader(requestBodyBytes)
	requestHTTP, err := http.NewRequest("POST", LOGIN_ENDPOINT, body)
	if err != nil {
		dlog.Error(err)
		return []byte{}, &customError{GENERIC_ERROR, "Could not make login request"}
	}
	requestHTTP.Header.Set("Accept", "application/json")
	requestHTTP.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(requestHTTP)
	if err != nil {
		dlog.Error(err)
		return []byte{}, &customError{GENERIC_ERROR, "Error logging user in"}
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadRequest { //Wrong Username and Passwords revert back a 400 response
		dlog.WithField("statusCode", response.StatusCode).Error("Bad login response code")
		return []byte{}, &customError{3, "Unauthorized login credentials. Please check username and password and try again."}
	} else if response.StatusCode != http.StatusOK {
		dlog.WithField("statusCode", response.StatusCode).Error("Bad login response code")
		return []byte{}, &customError{GENERIC_ERROR, "Error logging in"}
	}
	responeData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		dlog.Error(err)
		return []byte{}, &customError{GENERIC_ERROR, "Internal Error"}
	}
	return responeData, nil
}

func (ca *CloudfixAuth) storeCreds(loginParsed *LoginData) *customError {
	dlog := logger.DevLogger()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		dlog.Error(err)
		return &customError{GENERIC_ERROR, "Internal Error"}
	}
	homeDir += "/.cloudfix-creds"
	errDir := os.MkdirAll(homeDir, os.ModePerm)
	if errDir != nil {
		dlog.Error(errDir)
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	fileA, errC := os.Create(homeDir + "/access-token")
	if errC != nil {
		dlog.Error(errC)
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	_, errF := fileA.WriteString(loginParsed.AccessToken)
	if errF != nil {
		dlog.Error(errF)
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	fileR, errRT := os.Create(homeDir + "/refresh-token")
	if errRT != nil {
		dlog.Error(errRT)
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	_, errWRT := fileR.WriteString(loginParsed.RefreshToken)
	if errWRT != nil {
		dlog.Error(errWRT)
		return &customError{STORAGE_ERROR, "Internal Error"}
	}
	return nil
}

func (ca *CloudfixAuth) validate(token string) (int, *customError) {
	req, err := http.NewRequest("GET", ACCOUNTS_ENDPOINT, nil)
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
