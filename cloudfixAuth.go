package main

type CloudfixAuth struct {
	//no data fields required
}

//Member functions follow

func (ca *CloudfixAuth) getToken() ([]byte, error) {
	var token []byte
	return token, nil
}

func (ca *CloudfixAuth) handleLogin() ([]byte, error) {
	var token []byte
	return token, nil
}

func (ca *CloudfixAuth) storeCreds() error {
	return nil
}

func (ca *CloudfixAuth) validate(token []byte) (bool, error) {
	return true, nil
}
