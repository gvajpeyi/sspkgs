package eris

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.rackspace.com/SegmentSupport/sspkgs/identity"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

func getLogger(lvl logrus.Level) *logrus.Logger {
	var logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = new(logrus.TextFormatter) //default

	logger.SetOutput(logger.Writer())
	logger.SetLevel(lvl)
	logger.Out = os.Stdout
	return logger
}

func getHttpClient() *http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	return client
}

func getToken() (*string, error) {
	ssid := os.Getenv("SSID")
	sspw := os.Getenv("SSPW")

	if ssid == "" || sspw == "" {
		return nil, fmt.Errorf("no user or password supplied")
	}

	logit := new(log.Logger)

	logit.Out = os.Stdout
	logit.Formatter = new(logrus.TextFormatter) //default
	logit.SetOutput(logit.Writer())

	ll := logrus.ErrorLevel
	logit.SetLevel(ll)
	logit.Out = os.Stdout

	client := getHttpClient()

	id := identity.NewIdentityService(client, "https://identity-internal.api.rackspacecloud.com/v2.0", logit)
	token, err := id.AuthenticateWithPass(ssid, sspw, "Rackspace")

	if err != nil {
		logit.Debugf("error: %v", prettyPrint(err))

		return nil, err
	}
	logit.Debugf("token: %v", token.Access.Token.ID)

	return &token.Access.Token.ID, nil
}

func getBaseURL() string {

	return "https://prod.eris.rackspace.net"

}

func TestNewService(t *testing.T) {

	token, err := getToken()
	must(err, t)

	testCases := []struct {
		httpClient *http.Client
		baseURL    string
		idToken    *string
		logger     *logrus.Logger
		want       string
	}{
		{
			httpClient: getHttpClient(),
			baseURL:    getBaseURL(),
			idToken:    token,
			logger:     getLogger(logrus.DebugLevel),
			want:       "*eris.service",
		},
	}

	for _, tc := range testCases {

		eris, err := NewService(tc.httpClient, tc.baseURL, *tc.idToken, tc.logger)

		if err != nil {

			t.Errorf(
				"want:  %s   got:  %v", tc.want, nil)

			continue

		}
		serviceType := reflect.TypeOf(eris)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}

func TestSystems(t *testing.T) {

	token, err := getToken()
	must(err, t)

	testCases := []struct {
		httpClient *http.Client
		baseURL    string
		idToken    string
		logger     *logrus.Logger
		want       string
	}{
		{
			httpClient: getHttpClient(),
			baseURL:    getBaseURL(),
			idToken:    *token,
			logger:     getLogger(logrus.DebugLevel),
			want:       "*eris.SystemResponse",
		},
	}

	for _, tc := range testCases {

		e, err := NewService(tc.httpClient, tc.baseURL, tc.idToken, tc.logger)
		if err != nil {
			log.Debugf("error: %v", prettyPrint(err))

			t.Errorf("unable to create eris service")
			continue
		}

		erisSystem, err := e.Systems()
		if err != nil {
			log.Debugf("error: %v", prettyPrint(err))

		}
		serviceType := reflect.TypeOf(erisSystem)

		if err != nil {
			log.Debugf("error: %v", prettyPrint(err))

			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())

		}
		b, err := json.MarshalIndent(erisSystem, "", "\t")
		if err != nil {
			log.Debugf("error: %v", prettyPrint(err))

		}
		t.Logf("want:  %s   got:  %s\n%+v\n", tc.want, serviceType.String(), string(b))

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}

// func TestSites(t *testing.T) {
// 	token, err := getToken()
// 	must(err, t)
//
// 	nexpose, err := NewService(getHttpClient(), getBaseURL(), *token, getLogger(logrus.DebugLevel))
// 	if err != nil {
// 		t.Fatal("Failed to get Nexpose Service: ", err.Error())
// 	}
//
// 	testCases := []struct {
// 		want string
// 	}{
// 		{
// 			"*nexpose.Sites",
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		got, err := nexpose.Sites()
//
// 		if err != nil {
//
// 			if tc.want == "<nil>" {
// 				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			} else {
//
// 				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			}
// 			continue
//
// 		}
// 		serviceType := reflect.TypeOf(got)
//
// 		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())
//
// 		if serviceType.String() != tc.want {
// 			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
// 		}
// 	}
//
// }
//
// func TestSite(t *testing.T) {
//
// 	ssid, sspw, err := getCreds()
// 	must(err, t)
//
// 	nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
// 	if err != nil {
// 		t.Fatal("Failed to get Nexpose Service: ", err.Error())
// 	}
//
// 	testCases := []struct {
// 		siteID int32
// 		want   string
// 	}{
// 		{
// 			884,
// 			"*nexpose.Site",
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		got, err := nexpose.Site(tc.siteID)
//
// 		if err != nil {
//
// 			if tc.want == "<nil>" {
// 				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			} else {
//
// 				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			}
// 			continue
//
// 		}
// 		serviceType := reflect.TypeOf(got)
//
// 		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())
//
// 		if serviceType.String() != tc.want {
// 			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
// 		}
// 	}
//
// }
//
// func TestSiteAssets(t *testing.T) {
//
// 	ssid, sspw, err := getCreds()
// 	must(err, t)
//
// 	nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
// 	if err != nil {
// 		t.Fatal("Failed to get Nexpose Service: ", err.Error())
// 	}
//
// 	testCases := []struct {
// 		siteID int32
// 		want   string
// 	}{
// 		{
// 			884,
// 			"*nexpose.SiteAssets",
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		got, err := nexpose.SiteAssets(tc.siteID)
//
// 		if err != nil {
//
// 			if tc.want == "<nil>" {
// 				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			} else {
//
// 				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			}
// 			continue
//
// 		}
// 		serviceType := reflect.TypeOf(got)
//
// 		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())
//
// 		if serviceType.String() != tc.want {
// 			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
// 		}
// 	}
//
// }
//
// func TestAssetVulnerabilities(t *testing.T) {
//
// 	ssid, sspw, err := getCreds()
// 	must(err, t)
//
// 	nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
// 	if err != nil {
// 		t.Fatal("Failed to get Nexpose Service: ", err.Error())
// 	}
//
// 	testCases := []struct {
// 		assetID int64
// 		want    string
// 	}{
// 		{
// 			3576,
// 			"*nexpose.AssetVulnerabilities",
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		got, err := nexpose.AssetVulnerabilities(tc.assetID)
//
// 		if err != nil {
//
// 			if tc.want == "<nil>" {
// 				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			} else {
//
// 				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			}
// 			continue
//
// 		}
// 		serviceType := reflect.TypeOf(got)
//
// 		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())
//
// 		if serviceType.String() != tc.want {
// 			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
// 		}
// 	}
//
// }
//
// func TestVulnerabilityExceptions(t *testing.T) {
//
// 	ssid, sspw, err := getCreds()
// 	must(err, t)
//
// 	nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
// 	if err != nil {
// 		t.Fatal("Failed to get Nexpose Service: ", err.Error())
// 	}
//
// 	testCases := []struct {
// 		want string
// 	}{
// 		{
//
// 			"*nexpose.VulnExceptions",
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		got, err := nexpose.VulnerabilityExceptions()
//
// 		if err != nil {
//
// 			if tc.want == "<nil>" {
// 				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			} else {
//
// 				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			}
// 			continue
//
// 		}
// 		serviceType := reflect.TypeOf(got)
//
// 		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())
//
// 		if serviceType.String() != tc.want {
// 			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
// 		}
// 	}
//
// }
//
// func TestVulnerabilityException(t *testing.T) {
//
// 	ssid, sspw, err := getCreds()
// 	must(err, t)
//
// 	nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
// 	if err != nil {
// 		t.Fatal("Failed to get Nexpose Service: ", err.Error())
// 	}
//
// 	testCases := []struct {
// 		excepID int
// 		want    string
// 	}{
// 		{
// 			280250,
// 			"*nexpose.VulnException",
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		got, err := nexpose.VulnerabilityException(tc.excepID)
//
// 		if err != nil {
//
// 			if tc.want == "<nil>" {
// 				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			} else {
//
// 				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
// 			}
// 			continue
//
// 		}
// 		serviceType := reflect.TypeOf(got)
//
// 		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())
//
// 		if serviceType.String() != tc.want {
// 			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
// 		}
// 	}
//
// }

func must(err error, t *testing.T) {
	if err != nil {
		t.Fatal("error != nil: ", err.Error())
	}

}
