package nexpose

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
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
	Timeout: 10 * time.Second,
		Transport: tr,
	}

	return client
}

func getCreds() (string, string, error) {
	ssid := os.Getenv("SSID")
	sspw := os.Getenv("NEXPW")

	if ssid == "" || sspw == "" {
		return "", "", fmt.Errorf("no user or password supplied")
	}
	return ssid, sspw, nil
}

func getBaseURL() string {

	return "https://nexpose.secops.rackspace.com/api/3"

}

func TestNewService(t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)

	testCases := []struct {
		httpClient *http.Client
		baseURL    string
		user       string
		pass       string
		logger     *logrus.Logger
		want       string
	}{
		{
			httpClient: getHttpClient(),
			baseURL:    getBaseURL(),
			user:       ssid,
			pass:       sspw,
			logger:     getLogger(logrus.ErrorLevel),
			want:       "*nexpose.service",
		},
		{
			httpClient: getHttpClient(),
			baseURL:    getBaseURL(),
			user:       "",
			pass:       sspw,
			logger:     getLogger(logrus.ErrorLevel),
			want:       "<nil>",
		},
		{
			httpClient: getHttpClient(),
			baseURL:    getBaseURL(),
			user:       ssid,
			pass:       "",
			logger:     getLogger(logrus.ErrorLevel),
			want:       "<nil>",
		},
		{
			httpClient: getHttpClient(),
			baseURL:    "",
			user:       ssid,
			pass:       sspw,
			logger:     getLogger(logrus.ErrorLevel),
			want:       "<nil>",
		},
		{
			httpClient: nil,
			baseURL:    getBaseURL(),
			user:       ssid,
			pass:       sspw,
			logger:     getLogger(logrus.ErrorLevel),
			want:       "<nil>",
		},
		{
			httpClient: getHttpClient(),
			baseURL:    getBaseURL(),
			user:       ssid,
			pass:       sspw,
			logger:     nil,
			want:       "<nil>",
		},
		{
			httpClient: nil,
			baseURL:    "",
			user:       "",
			pass:       "",
			logger:     nil,
			want:       "<nil>",
		},
	}

	for _, tc := range testCases {

		nexpose, err := NewService(tc.httpClient, tc.baseURL, tc.user, tc.pass, tc.logger)

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v", tc.want, nil)
			} else {
				t.Errorf(
					"want:  %s   got:  %v", tc.want, nil)
			}
			continue

		}
		serviceType := reflect.TypeOf(nexpose)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}


func TestSites(t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)
	
			nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
			if err != nil{
				t.Fatal("Failed to get Nexpose Service: ", err.Error())
			}
			
			
			

	testCases := []struct {
		
		want       string
	}{
		{
			"*nexpose.Sites",
		},
	}

	for _, tc := range testCases {

		got, err := nexpose.Sites()

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			} else {
				
				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			}
			continue

		}
		serviceType := reflect.TypeOf(got)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}


func TestSite(t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)
	
			nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
			if err != nil{
				t.Fatal("Failed to get Nexpose Service: ", err.Error())
			}
			
			
			

	testCases := []struct {
		siteID  int32
		want       string
	}{
		{
			884,
			"*nexpose.Site",
		},
	}

	for _, tc := range testCases {

		got, err := nexpose.Site(tc.siteID)

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			} else {
				
				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			}
			continue

		}
		serviceType := reflect.TypeOf(got)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}


func TestSiteAssets(t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)
	
			nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
			if err != nil{
				t.Fatal("Failed to get Nexpose Service: ", err.Error())
			}
			
			
			

	testCases := []struct {
		siteID  int32
		want       string
	}{
		{
			884,
			"*nexpose.SiteAssets",
		},
	}

	for _, tc := range testCases {

		got, err := nexpose.SiteAssets(tc.siteID)

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			} else {
				
				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			}
			continue

		}
		serviceType := reflect.TypeOf(got)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}


func TestAssetVulnerabilities(t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)
	
			nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
			if err != nil{
				t.Fatal("Failed to get Nexpose Service: ", err.Error())
			}
			
			
			

	testCases := []struct {
		assetID  int64
		want       string
	}{
		{
			3576,
			"*nexpose.AssetVulnerabilities",
		},
	}

	for _, tc := range testCases {

		got, err := nexpose.AssetVulnerabilities(tc.assetID)

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			} else {
				
				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			}
			continue

		}
		serviceType := reflect.TypeOf(got)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}

func TestVulnerabilityExceptions(t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)
	
			nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
			if err != nil{
				t.Fatal("Failed to get Nexpose Service: ", err.Error())
			}
			
			
			

	testCases := []struct {
		want       string
	}{
		{
			
			"*nexpose.VulnExceptions",
		},
	}

	for _, tc := range testCases {

		got, err := nexpose.VulnerabilityExceptions()

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			} else {
				
				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			}
			continue

		}
		serviceType := reflect.TypeOf(got)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}


func TestVulnerabilityException (t *testing.T) {

	ssid, sspw, err := getCreds()
	must(err, t)
	
			nexpose, err := NewService(getHttpClient(), getBaseURL(), ssid, sspw, getLogger(logrus.DebugLevel))
			if err != nil{
				t.Fatal("Failed to get Nexpose Service: ", err.Error())
			}
			
			
			

	testCases := []struct {
		excepID int
		want       string
	}{
		{
			280250,
			"*nexpose.VulnException",
		},
	}

	for _, tc := range testCases {

		got, err := nexpose.VulnerabilityException(tc.excepID)

		if err != nil {

			if tc.want == "<nil>" {
				t.Logf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			} else {
				
				t.Errorf("want:  %s   got:  %v:  Error:  %s", tc.want, nil, err.Error())
			}
			continue

		}
		serviceType := reflect.TypeOf(got)

		t.Logf("want:  %s   got:  %s", tc.want, serviceType.String())

		if serviceType.String() != tc.want {
			t.Errorf("want:  %s  got:  %s", tc.want, serviceType.String())
		}
	}

}


func must(err error, t *testing.T) {
	if err != nil {
		t.Fatal("error != nil: ", err.Error())
	}

}
