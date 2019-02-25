package identity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	
	log "github.com/sirupsen/logrus"
)


//https: // identity-internal.api.rackspacecloud.com/v2.0/tokens

type Response struct {
	Access struct {
		ServiceCatalog []interface{} `json:"serviceCatalog"`
		User           struct {
			RAXAUTHDefaultRegion string `json:"RAX-AUTH:defaultRegion"`
			Roles                []struct {
				Name string `json:"name"`
				ID   string `json:"id,omitempty"`
			} `json:"roles"`
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"user"`
		Token struct {
			Expires                time.Time `json:"expires"`
			RAXAUTHIssued          time.Time `json:"RAX-AUTH:issued"`
			RAXAUTHAuthenticatedBy []string  `json:"RAX-AUTH:authenticatedBy"`
			ID                     string    `json:"id"`
		} `json:"token"`
	} `json:"access"`
}

type Request struct {
	Auth struct {
		RAXAUTHDomain struct {
			Name string `json:"name"`
		} `json:"RAX-AUTH:domain"`
		PasswordCredentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"passwordCredentials"`
	} `json:"auth"`
}

type FailedRequest struct{
BadRequest struct{
	Code int `json:"code"`
	Message string `json:"message"`
} `json:"badRequest"`
}


type IdentityClient interface {
	Authenticate(req Request)(Response,error)
	AuthenticateWithPass(user string, pass string, domain string) (Response, error)
	VerifyToken(token string) (bool, *Response, error)
	AuthIfInvalidToken(reqPayload Request, token string) (bool, *Response, error)

}


type IdentityService interface {
	IdentityClient
}

var _ IdentityService = &identityService{}

type identityService struct {
	IdentityClient
}

var _ IdentityClient = &identityClient{}

type identityClient struct {
	client  *http.Client
	baseURL string
	logger *log.Logger
}



func NewIdentityService(client *http.Client, baseURL string, logger *log.Logger ) (IdentityService){
		logger.WithFields(log.Fields{
		"func": "NewIdentityService",
		"baseURL": baseURL,
		
	}).Info("")
	
	id := &identityClient{ client, baseURL, logger}
	return &identityService{IdentityClient: id}

}


func (id *identityClient) VerifyToken(token string) (bool, *Response, error) {

	id.logger.WithFields(log.Fields{
		"func": "id.verifyToken",
		"token": token,
	}).Info("")

	qUrl := fmt.Sprintf("%s/%s/%s", id.baseURL, "tokens", token)

	req, err := http.NewRequest("GET", qUrl, nil)
	req.Header.Set("X-Auth-Token", token)

	resp, err := id.client.Do(req)
	if err != nil {
		return false, &Response{}, fmt.Errorf(fmt.Sprintf("error: %v    Response Body: %+v", err, resp))
	}

	body, err := ioutil.ReadAll(resp.Body)

	idResponse := Response{}
	err = json.Unmarshal(body, &idResponse)
	if err != nil {
		return false, &Response{}, fmt.Errorf(fmt.Sprintf("error: %v    Response Body: %+v", err, resp))
	}

	if resp.StatusCode == 200 {
		return true, &idResponse, nil
	}

	return false, &Response{}, nil

}

func (id *identityClient) AuthenticateWithPass(user string, pass string, domain string) (Response, error) {

	authReq := Request{}
	authReq.Auth.RAXAUTHDomain.Name = domain
	authReq.Auth.PasswordCredentials.Username = user
	authReq.Auth.PasswordCredentials.Password = pass

		id.logger.WithFields(log.Fields{
		"func": "id.AuthenticateWithPass",
		"id": user,
		"password not null": pass != "",
	}).Info("")

	payload, err := json.Marshal(authReq)
	
	
	if err != nil {
		id.logger.WithFields(log.Fields{
		"func": "id.AuthenticateWithPass",
		"id": user,
		"failing code": "payload, err := json.Marshal(authReq)",
		"authReq": authReq,
	}).Error(err)
		return Response{}, err
	}

	qUrl := fmt.Sprintf("%s/%s", id.baseURL, "tokens")

	req, err := http.NewRequest("POST", qUrl, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := id.client.Do(req)
	if err != nil {
		
		id.logger.WithFields(log.Fields{
		"func": "id.AuthenticateWithPass",
		"id": user,
		"failing code": "	resp, err := id.client.Do(req)",
		"req": req,
	}).Error(err)
		
		return Response{}, err
	}
	if resp.StatusCode == 400 {

		body, err := ioutil.ReadAll(resp.Body)
		failedReq := FailedRequest{}
		err = json.Unmarshal(body, &failedReq)

		if err == nil {
			err = errors.New(failedReq.BadRequest.Message)

		id.logger.WithFields(log.Fields{
		"func": "id.AuthenticateWithPass",
		"id": user,
		"failing code": "err = json.Unmarshal(body, &failedReq)",
		"failedReq": failedReq,
	}).Error(err)
		
		}

		return Response{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	idResponse := Response{}

	err = json.Unmarshal(body, &idResponse)
	if err != nil {
		
				id.logger.WithFields(log.Fields{
		"func": "id.AuthenticateWithPass",
		"id": user,
		"failing code": "err = json.Unmarshal(body, &idResponse)",
		"idResponse": idResponse,
	}).Error(err)
		
		return Response{}, err
	}
	return idResponse, nil

}


func (id *identityClient) AuthIfInvalidToken(reqPayload Request, token string)(bool, *Response, error){

	// if existing token is valid, return true and response struct

	isValidToken, idResponse, err := id.VerifyToken(token)
	if isValidToken {
		return true, idResponse, nil
	}

	if err != nil{
		return false, nil, err
	}


	// if existing token is not valid, authenticate and return bool and resp struct


	authResponse, err := id.Authenticate(reqPayload)
	if err != nil{
		 return false, nil, err
	}

	return true, &authResponse, nil



}

// Authenticate accepts an ID Request
func (id *identityClient) Authenticate(reqPayload Request) (Response,  error) {
	

	payload, err := json.Marshal(reqPayload)
	if err != nil {
		return Response{}, err
	}
	
	
	

	qUrl := fmt.Sprintf("%s/%s", id.baseURL, "tokens")


	
	req, err := http.NewRequest("POST", qUrl, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")


	
	resp, err := id.client.Do(req)
	if err != nil {
		
		
				
				id.logger.WithFields(log.Fields{
		"func": "id.Authenticate",
		"failing code": "resp, err := id.client.Do(req)",
		"req": req,
	}).Error(err)
		
		return Response{}, err
	}
	if resp.StatusCode == 400 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
				id.logger.WithFields(log.Fields{
		"func": "id.Authenticate",
		"failing code": "		body, err := ioutil.ReadAll(resp.Body)",
		"resp.Body": resp.Body,
	}).Error(err)
			return Response{}, err
		}

		failedReq := FailedRequest{}
		err = json.Unmarshal(body, &failedReq)
		if err != nil{
				
				id.logger.WithFields(log.Fields{
		"func": "id.Authenticate",
		"failing code": "err = json.Unmarshal(body, &failedReq)",
		"failedReq": failedReq,
	}).Error(err)
			return Response{}, err
		}


		return Response{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	
	idResponse := Response{}

	err = json.Unmarshal(body, &idResponse)
	if err != nil {
			
				id.logger.WithFields(log.Fields{
		"func": "id.Authenticate",
		"failing code": "err = json.Unmarshal(body, &idResponse)",
		"idResponse": idResponse,
	}).Error(err)
		return Response{}, err
	}
	return idResponse, nil

}
