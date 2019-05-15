package coreods

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)
func setupODSTests() (ODSService, error) {
	odsUsername := os.Getenv("ODSID")
	odsPassword := os.Getenv("ODSPW")
	odsHost := os.Getenv("ODSHOST")
	odsPortS := os.Getenv("ODSPORT")
	os.Getenv("CoreProdPassword")
	//comment
	
	odsPort, err := strconv.Atoi(odsPortS)
	if err != nil {
		return nil, err
	}

	fmt.Println("ODSPW: ", odsPassword)
	
	if odsUsername == "" || odsPassword == "" || odsPortS == "" {
		return nil, fmt.Errorf("database env variables not configured")
	}
	dbConfig := ODSConfig{
		Host:             odsHost,
		Port:             odsPort,
		User:             odsUsername,
		Password:         odsPassword,
		DBName:           "Operational_reporting_CORE",

	}

	ods, err := NewODSService(dbConfig)
	return ods, err

}



func setupDmartTests() (ODSService, error) {
	odsUsername := os.Getenv("ODSID")
	odsPassword := os.Getenv("DMPW")
	odsHost := os.Getenv("DMHOST")
	odsPortS := os.Getenv("ODSPORT")
	os.Getenv("CoreProdPassword")
	//comment
	odsPort, err := strconv.Atoi(odsPortS)
	if err != nil {
		return nil, err
	}

	if odsUsername == "" || odsPassword == "" || odsPortS == "" {
		return nil, fmt.Errorf("database env variables not configured")
	}
	dbConfig := ODSConfig{
		Host:             odsHost,
		Port:             odsPort,
		User:             odsUsername,
		Password:         odsPassword,
		DBName:           "Corporate_DMART",

	}

	ods, err := NewODSService(dbConfig)
	return ods, err

}

func TestDBPing(t *testing.T) {

	dbs, err := setupDmartTests()
	if err != nil {
		t.Fatal(err)

	}

	dbs.Ping()

}

func TestOdsDB_DeviceDetails(t *testing.T) {
	dbs, err := setupODSTests()
	if err != nil {
		
		t.Fatal(err)

	}

	type TestArgs struct {
		devInputList string
	}

	testCases := []struct {
		T    TestArgs
		Want int
	}{
		{T: TestArgs{devInputList: "700656,939716,939659"}, Want: 3},
	}

	for _, tc := range testCases {
		got, err := dbs.DeviceDetails(tc.T.devInputList)
		t.Log(got)
		if err != nil || len(*got) != tc.Want {
			t.Errorf("%s:  Got: %v; Want:  %v", tc.T.devInputList, got, tc.Want)
			continue
		}
		t.Logf("%s:  Got: %v; Want:  %v", tc.T.devInputList, got, tc.Want)

		// if got != tc.Want {
		// 	t.Errorf("err:%s:  Got: %v; Want:  %v", tc.T.convertFromCurrency, got, tc.Want)
		// }
	}

}

func TestOdsDB_ExchangeRate(t *testing.T) {

	dbs, err := setupDmartTests()
	if err != nil {
		t.Fatal(err)

	}

	type TestArgs struct {
		convertFromCurrency string
		month               int
		year                int
	}

	testCases := []struct {
		T    TestArgs
		Want float64
	}{

		{T: TestArgs{"AUD", 1, 2019}, Want: 1.39971757985588},
		{T: TestArgs{"EUR", 1, 2019}, Want: 0.879225},
		{T: TestArgs{"USD", 1, 2019}, Want: 1},
		{T: TestArgs{"HKD", 1, 2019}, Want: 7.84073551516471},
		{T: TestArgs{"GBP", 8, 2017}, Want: 0.771124120667742},
	}

	for _, tc := range testCases {
		got, err := dbs.ExchangeRate(tc.T.convertFromCurrency, tc.T.month, tc.T.year)
		if err != nil || got != tc.Want {
			t.Errorf("%s:  Got: %v; Want:  %v", tc.T.convertFromCurrency, got, tc.Want)
			continue
		}
		t.Logf("%s:  Got: %v; Want:  %v", tc.T.convertFromCurrency, got, tc.Want)

		// if got != tc.Want {
		// 	t.Errorf("err:%s:  Got: %v; Want:  %v", tc.T.convertFromCurrency, got, tc.Want)
		// }
	}

}

func TestOdsDB_DeviceCount(t *testing.T) {

	dbs, err := setupDmartTests()
	if err != nil {
		t.Fatal(err)

	}

	type TestArgs struct {
		Start time.Time
		End   time.Time
		Dc    string
	}

	startTime := time.Date(2018, time.Month(12), 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2018, time.Month(12), 31, 23, 59, 59, 0, time.UTC)

	testCases := []struct {
		T    TestArgs
		Want DCServerCount
	}{

		{T: TestArgs{startTime, endTime, "ORD1"}, Want: DCServerCount{"ORD1", 23830}},
		{T: TestArgs{startTime, endTime, "IAD3"}, Want: DCServerCount{"IAD3", 20151}},
		{T: TestArgs{startTime, endTime, "DFW3"}, Want: DCServerCount{"DFW3", 14218}},
		{T: TestArgs{startTime, endTime, "LON3"}, Want: DCServerCount{"LON3", 14158}},
		{T: TestArgs{startTime, endTime, "LON5"}, Want: DCServerCount{"LON5", 5373}},
		{T: TestArgs{startTime, endTime, "DFW2"}, Want: DCServerCount{"DFW2", 2811}},
		{T: TestArgs{startTime, endTime, "HKG1"}, Want: DCServerCount{"HKG1", 1972}},
		{T: TestArgs{startTime, endTime, "SYD2"}, Want: DCServerCount{"SYD2", 1573}},
		{T: TestArgs{startTime, endTime, "IAD2"}, Want: DCServerCount{"IAD2", 1345}},
		{T: TestArgs{startTime, endTime, "SAT6"}, Want: DCServerCount{"SAT6", 586}},
		{T: TestArgs{startTime, endTime, "CDC1"}, Want: DCServerCount{"CDC1", 329}},
		{T: TestArgs{startTime, endTime, "FRA1"}, Want: DCServerCount{"FRA1", 257}},
		{T: TestArgs{startTime, endTime, "IAD60"}, Want: DCServerCount{"IAD60", 239}},
		{T: TestArgs{startTime, endTime, "YYZ2"}, Want: DCServerCount{"YYZ2", 235}},
		{T: TestArgs{startTime, endTime, "YYZ1"}, Want: DCServerCount{"YYZ1", 103}},
		{T: TestArgs{startTime, endTime, "ATL60"}, Want: DCServerCount{"ATL60", 48}},
		{T: TestArgs{startTime, endTime, "ATL61"}, Want: DCServerCount{"ATL61", 34}},
		{T: TestArgs{startTime, endTime, "DEN60"}, Want: DCServerCount{"DEN60", 34}},
		{T: TestArgs{startTime, endTime, "SYD4"}, Want: DCServerCount{"SYD4", 31}},
		{T: TestArgs{startTime, endTime, "LON4"}, Want: DCServerCount{"LON4", 30}},
		{T: TestArgs{startTime, endTime, "DFW1"}, Want: DCServerCount{"DFW1", 24}},
		{T: TestArgs{startTime, endTime, "AUS2"}, Want: DCServerCount{"AUS2", 23}},
		{T: TestArgs{startTime, endTime, "SIN80"}, Want: DCServerCount{"SIN80", 19}},
		{T: TestArgs{startTime, endTime, "AUS3"}, Want: DCServerCount{"AUS3", 12}},
		{T: TestArgs{startTime, endTime, "TST1"}, Want: DCServerCount{"TST1", 6}},
		{T: TestArgs{startTime, endTime, "BCB2"}, Want: DCServerCount{"BCB2", 5}},
		{T: TestArgs{startTime, endTime, "DC-NA"}, Want: DCServerCount{"DC-NA", 4}},
		{T: TestArgs{startTime, endTime, "SMDC"}, Want: DCServerCount{"SMDC", 3}},
		{T: TestArgs{startTime, endTime, "AMS1"}, Want: DCServerCount{"AMS1", 3}},
		{T: TestArgs{startTime, endTime, "SYD1"}, Want: DCServerCount{"SYD1", 2}},
		{T: TestArgs{startTime, endTime, "FRA30"}, Want: DCServerCount{"FRA30", 2}},
		{T: TestArgs{startTime, endTime, "ZRH1"}, Want: DCServerCount{"ZRH1", 2}},
		{T: TestArgs{startTime, endTime, "STO2"}, Want: DCServerCount{"STO2", 1}},
		{T: TestArgs{startTime, endTime, "HKO1"}, Want: DCServerCount{"HKO1", 1}},
		{T: TestArgs{startTime, endTime, "STO1"}, Want: DCServerCount{"STO1", 1}},
		{T: TestArgs{startTime, endTime, "INVALID"}, Want: DCServerCount{}},
	}

	for _, tc := range testCases {
		got, _ := dbs.DataCenterServerCount(tc.T.Start, tc.T.End, tc.T.Dc)

		t.Logf("Got: %+v; Want:  %+v", got, tc.Want)

		if got != tc.Want {
			t.Errorf("Got: %v; Want:  %v", got, tc.Want)
		}
	}

}
