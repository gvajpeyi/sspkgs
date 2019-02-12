package creditmemo
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.rackspace.com/SegmentSupport/raxss/pkg/metricsdb"
// 	"io/ioutil"
// 	"strconv"
// 	"testing"
//
// )
//
// func TestGetResponsibleDepartments(t *testing.T) {
//
// 	var searchString string = "DCOPS"
//
// 	responsibleDepartments, err := GetResponsibleDepartments("AAC0MahOP-xzQ8tbJ4lX3DZs9sTjMYnyfj8S3L7On5pzl09Z-lqBJWiBQKz7q9K6gUudwyteZJ0ObEeXrd1CkCuINahZZpAhzATsZt8xoaWH8DKG_iKE4otz", &searchString)
// 	if err != nil {
// 		t.Log(err.Error())
// 		t.Fail()
// 	}
//
// 	fmt.Println("count: ", len(responsibleDepartments))
// 	for k, v := range responsibleDepartments {
// 		fmt.Println(k, v)
// 	}
//
// }
//
// func TestGetCreditMemoRequests(t *testing.T) {
//
// 	dcMap := make(map[int]string)
// 	dcMap[33] = "DCOPS - HKG"
// 	dcMap[34] = "DCOPS - IAD"
// 	dcMap[35] = "DCOPS - LON3"
// 	dcMap[36] = "DCOPS - LON5"
// 	dcMap[37] = "DCOPS - ORD"
// 	dcMap[38] = "DCOPS - SYD"
// 	dcMap[32] = "DCOPS - DFW"
//
// 	/*
// 		33 DCOPS - HKG
// 		34 DCOPS - IAD
// 		35 DCOPS - LON3
// 		36 DCOPS - LON5
// 		37 DCOPS - ORD
// 		38 DCOPS - SYD
// 		32 DCOPS - DFW
//
// 	*/
//
//
//
//
// 	month := 3
// 	year := 2018
// 	for k, _ := range dcMap {
// 		searchString := make(map[string]string)
// 		searchString["responsibleDepartmentId"] = strconv.Itoa(k)
// 		searchString["status"] = "closed"
// 		searchString["workflowState"] = "completed"
// 		cmr, err := GetCreditMemoRequests("AAC0MahOP-xzQ8tbJ4lX3DZs9sTjMYnyfj8S3L7On5pzl09Z-lqBJWiBQKz7q9K6gUudwyteZJ0ObEeXrd1CkCuINahZZpAhzATsZt8xoaWH8DKG_iKE4otz", searchString, month, year)
// 		if err != nil {
// 			t.Log(err.Error())
// 			t.Fail()
// 		}
// 		InsertMetricIntoDatabase(cmr)
//
// 		//for _, cm := range cmr.Requests {
// 		//	//layout := "2006-01-02 15:04:05 -0700 MST"
// 		//	//closedTime, err := time.Parse(layout, cm.Request.UpdatedDatetime.String())
// 		//	if err != nil {
// 		//		continue
// 		//	}
// 		//	validResponses := []string{}
// 		//	validResponses = append(validResponses, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%f\t%s\t%s\t%s", cm.Request.ResponsibleDepartment.Name, cm.Request.UpdatedDatetime, cm.Request.ClosedDate, cm.Request.ContractingEntity.Code, cm.Request.LastApprovedBy.DisplayName, cm.Request.Currency.Code, cm.Request.WorkflowState, cm.Request.Amount, cm.Request.Status, cm.Request.ReasonCode.Name, cm.Request.ServiceFailure.Name))
// 		//	respDeptTotal += cm.Request.Amount + cm.Request.AdditionalAmount
// 		//	if len(validResponses) > 0 {
// 		//		responses = append(responses, validResponses)
// 		//	}
// 		//}
//
// 	}
// }
//
// func TestGetCreditMemosForRegion(t *testing.T) {
// 	dcs := []string{
// 		"DCOPS - DFW",
// 		"DCOPS - HKG",
// 		"DCOPS - IAD",
// 		"DCOPS - LON3",
// 		"DCOPS - LON5",
// 		"DCOPS - ORD",
// 		"DCOPS - SYD",
// 	}
//
// 	raw, err := ioutil.ReadFile("../config/db.json")
// 	if err != nil {
// 		fmt.Printf("credit memo failure to get appconfig: %s\n", err.Error())
//
// 	}
// 	var metricDBConfig  = new (metricsdb.MetricsDatabase)
//
//
// 	err = json.Unmarshal(raw, &metricDBConfig)
// 	if err != nil {
// 		fmt.Printf("credit memo failure to load db settings: %s\n", err.Error())
//
// 	}
//
//
//
// 	var creditMemoDB = NewCreditMemoDatabase(metricDBConfig)
//
//
// 	for _, dc := range dcs {
// 		counts, err := creditMemoDB.GetCurrentMonthCounts(dc)
// 		if err != nil {
// 			t.Log(err.Error())
// 			t.Fail()
// 		}
// 		fmt.Printf("%+v\n", *counts)
//
// 	}
// }
