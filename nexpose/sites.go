package nexpose

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"io/ioutil"
	"net/http"
)

type Sites struct {
	Links     []Links `json:"links"`
	Page      `json:"page"`
	Resources []Resources `json:"resources"`
}

type Links struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type Page struct {
	Number         int `json:"number"`
	Size           int `json:"size"`
	TotalPages     int `json:"totalPages"`
	TotalResources int `json:"totalResources"`
}

type Resources struct {
	Assets          int             `json:"assets"`
	ID              int             `json:"id"`
	Importance      string          `json:"importance"`
	LastScanTime    string          `json:"lastScanTime"`
	Links           []Links         `json:"links"`
	Name            string          `json:"name"`
	RiskScore       float64         `json:"riskScore"`
	ScanEngine      int             `json:"scanEngine"`
	ScanTemplate    string          `json:"scanTemplate"`
	Type            string          `json:"type"`
	Vulnerabilities Vulnerabilities `json:"vulnerabilities"`
}

type Vulnerabilities struct {
	Critical int `json:"critical"`
	Moderate int `json:"moderate"`
	Severe   int `json:"severe"`
	Total    int `json:"total"`
}

type Site struct {
	Assets          int             `json:"assets"`
	ConnectionType  string          `json:"connectionType"`
	Description     string          `json:"description"`
	ID              int32           `json:"id"`
	Importance      string          `json:"importance"`
	LastScanTime    string          `json:"lastScanTime"`
	Links           []Links         `json:"links"`
	Name            string          `json:"name"`
	RiskScore       float64         `json:"riskScore"`
	ScanEngine      int32           `json:"scanEngine"`
	ScanTemplate    string          `json:"scanTemplate"`
	Type            string          `json:"type"`
	Vulnerabilities Vulnerabilities `json:"vulnerabilities"`
}

type SiteAssets struct {
	Links     []Links              `json:"links"`
	Page      Page                 `json:"page"`
	Resources []SiteAssetResources `json:"resources"`
}

type Addresses struct {
	IP  string `json:"ip"`
	Mac string `json:"mac"`
}
type Configurations struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Databases struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}
type Attributes struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Files struct {
	Attributes []Attributes `json:"attributes"`
	Name       string       `json:"name"`
	Size       int          `json:"size"`
	Type       string       `json:"type"`
}
type History struct {
	Date                     time.Time `json:"date"`
	Description              string    `json:"description"`
	ScanID                   int       `json:"scanId"`
	Type                     string    `json:"type"`
	User                     string    `json:"user"`
	Version                  int       `json:"version"`
	VulnerabilityExceptionID string    `json:"vulnerabilityExceptionId"`
}
type HostNames struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}
type Ids struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}
type Cpe struct {
	Edition   string `json:"edition"`
	Language  string `json:"language"`
	Other     string `json:"other"`
	Part      string `json:"part"`
	Product   string `json:"product"`
	SwEdition string `json:"swEdition"`
	TargetHW  string `json:"targetHW"`
	TargetSW  string `json:"targetSW"`
	Update    string `json:"update"`
	V22       string `json:"v2.2"`
	V23       string `json:"v2.3"`
	Vendor    string `json:"vendor"`
	Version   string `json:"version"`
}
type OsFingerprint struct {
	Architecture   string           `json:"architecture"`
	Configurations []Configurations `json:"configurations"`
	Cpe            Cpe              `json:"cpe"`
	Description    string           `json:"description"`
	Family         string           `json:"family"`
	ID             int              `json:"id"`
	Product        string           `json:"product"`
	SystemName     string           `json:"systemName"`
	Type           string           `json:"type"`
	Vendor         string           `json:"vendor"`
	Version        string           `json:"version"`
}
type UserGroups struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Users struct {
	FullName string `json:"fullName"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}
type Pages struct {
	LinkType string `json:"linkType"`
	Path     string `json:"path"`
	Response int    `json:"response"`
}
type WebApplications struct {
	ID          int     `json:"id"`
	Pages       []Pages `json:"pages"`
	Root        string  `json:"root"`
	VirtualHost string  `json:"virtualHost"`
}
type Services struct {
	Configurations  []Configurations  `json:"configurations"`
	Databases       []Databases       `json:"databases"`
	Family          string            `json:"family"`
	Links           []Links           `json:"links"`
	Name            string            `json:"name"`
	Port            int               `json:"port"`
	Product         string            `json:"product"`
	Protocol        string            `json:"protocol"`
	UserGroups      []UserGroups      `json:"userGroups"`
	Users           []Users           `json:"users"`
	Vendor          string            `json:"vendor"`
	Version         string            `json:"version"`
	WebApplications []WebApplications `json:"webApplications"`
}
type Software struct {
	Configurations []Configurations `json:"configurations"`
	Cpe            Cpe              `json:"cpe"`
	Description    string           `json:"description"`
	Family         string           `json:"family"`
	ID             int              `json:"id"`
	Product        string           `json:"product"`
	Type           string           `json:"type"`
	Vendor         string           `json:"vendor"`
	Version        string           `json:"version"`
}

type SiteAssetResources struct {
	Addresses                  []Addresses      `json:"addresses"`
	AssessedForPolicies        bool             `json:"assessedForPolicies"`
	AssessedForVulnerabilities bool             `json:"assessedForVulnerabilities"`
	Configurations             []Configurations `json:"configurations"`
	Databases                  []Databases      `json:"databases"`
	Files                      []Files          `json:"files"`
	History                    []History        `json:"history"`
	HostName                   string           `json:"hostName"`
	HostNames                  []HostNames      `json:"hostNames"`
	ID                         int              `json:"id"`
	Ids                        []Ids            `json:"ids"`
	IP                         string           `json:"ip"`
	Links                      []Links          `json:"links"`
	Mac                        string           `json:"mac"`
	Os                         string           `json:"os"`
	OsFingerprint              OsFingerprint    `json:"osFingerprint"`
	RawRiskScore               float64          `json:"rawRiskScore"`
	RiskScore                  float64          `json:"riskScore"`
	Services                   []Services       `json:"services"`
	Software                   []Software       `json:"software"`
	Type                       string           `json:"type"`
	UserGroups                 []UserGroups     `json:"userGroups"`
	Users                      []Users          `json:"users"`
	Vulnerabilities            Vulnerabilities  `json:"vulnerabilities"`
}

func makeGetRequest(c *client, url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.user, c.pass)
	resp, err := c.httpClient.Do(req)
	log.Debug("Response Status: ", resp.StatusCode, resp.Status)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf("error: %v    Request: %+v   Response Body: %+v", err, req, resp))
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil

}

// Sites will  return details for all sites to which the account has privileges
func (c *client) Sites() (*Sites, error) {
	c.logger.WithFields(log.Fields{
		"func": "nexpose.Sites()",
	}).Debug("Calling nexpose to get all site info.")

	url := fmt.Sprintf("%s/sites", c.baseURL)

	body, err := makeGetRequest(c, url)
	if err != nil {
		return nil, err
	}

	sites := Sites{}
	err = json.Unmarshal(body, &sites)
	if err != nil {
		return nil, err
	}
	return &sites, nil
}

// Site will return details for the site that has the ID that is passed in
func (c *client) Site(siteID int32) (*Site, error) {

	c.logger.WithFields(log.Fields{
		"func": "nexpose.Sites()",
	}).Debug("Calling nexpose to get specific site info.")

	url := fmt.Sprintf("%s/sites/%d", c.baseURL, siteID)

	body, err := makeGetRequest(c, url)
	if err != nil {
		return nil, err
	}
	site := Site{}
	err = json.Unmarshal(body, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}
