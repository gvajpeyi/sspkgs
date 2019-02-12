package helpanto
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"strconv"
// 	"testing"
//
// 	"github.rackspace.com/SegmentSupport/raxss/pkg/coreods"
// )
//
//
// type Currency string
//
// func setupTests() (coreods.ODSService, error) {
// 	odsUsername := os.Getenv("ODSID")
// 	odsPassword := os.Getenv("ODSPW")
// 	// odsHost := os.Getenv("ODSHOST")
// 	odsPortS := os.Getenv("ODSPORT")
//
// 	odsPort, err := strconv.Atoi(odsPortS)
// 	if err != nil {
// 		return nil, err
// 	}
//
//
// 	if odsUsername == "" || odsPassword == "" || odsPortS == ""{
// 		return nil, fmt.Errorf("database env variables not configured")
// 	}
// 	dbConfig := coreods.ODSConfig{
// 		Host:     "EBI-DATAMART",
// 		Port:     odsPort,
// 		User:     odsUsername,
// 		Password: odsPassword,
// 		DBName:   "Corporate_DMART",
// 	}
//
//
// 	ods, err := coreods.NewODSService(dbConfig)
//
// 	return ods, err
//
// }
//
//
//
// func TestConvertToUSD(t *testing.T) {
//
// 	dbs, err := setupTests()
// 	if err != nil {
// 		t.Fatal(err)
//
// 	}
// 	/*
//
// 			Total Pre Tax Amount (USD)	AUD	Closed Date
// 		2,816.67	3,630.00	3/14/18	2816.672242	max val
// 		1,324.52	1,805.19	8/27/18	1324.521518	max val
// 		456.85   	594.61	    4/17/18	456.8536485
//
//
// 	*/
//
// 	usd, err := ConvertToUSD(connstring, "AUD", 4, 2018, 594.61)
// 	if err != nil {
// 		t.Log("Failed: ", err.Error())
//
// 	}
// 	fmt.Println(" ")
// 	fmt.Println("Final:  ", *usd)
// }
// //
// // type EbiDatabase struct {
// // 	Name             string `json:"name"`
// // 	Port             int    `json:"port"`
// // 	IPAddress        string `json:"ip_address"`
// // 	DatabaseName     string `json:"database_name"`
// // 	DatabaseUser     string `json:"database_user"`
// // 	DatabasePassword string `json:"database_password"`
// // }
// //
// // func (s *EbiDatabase) GenerateConnectionString() string {
// // 	return fmt.Sprintf("odbc:server=%s; port=%d; user id=%s;password=%s; database=%s;log=3;encrypt=false;TrustServerCertificate=true", s.IPAddress, s.Port, s.DatabaseUser, s.DatabasePassword, s.DatabaseName)
// // }
//
// func GetEbiConfigFile(filename string) *EbiDatabase {
// 	//	log.Info("In Get Consoles")
// 	if filename == "" {
// 		filename = "ebiConfig.json"
// 	}
// 	raw, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		os.Exit(1)
// 	}
//
// 	var ebi *EbiDatabase
// 	err = json.Unmarshal(raw, &ebi)
// 	if err != nil {
// 		os.Exit(1)
// 	}
//
// 	return ebi
// }
//
//
//
