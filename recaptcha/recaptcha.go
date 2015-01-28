package recaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var verifyUrl string = "https://www.google.com/recaptcha/api/siteverify?secret=%s&response=%s"

type recaptchaResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// Checks the validity of a recaptcha response
func CheckRecaptcha(pvtKey string, response string) bool {
	requestUrl := fmt.Sprintf(verifyUrl, pvtKey, response)
	googleResponse, err := http.Get(requestUrl)
	if err != nil {
		return false
	}
	defer googleResponse.Body.Close()
	responseContent, err := ioutil.ReadAll(googleResponse.Body)
	if err != nil {
		return false
	}
	var responseJson recaptchaResponse
	json.Unmarshal(responseContent, &responseJson)
	return responseJson.Success
}
