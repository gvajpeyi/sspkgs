package cyclecounts_test
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/inventory"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/metricsdb"
//
//	"io/ioutil"
//	"testing"
//)
//
//func TestGetCycleCountsForDC(t *testing.T) {
//
//
//	dc := "IAD"
//
//	//dcs := []string{
//	//	"DCOPS - DFW",
//	//	"DCOPS - HKG",
//	//	"DCOPS - IAD",
//	//	"DCOPS - LON",
//	//	"DCOPS - LON",
//	//	"DCOPS - ORD",
//	//	"DCOPS - SYD",
//	//}
//
//	raw, err := ioutil.ReadFile("../config/db.json")
//	if err != nil {
//		fmt.Printf("cycle count failure to get appconfig: %s\n", err.Error())
//
//	}
//	var metricDBConfig  = new (metricsdb.MetricsDatabase)
//
//
//	err = json.Unmarshal(raw, &metricDBConfig)
//	if err != nil {
//		fmt.Printf("cycle count failure to load db settings: %s\n", err.Error())
//
//	}
//
//
//
//	var cycleCountDB = cyclecounts.NewCycleCountDatabase(metricDBConfig)
//
//
//		counts, err := cycleCountDB.GetCycleCounts(dc)
//		if err != nil {
//			t.Log(err.Error())
//			t.Fail()
//		}
//		fmt.Printf("%+v\n", *counts)
//
//
//}
//
