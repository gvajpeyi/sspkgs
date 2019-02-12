package hitch

//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"testing"
//)
//
//func TestIncidentsByDates_StartDateOnly(t *testing.T) {
//
//	_, err := GetIncidentsByDates("2018-09-01T00:00:00Z", "")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//}
//
//func TestIncidentsByDates_EndDateOnly(t *testing.T) {
//
//	_, err := GetIncidentsByDates("", "2018-08-01T00:00:00Z")
//	if err == nil {
//		t.Fatal(err.Error())
//	}
//
//}
//
//func TestIncidentsByDates_StartDateAndEndDate(t *testing.T) {
//
//	_, err := GetIncidentsByDates("2018-08-01T00:00:00Z", "2018-08-01T00:00:00Z")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//}
//
//func TestGetIncidentsByDates_OutOfOrderDates(t *testing.T) {
//
//	_, err := GetIncidentsByDates("2018-08-01T00:00:00Z", "2018-08-01T00:00:00Z")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//}
//
//func TestRemoveDuplicateTickets(t *testing.T) {
//
//	raw, err := ioutil.ReadFile("./exampleResponse.json")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	var testData HitchResponse
//	json.Unmarshal(raw, &testData)
//
//	beforeRemoval := len(testData.Results)
//	testData.RemoveDuplicateTickets()
//	afterRemoval := len(testData.Results)
//
//	if beforeRemoval <= afterRemoval {
//		t.Fatal("No records removed")
//	}
//
//}

//
//func TestCountPerDataCenter(t *testing.T) {
//
//	raw, err := ioutil.ReadFile("./exampleResponse.json")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	var testData gohitch.HitchResponse
//	json.Unmarshal(raw, &testData)
//
//	//testData, err := gohitch.GetIncidentsByDates("2018-06-01T00:00:00Z", "2018-06-30T11:59:59Z")
//	//if err != nil{
//	//	t.Fatal(err.Error())
//	//}
//
//	testData.RemoveDuplicateTickets()
//	failures := []string{}
//	dfw1, dfw2, dfw3, hkg1, iad2, iad3, lon1, lon3, lon5, ord1, syd2, fra1 := testData.CountPerDataCenter()
//
//	if dfw1 != 1 {
//		failures = append(failures, "dfw1")
//	}
//	if dfw2 != 1 {
//		failures = append(failures, "dfw2")
//	}
//	if dfw3 != 1 {
//		failures = append(failures, "dfw3")
//	}
//	if hkg1 != 1 {
//		failures = append(failures, "hkg1")
//	}
//	if iad2 != 1 {
//		failures = append(failures, "iad2")
//	}
//	if iad3 != 1 {
//		failures = append(failures, "iad3")
//	}
//	if lon1 != 1 {
//		failures = append(failures, "lon1")
//	}
//	if lon3 != 1 {
//		failures = append(failures, "lon3")
//	}
//	if lon5 != 1 {
//		failures = append(failures, "lon5")
//	}
//	if ord1 != 1 {
//		failures = append(failures, "ord1")
//	}
//	if syd2 != 1 {
//		failures = append(failures, "syd2")
//	}
//	if fra1 != 1 {
//		failures = append(failures, "fra1")
//	}
//
//	fmt.Printf("\ndfw1:%d dfw2:%d dfw3:%d hkg1:%d iad2:%d iad3:%d lon1:%d lon3:%d lon5:%d ord1:%d syd2:%d fra1:%d\n", dfw1, dfw2, dfw3, hkg1, iad2, iad3, lon1, lon3, lon5, ord1, syd2, fra1)
//
//	if len(failures) > 0 {
//		t.Logf("\nFollowing have failed: %s\n", failures)
//		t.Fail()
//	}
//
//}

//func TestWriteCountsToDatabase(t *testing.T) {
//
//	testData, err := GetIncidentsByDates("2018-08-12T00:00:00Z", "2018-08-20T11:59:59Z")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//	fmt.Printf("incidents\n%+v\n\n", testData.Results)
//	testData.RemoveDuplicateTickets()
//
//	err = testData.WriteCountsToDatabase()
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//
//}
