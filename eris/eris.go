package eris

import (
	"fmt"
	"github.com/sirupsen/logrus"

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
	idToken    string
	logger     *logrus.Logger
}

type SiteResponse struct{}

type Client interface {
	Systems() (*SystemResponse, error)
}

// New Service creates a new instance of a nexpose client that allows for several calls to the nexpose api.   The focus of this service is currently read only calls.
func NewService(httpClient *http.Client, baseURL string, token string, logger *log.Logger) (Service, error) {

	if httpClient == nil {
		return nil, fmt.Errorf("http.Client is nil")

	}
	if token == "" {
		logger.WithFields(log.Fields{
			"func":    "NewService",
			"baseURL": baseURL,
		}).Debug("Token is empty")
		return nil, fmt.Errorf("token is nil")
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

	c := &client{httpClient, baseURL, token, logger}

	return &service{Client: c}, nil

}
