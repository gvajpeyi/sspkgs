package nexpose

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

//ssl-only-weak-ciphers

type VulnExceptions struct {
	Links     []Links       `json:"links"`
	Page      Page          `json:"page"`
	Resources []VeResources `json:"resources"`
}

type Review struct {
	Comment string  `json:"comment"`
	Date    string  `json:"date"`
	Links   []Links `json:"links"`
	Name    string  `json:"name"`
	User    int     `json:"user"`
}
type Scope struct {
	ID            int     `json:"id"`
	Key           string  `json:"key"`
	Links         []Links `json:"links"`
	Port          int     `json:"port"`
	Type          string  `json:"type"`
	Vulnerability string  `json:"vulnerability"`
}
type Submit struct {
	Comment string  `json:"comment"`
	Date    string  `json:"date"`
	Links   []Links `json:"links"`
	Name    string  `json:"name"`
	Reason  string  `json:"reason"`
	User    int     `json:"user"`
}
type VeResources struct {
	Expires string  `json:"expires"`
	ID      int     `json:"id"`
	Links   []Links `json:"links"`
	Review  Review  `json:"review"`
	Scope   Scope   `json:"scope"`
	State   string  `json:"state"`
	Submit  Submit  `json:"submit"`
}

type VulnException struct {
	Expires string  `json:"expires"`
	ID      int     `json:"id"`
	Links   []Links `json:"links"`
	Review  Review  `json:"review"`
	Scope   Scope   `json:"scope"`
	State   string  `json:"state"`
	Submit  Submit  `json:"submit"`
}

// Sites will  return details for all sites to which the account has privileges
func (c *client) VulnerabilityExceptions() (*VulnExceptions, error) {
	c.logger.WithFields(log.Fields{
		"func": "nexpose.Sites()",
	}).Debug("Calling nexpose to get all site info.")

	url := fmt.Sprintf("%s/vulnerability_exceptions", c.baseURL)

	body, err := makeGetRequest(c, url)
	if err != nil {
		return nil, err
	}

	vulnExceptions := VulnExceptions{}
	err = json.Unmarshal(body, &vulnExceptions)
	if err != nil {
		return nil, err
	}
	return &vulnExceptions, nil
}

func (c *client) VulnerabilityException(exceptionID int) (*VulnException, error) {
	c.logger.WithFields(log.Fields{
		"func": "nexpose.Sites()",
	}).Debug("Calling nexpose to get all site info.")

	url := fmt.Sprintf("%s/vulnerability_exceptions/%d", c.baseURL, exceptionID)

	body, err := makeGetRequest(c, url)
	if err != nil {
		return nil, err
	}

	vulnException := VulnException{}
	err = json.Unmarshal(body, &vulnException)
	if err != nil {
		return nil, err
	}
	return &vulnException, nil
}
