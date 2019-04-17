package act

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"github.rackspace.com/SegmentSupport/sspkgs/coreods"
	"github.rackspace.com/SegmentSupport/sspkgs/identity"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var boolPtr *bool

func actServiceSetup() (*actClient, error) {

	var actURL = "https://act-api.gscs.rackspace.com/v1/"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{

		Transport: tr,
	}

	odsUsername := os.Getenv("ODSID")
	odsPassword := os.Getenv("ODSPW")
	odsHost := os.Getenv("ODSHOST")
	odsPort, nil := strconv.Atoi(os.Getenv("ODSPORT"))

	var logger = logrus.New()

	logger.Out = os.Stdout

	logger.Formatter = new(logrus.TextFormatter) //default
	logger.SetOutput(logger.Writer())
	//logger.SetReportCaller(true)

	loglevel := os.Getenv("LOGLEVEL")
	ll := logger.Level
	switch loglevel {
	case "debug":
		ll = logrus.DebugLevel
	case "info":
		ll = logrus.InfoLevel
	case "error":
		ll = logrus.ErrorLevel
	case "fatal":
		ll = logrus.FatalLevel
	default:
		ll = logrus.WarnLevel
	}

	logger.SetLevel(ll)

	logger.Out = os.Stdout

	odsConfig := coreods.ODSConfig{
		Host:     odsHost,
		Port:     odsPort,
		User:     odsUsername,
		Password: odsPassword,
		DBName:   "",
	}
	ods, err := coreods.NewODSService(odsConfig)
	if err != nil {
		log.Fatalf("unable to setup service: %v", err)
	}

	ssUser := os.Getenv("SSID")
	ssPass := os.Getenv("SSPW")

	id := identity.NewIdentityService(client, "https://identity-internal.api.rackspacecloud.com/v2.0", logger)
	tempToken, err := id.AuthenticateWithPass(ssUser, ssPass, "Rackspace")

	internalToken := &tempToken.Access.Token.ID
	act := &actClient{client, actURL, ssUser, ssPass, "Rackspace", ods, id, logger, internalToken}

	return act, err

}

func TestResponsibleDepartments(t *testing.T) {
	act, err := actServiceSetup()
	if err != nil {
		t.Fatal("Unable to setup service: ", err.Error())
	}
	filter := "DCOPS"
	d, err := act.GetDepartments(&filter)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	t.Logf("%+v", d)

}

func TestCreditRequests(t *testing.T) {
	act, err := actServiceSetup()
	if err != nil {
		t.Fatal("Unable to setup service: ", err.Error())
	}
	var pageLink *string = nil
	var qp = "requests?responsibleDepartmentId=120&status=closed&workflowState=completed"
	queryParams := &qp

Loop:
	for {
		resp, err := act.GetCreditRequests(queryParams, queryParams)
		if err != nil {
			t.Fatal("failed on get credit requests: ", err.Error())
		}

		pageLink = &resp.Links.Next.Href
		t.Log("page link", *pageLink)
		if *pageLink == "" {
			break Loop
		}
		t.Logf("%+v", resp.Requests)
	}

}

func TestGetIDForDC(t *testing.T) {
	act, err := actServiceSetup()
	if err != nil {
		t.Fatal("Unable to setup service: ", err.Error())
	}

	/*
			156: DCOPS CWL
		33: DCOPS - HKG
		34: DCOPS - IAD
		150: DCOPS Den
		157: DCOPS LBA
		35: DCOPS - LON3
		37: DCOPS - ORD
		32: DCOPS - DFW
		36: DCOPS - LON5
		149: DCOPS NYC
		154: DCOPS SHA
		152: DCOPS SJC
		155: DCOPS AMS
		151: DCOPS MCI
		153: DCOPS SIN
		38: DCOPS - SYD

	*/

	filter := "DCOPS"
	deptMap, err := act.GetDepartments(&filter)

	dc := "ORD"
	cms, err := getIDForDC(deptMap, dc)
	if err != nil {
		t.Errorf("Want: 37 (ORD)  Got: %d (%s) ", cms, dc)
	}
	if cms != 37 {
		t.Errorf("Want: 37 (ORD)  Got: %d (%s) ", cms, dc)
	}
	dc = "LON5"
	cms, err = getIDForDC(deptMap, dc)
	if err != nil {
		t.Errorf("Want: 36 (ORD)  Got: %d (%s) ", cms, dc)
	}
	if cms != 36 {
		t.Errorf("Want: 36 (ORD)  Got: %d (%s) ", cms, dc)
	}
	dc = "OR"
	cms, err = getIDForDC(deptMap, dc)
	if err == nil {
		t.Errorf("Want: DC Not Found: OR error  Got: %d (%s) ", cms, dc)
	}

	dc = "LON"
	cms, err = getIDForDC(deptMap, dc)
	if err == nil {
		t.Errorf("Want: DC Not Found: OR error  Got: %d (%s) ", cms, dc)
	}
}
