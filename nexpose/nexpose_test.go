package nexpose

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"reflect"
	"testing"
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

		Transport: tr,
	}

	return client
}

func getCreds() (string, string, error) {
	ssid := os.Getenv("SSID")
	sspw := os.Getenv("SSPW")

	if ssid == "" || sspw == "" {
		return "", "", fmt.Errorf("no user or password supplied")
	}
	return ssid, sspw, nil
}

func getBaseURL() string {

	return "https://nexpose.secops.rackspace.com/api/3"

}

func TestNewService(t *testing.T) {

	hc := getHttpClient()
	baseURL := getBaseURL()
	ssid, sspw, err := getCreds()
	must(err, t)
	logger := getLogger(logrus.DebugLevel)
	nexpose := NewService(hc, baseURL, ssid, sspw, logger)

	t.Log(reflect.TypeOf(nexpose))

	t.Error("somethinmg something")
}

func must(err error, t *testing.T) {
	if err != nil {
		t.Fatal("error != nil: ", err.Error())
	}

}
