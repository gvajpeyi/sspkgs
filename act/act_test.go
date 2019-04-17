package uptime

import (
	"flag"
	"fmt"
	"github.rackspace.com/SegmentSupport/metrics-dashboard/config"
	"github.rackspace.com/SegmentSupport/sspkgs/coreods"
	"github.rackspace.com/SegmentSupport/sspkgs/identity"
	"log"
	"os"
	"testing"
)

var boolPtr *bool
var cfg *config.Config

func actServiceSetup() (*Services, error) {

	if boolPtr == nil || cfg == nil {
		boolPtr = flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
		flag.Parse()
		tempCfg := config.LoadConfig(*boolPtr)
		cfg = &tempCfg
	}

	odsCfg := config.LoadDBConfig(os.Getenv("ODSID"), os.Getenv("ODSPW"), os.Getenv("ODSHOST"), os.Getenv("ODSPORT"),  os.Getenv("ODSDB"))

	fmt.Println("ODSDB Host:  ", odsCfg.Host)
	fmt.Println("ODSDB User:  ", odsCfg.User)
	fmt.Println("ODSDB Pass:  ", odsCfg.Password)
	fmt.Println("ODSDB DBName:  ", odsCfg.DBName)
	fmt.Println("ODSDB Port:  ", odsCfg.Port)

	ods, err := coreods.NewODSService(odsCfg)
	if err != nil {
		log.Fatalf("unable to setup service: %v", err)
	}

	cfg.ACTUser = os.Getenv("SSID")
	cfg.ACTPass = os.Getenv("SSPW")

	id := identity.NewIdentityService(cfg.HttpClient, cfg.IdentityURL, cfg.Logger)

	services, err := NewServices(

		WithACT(cfg.ACTUrl, cfg.ACTUser, cfg.ACTPass, cfg.ACTDomain, ods, id, cfg.Logger),
	)

	if err != nil {
		return nil, err
	}

	return services, err

}

func TestResponsibleDepartments(t *testing.T) {
	s, err := actServiceSetup()
	if err != nil {
		t.Fatal("Unable to setup service: ", err.Error())
	}
	filter := "DCOPS"
	_, err = s.ACT.GetDepartments(&filter)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

}

func TestCreditRequests(t *testing.T) {
	s, err := actServiceSetup()
	if err != nil {
		t.Fatal("Unable to setup service: ", err.Error())
	}
	var pageLink *string = nil
	var qp = "?responsibleDepartmentId=35&status=closed&workflowState=completed"
	queryParams := &qp

Loop:
	for {
		resp, err := s.ACT.GetCreditRequests(pageLink, queryParams)
		if err != nil {
			t.Fatal("failed on get credit requests: ", err.Error())
		}

		pageLink = &resp.Links.Next.Href

		if *pageLink == "" {
			break Loop
		}

	}

}

func TestCurrentMonthSums(t *testing.T) {
	s, err := actServiceSetup()
	if err != nil {
		t.Fatal("Unable to setup service: ", err.Error())
	}
	cms, err := s.ACT.OverviewSums("DFW")

	if err != nil {
		t.Errorf("Want: something   Got: something else: %s", err.Error())
	}
	t.Logf("cms:\n%+v\n", cms)

}

func TestGetIDForDC(t *testing.T) {
	s, err := actServiceSetup()
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
	deptMap, err := s.ACT.GetDepartments(&filter)

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
