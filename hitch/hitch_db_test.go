package hitch
//import (
//	"encoding/json"
//	"fmt"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/metricsdb"
//
//	"io/ioutil"
//	"testing"
//)
//
//var hitchReports = []*HitchReport{}
//var metricDBConfig = new(metricsdb.MetricsDatabase)
//var hitchDB = new(HitchDB)
//var initError error
//
//func TestGetAllIncidents(t *testing.T) {
//
//	raw, initError := ioutil.ReadFile("../config/db.json")
//	if initError != nil {
//		fmt.Printf("hitchReports failure: %s\n", initError.Error())
//	}
//	initError = json.Unmarshal(raw, &metricDBConfig)
//	if initError != nil {
//		fmt.Printf("hitchReports failure: %s\n", initError.Error())
//
//	}
//	hitchDB = NewHitchDatabase(metricDBConfig)
//
//	fmt.Printf("Hitch db: %+v\n", *hitchDB)
//
//	if initError != nil {
//		t.Errorf("Failed: %s", initError.Error())
//	}
//
//	reports, _, err := hitchDB.GetHitchIncidents("DFW3", nil)
//	if err != nil {
//		t.Errorf("hitchReports failure: %s\n", err.Error())
//	}
//
//	fmt.Printf("reports:  %+v\n", reports)
//
//}
