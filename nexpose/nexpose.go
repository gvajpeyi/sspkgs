package nexpose

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Service interface {
	Client
}

var _ Service = &service{}

type service struct {
	Client
}

var _ Client = &client{}

type client struct {
	httpClient *http.Client
	baseURL    string
	user    string
	pass    string
	logger     *log.Logger
}

type SiteResponse struct{}

type Client interface {
	Sites() (*Sites, error)
	Site(siteID int32) (*Site, error)
	SiteAssets(siteID int32) (*SiteAssets, error)
	AssetVulnerabilities(assetID int64) (*AssetVulnerabilities, error)
	VulnerabilityExceptions() (*VulnExceptions, error)
	VulnerabilityException(exceptionID int) (*VulnException, error)
}

// New Service creates a new instance of a nexpose client that allows for several calls to the nexpose api.   The focus of this service is currently read only calls.
func NewService(httpClient *http.Client, baseURL string, user string, pass string, logger *log.Logger) (Service, error) {

	if httpClient == nil {
		return nil, fmt.Errorf("http.Client is nil")

	}
	if pass == "" || user == "" {
		return nil, fmt.Errorf("user or pass is nil")
	}
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")

	}

	logger.WithFields(log.Fields{
		"func":    "NewService",
		"baseURL": baseURL,
	}).Debug("New Service Called")

	c := &client{httpClient, baseURL, user, pass, logger}

	return &service{Client: c}, nil

}

