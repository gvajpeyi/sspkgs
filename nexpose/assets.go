package nexpose

import (
	"encoding/json"
	"fmt"
			log "github.com/sirupsen/logrus"

	
)

type AssetVulnerabilities struct {
	Links     []Links     `json:"links"`
	Page      Page        `json:"page"`
	Resources []AVResources `json:"resources"`
}

type AVResults struct {
	CheckID    string  `json:"checkId"`
	Exceptions []int   `json:"exceptions"`
	Key        string  `json:"key"`
	Links      []Links `json:"links"`
	Port       int     `json:"port"`
	Proof      string  `json:"proof"`
	Protocol   string  `json:"protocol"`
	Status     string  `json:"status"`
}
type AVResources struct {
	ID        string    `json:"id"`
	Instances int       `json:"instances"`
	Links     []Links   `json:"links"`
	Results   []AVResults `json:"results"`
	Status    string    `json:"status"`
}



// SiteAssets will return details for all assets associated with the site that has the ID that is passed in
func (c *client) SiteAssets(siteID int32) (*SiteAssets, error) {

		c.logger.WithFields(log.Fields{
		"func": "nexpose.Sites()",
	}).Debug("Calling nexpose to get all site info.")
	
	path:=fmt.Sprintf("/sites/%d/assets", siteID)
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	body, err := makeGetRequest(c, url)
	if err != nil{
		return nil, err
	}
	
	siteAssets := SiteAssets{}
	err = json.Unmarshal(body, &siteAssets)
	if err != nil{
		return nil, err
	}
	return &siteAssets, nil
}





func (c *client) AssetVulnerabilities(assetID int64)(*AssetVulnerabilities, error){
	
	path:=fmt.Sprintf("/assets/%d/vulnerabilities", assetID)
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	
	body, err := makeGetRequest(c, url)
	if err != nil{
		return nil, err
	}
	siteAssets := AssetVulnerabilities{}
	err = json.Unmarshal(body, &siteAssets)
	if err != nil{
		return nil, err
	}
	return &siteAssets, nil
}