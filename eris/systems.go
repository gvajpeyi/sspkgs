package eris

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type SystemQuery struct {
	Filterby struct {
		Name struct {
			Regex   string `json:"$regex"`
			Options string `json:"$options"`
		} `json:"name"`
	} `json:"filterby"`
	Sortby struct {
		Name int `json:"name"`
	} `json:"sortby"`
	PageNumber int `json:"pageNumber"`
}

type SystemResponse struct {
	Items              []SystemItem `json:"items"`
	PageNumber         int          `json:"pageNumber"`
	TotalNumberOfItems int          `json:"totalNumberOfItems"`
}
type SystemItem struct {
	ID               string `json:"_id"`
	Name             string `json:"name"`
	V                int    `json:"__v"`
	IsPCI            bool   `json:"isPCI,omitempty"`
	ConfirmedHistory []struct {
	} `json:"-"`
	LastConfirmedList  string      `json:"lastConfirmedList,omitempty"`
	LastConfirmedDate  string      `json:"lastConfirmedDate,omitempty"`
	LastRequestDate    interface{} `json:"lastRequestDate,omitempty"`
	LastUpdateDate     string      `json:"lastUpdateDate,omitempty"`
	RecurringScanMonth string      `json:"recurringScanMonth,omitempty"`
	RecurringScanDate  string      `json:"recurringScanDate,omitempty"`
	Recurring          bool        `json:"recurring,omitempty"`
	RPN                string      `json:"RPN,omitempty,omitempty"`
	SecondarySME       string      `json:"secondarySME,omitempty"`
	SME                string      `json:"SME,omitempty"`
	SecondaryManager   string      `json:"secondaryManager,omitempty"`
	Manager            string      `json:"manager,omitempty"`
	SecondaryPOC       string      `json:"secondaryPOC,omitempty"`
	POC                string      `json:"POC,omitempty"`
	Description        string      `json:"description,omitempty"`
	SiteID             string      `json:"siteId,omitempty"`
}

type SystemNames struct {
	Systems []System `json:"systems"`
}
type System struct {
	Name         string `json:"name"`
	BusinessUnit string `json:"business_unit"`
}

func prettyPrint(i interface{}) string {

	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)

}

func (c *client) request(url string, payload string) ([]byte, error) {

	var jsonStr = []byte(payload)
	url = fmt.Sprintf("%s%s", c.baseURL, url)
	c.logger.Debugf("url: %v\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		c.logger.Debugf("error: %v\n", prettyPrint(err))

		return nil, fmt.Errorf("Error building request: %+v ", err.Error())
	}
	c.logger.Debug("idToken: ", c.idToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", fmt.Sprintf("%s", c.idToken))

	resp, err := c.httpClient.Do(req)

	c.logger.Debug("Response Status: ", resp.StatusCode, resp.Status)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Debugf("error: %v", prettyPrint(err))

		return nil, fmt.Errorf(fmt.Sprintf("error: %v    Request: %+v   Response Body: %+v", err, req, resp))
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil

	// return resp.Body, nil
	//
	// defer resp.Body.Close()
	//
	// body, _ := ioutil.ReadAll(resp.Body)
	// var systemResponse SystemResponse
	// err = json.Unmarshal(body, &systemResponse)
	//
	// if err != nil {
	// 	log.Printf("error decoding  response: %v", err)
	// 	if e, ok := err.(*json.SyntaxError); ok {
	// 		log.Errorf("syntax error at byte offset %d", e.Offset)
	// 	}
	// 	log.Errorf("body response: %q", body)
	// 	return nil, fmt.Errorf("Error Parsing Json: %+v ", err.Error())
	// }
	//
	// return systemResponse, nil
}

func (c *client) Systems() (*SystemResponse, error) {

	url := "/system/get"
	// var systemQuery SystemQuery
	//
	// systemQuery.Filterby.Name.Regex = fmt.Sprintf("", )
	// systemQuery.Filterby.Name.Options = "i"
	// systemQuery.Sortby.Name = 1
	// systemQuery.PageNumber = 0

	// payload, err := json.Marshal(systemQuery)
	// fmt.Printf("\n\nmarshalled: \n%s\n\n\n", payload)

	//payload := fmt.Sprintf(`{"filterby":"{}"}`)

	// start of request to eris
	// initialRequest := true
	// totalItems := 0
	// pageRequest := 0
	// receivedItems := 0

	payload := fmt.Sprintf(`{"filterby":"{\"name\":{\"$regex\": \"^((?!ASV-).)*$\", \"$options\": \"i\"}}","sortby":{"name":1}, "pageNumber":0}`)

	var systemResponse SystemResponse
	body, err := c.request(url, payload)
	if err != nil {
		log.Debugf("error: %v", prettyPrint(err))

		return nil, err
	}

	err = json.Unmarshal(body, &systemResponse)
	if err != nil {
		log.Debugf("error: %v", prettyPrint(err))

		return nil, err
	}
	return &systemResponse, nil

	//
	//
	// log.Debugf("erisPageData: %+v", erisPageData)
	// for x := 0; x < len(erisPageData.Items); x++ {
	// 	parsedData := strings.Split(erisPageData.Items[x].Name, " - ")
	// 	if len(parsedData) < 2 {
	// 		system := System{Name: parsedData[0], BusinessUnit: parsedData[0]}
	// 		systemNames.Systems = append(systemNames.Systems, system)
	//
	// 	} else {
	// 		system := System{Name: parsedData[1], BusinessUnit: parsedData[0]}
	// 		systemNames.Systems = append(systemNames.Systems, system)
	// 	}
	// }

	//
	//
	// totalItems = erisPageData.TotalNumberOfItems
	// receivedItems = len(systemNames.Systems)
	// log.Infof("receivedItems: %v     totalItems: %v", receivedItems, totalItems)
	// //more data to get.  Increment page and grab it

}
