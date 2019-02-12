package hitch
//
//import (
//	"bytes"
//	"crypto/tls"
//	"encoding/json"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/jmoiron/sqlx"
//	"github.rackspace.com/SegmentSupport/raxss/internal/helpanto"
//	"io/ioutil"
//	"net/http"
//	"sort"
//	"strings"
//)
//
//func endDateRequest(startDateUnix int64, endDateUnix int64) (HitchResponse, error) {
//
//	var queryWithEndDate QueryWithEndDate
//	queryWithEndDate.CreatedAt.Gte.Date = startDateUnix
//	queryWithEndDate.CreatedAt.Lte.Date = endDateUnix
//	queryWithEndDate.RackspaceCaused = true
//	queryWithEndDate.IncidentCause = "Racker Error"
//	queryWithEndDate.Segment = "datacenter"
//	queryWithEndDate.Departments.In = []string{"DCOPS"}
//	queryWithEndDate.Status = "closed"
//	payloadBytes, err := json.Marshal(queryWithEndDate)
//	if err != nil {
//		return HitchResponse{}, err
//	}
//	fmt.Printf("\n\n%v\n\n", string(payloadBytes))
//	hitchResponse, err := makeHitchPostRequest(payloadBytes)
//	if err != nil {
//		return HitchResponse{}, err
//	}
//	return hitchResponse, nil
//
//}
//
//func startDateOnlyRequest(startDateUnix int64) (HitchResponse, error) {
//	var queryWithOnlyStartDate QueryWithOnlyStartDate
//	queryWithOnlyStartDate.CreatedAt.Gte.Date = startDateUnix
//	queryWithOnlyStartDate.RackspaceCaused = true
//	queryWithOnlyStartDate.IncidentCause = "Racker Error"
//	queryWithOnlyStartDate.Segment = "datacenter"
//	queryWithOnlyStartDate.Status = "closed"
//	queryWithOnlyStartDate.Departments.In = []string{"DCOPS"}
//
//	payloadBytes, err := json.Marshal(queryWithOnlyStartDate)
//	if err != nil {
//		return HitchResponse{}, err
//	}
//
//	hitchResponse, err := makeHitchPostRequest(payloadBytes)
//	if err != nil {
//		return HitchResponse{}, err
//	}
//	return hitchResponse, nil
//
//}
//
//func makeHitchPostRequest(payloadBytes []byte) (HitchResponse, error) {
//
//	url := ""
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := &http.Client{
//
//		Transport: tr,
//	}
//
//	url = "https://hitch.res.rackspace.com/api/incidents/query"
//
//	var resp *http.Response
//	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
//
//	if err != nil {
//		return HitchResponse{}, err
//
//	}
//
//	request.Header.Add("Content-Type", "application/json")
//	resp, err = client.Do(request)
//
//	if err != nil {
//		return HitchResponse{}, err
//
//	}
//
//	if resp.StatusCode > 299 || resp.StatusCode < 200 {
//		return HitchResponse{}, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
//	}
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		return HitchResponse{}, err
//	}
//
//	var hitchResponse HitchResponse
//	err = json.Unmarshal(body, &hitchResponse)
//	if err != nil {
//		return HitchResponse{}, err
//
//	}
//
//	return hitchResponse, nil
//
//}
//
//func GetIncidentsByDates(startDate string, endDate string) (HitchResponse, error) {
//
//	var hitchResponse HitchResponse
//	var endDateUnix int64
//
//	startDateUnix, err := helpanto.DateStringToUnixMicro(startDate)
//	if err != nil {
//		return HitchResponse{}, err
//	}
//
//	if endDate != "" {
//		endDateUnix, err = helpanto.DateStringToUnixMicro(endDate)
//		if err != nil {
//			return HitchResponse{}, err
//		}
//		hitchResponse, err = endDateRequest(startDateUnix, endDateUnix)
//		if err != nil {
//			return HitchResponse{}, err
//		}
//	} else {
//
//		hitchResponse, err = startDateOnlyRequest(startDateUnix)
//
//	}
//
//	return hitchResponse, nil
//
//}
//
//func (h *HitchResponse) RemoveDuplicateTickets() {
//
//	sort.Sort(ByTicketNumberLastEmailDate(h.Results))
//	idsToRemove := []string{}
//	lastTicketNumber := ""
//	for _, incident := range h.Results {
//		if lastTicketNumber == incident.TicketNumber {
//			idsToRemove = append(idsToRemove, incident.ID.Oid)
//			lastTicketNumber = incident.TicketNumber
//
//		} else {
//			lastTicketNumber = incident.TicketNumber
//		}
//
//	}
//
//	for i := len(h.Results) - 1; i >= 0; i-- {
//
//		for j := len(idsToRemove) - 1; j >= 0; j-- {
//			if h.Results[i].ID.Oid == idsToRemove[j] {
//				h.Results = append(h.Results[:i], h.Results[i+1:]...)
//				break
//
//			}
//
//		}
//
//	}
//
//}
//
//func (h *HitchResponse) WriteCountsToDatabase() error {
//	dbDSN := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v", "dc_metrics_hitch", "kpCizmbVDdsSM4wXZRUK", "127.0.0.1", 3306, "metrics")
//
//	db, err := sqlx.Connect("mysql", dbDSN)
//	if err != nil {
//		return err
//	}
//
//	tx, err := db.Begin()
//	if err != nil {
//		return err
//	}
//
//	sql := `Insert into dc_metrics_hitch (
//    ticket_number,
//	temperature,
//	outage_status,
//	start_date,
//	short_description,
//	segment,
//	rackspace_caused,
//	outage_type,
//	outage_cause,
//	outage_impact,
//	last_email_sent,
//	escalated,
//	end_date,
//	departments,
//	data_center,
//	customer_caused,
//	created_by_name,
//	created_by_email,
//	time,
//	account_number,
//	account_name,
//	account_segment,
//	support_team)
//	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
//
//	stmt, err := tx.Prepare(sql)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	defer stmt.Close()
//	for _, incident := range h.Results {
//		startTime := helpanto.UnixMicroToDateString(incident.StartDate.Date)
//		startTime, _ = helpanto.ConvertDateForMariaDBInsert(startTime)
//
//		endTime := helpanto.UnixMicroToDateString(incident.EndDate.Date)
//		endTime, _ = helpanto.ConvertDateForMariaDBInsert(endTime)
//
//		createdTime := helpanto.UnixMicroToDateString(incident.CreatedAt.Date)
//		createdTime, _ = helpanto.ConvertDateForMariaDBInsert(createdTime)
//
//		lastEmailSent := helpanto.UnixMicroToDateString(incident.LastEmail.Date)
//		lastEmailSent, _ = helpanto.ConvertDateForMariaDBInsert(lastEmailSent)
//
//		_, err = stmt.Exec(incident.TicketNumber, incident.Temperature, incident.Status, startTime,
//			incident.ShortDescription, incident.Segment, incident.RackspaceCaused, incident.IncidentType, incident.IncidentCause,
//			incident.IncidentImpact, lastEmailSent, incident.Escalated, endTime, fmt.Sprintf("%s", incident.Departments),
//			incident.DataCenter, incident.CustomerCaused, incident.CreatedByName, incident.CreatedByEmail, createdTime, incident.AccountNumber, incident.AccountName, incident.AccountSegment, incident.SupportTeam)
//		if err != nil {
//			if strings.Contains(err.Error(), "Error 1062: Duplicate entry ") {
//				continue
//			}
//			tx.Rollback()
//
//			return err
//		}
//
//	}
//
//	err = tx.Commit()
//	if err != nil {
//
//		tx.Rollback()
//		return err
//	}
//
//	return nil
//}
//
///*func (h *HitchResponse) CountPerDataCenter() (int, int, int, int, int, int, int, int, int, int, int, int) {
//
//	dfw1 := 0
//	dfw2 := 0
//	dfw3 := 0
//	hkg1 := 0
//	iad2 := 0
//	iad3 := 0
//	lon1 := 0
//	lon3 := 0
//	lon5 := 0
//	ord1 := 0
//	syd2 := 0
//	fra1 := 0
//
//	for _, incident := range h.Results {
//
//		switch strings.ToUpper(incident.DataCenter) {
//		case "DFW1":
//			dfw1++
//		case "DFW2":
//			dfw2++
//		case "DFW3":
//			dfw3++
//		case "HKG1":
//			hkg1++
//		case "IAD2":
//			iad2++
//		case "IAD3":
//			iad3++
//		case "LON1":
//			lon1++
//		case "LON3":
//			lon3++
//		case "LON5":
//			lon5++
//		case "ORD1":
//			ord1++
//		case "SYD2":
//			syd2++
//		case "FRA1":
//			fra1++
//
//		}
//
//	}
//
//	return dfw1, dfw2, dfw3, hkg1, iad2, iad3, lon1, lon3, lon5, ord1, syd2, fra1
//
//}
//*/
////func (h *HitchResponse) GenerateCsv(filename string) string {
////
////	fileOut := ""
////	rowOut := ""
////
////	fileHandle, _ := os.Create(filename)
////	writer := bufio.NewWriter(fileHandle)
////	defer fileHandle.Close()
////	headers := "DCOPS At Fault;ID;AcctNumber;AcctName; AcctLeadTech; AcctLeadTechEmail;AcctLeadTechName;AcctManager;AcctMgrEmail;AcctMgrName;AccntSegment;AcctSvcLevel;AmSupervisorEmail;CreatedTime;CreatedBy;CreatedByEmail;CreatedByName;CustomerCaused;DataCenter;Departments;EndDate;Escalated;LastEmailDate;Log;RackspaceCaused;Segment;StartDate;Status;Subscriptions;SupportTeam;Temperature;TicketNumber;TicketSource\n"
////	fmt.Fprintln(writer, headers)
////
////	for i, rowIn := range h.Results {
////
////		fmt.Println("Row Cnt: ", i)
////		rowOut = fmt.Sprintf("%t", atFault(rowIn.Departments, "DCOPS", rowIn.Segment, rowIn.OutageCause,  rowIn.RackspaceCaused))
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.ID.Oid)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountNumber)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountName)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountLeadTech)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountLeadTechEmail)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountLeadTechName)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountManager)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountManagerEmail)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountManagerName)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountSegment)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AccountServiceLevel)
////		//	rowOut  = fmt.Sprintf("\"%s;%s\"", rowOut, rowIn.AdditionalInfo)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.AmSupervisorEmail)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, time.Unix(rowIn.CreatedAt.Date, 0).String())
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.CreatedBy)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.CreatedByEmail)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.CreatedByName)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%t", rowIn.CustomerCaused))
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.DataCenter)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%s", rowIn.Departments))
////		//rowOut  = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%d", rowIn.Devices))
////		rowOut = fmt.Sprintf("%s;%s", rowOut, time.Unix(rowIn.EndDate.Date, 0).String())
////		rowOut = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%t", rowIn.Escalated))
////		rowOut = fmt.Sprintf("%s;%s", rowOut, time.Unix(rowIn.LastEmail.Date, 0).String())
////		rowOut = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%s", rowIn.Log))
////		//	rowOut  = fmt.Sprintf("%s;%s", rowOut, rowIn.OutageCause)
////		//rowOut  = fmt.Sprintf("%s;%s", rowOut, rowIn.OutageImpact)
////		//	rowOut  = fmt.Sprintf("%s;%s", rowOut, rowIn.OutageStatus)
////		//		rowOut  = fmt.Sprintf("%s;%s", rowOut, rowIn.OutageType)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%t", rowIn.RackspaceCaused))
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.Segment)
////		//rowOut  = fmt.Sprintf("%s;%s", rowOut, rowIn.ShortDescription)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, time.Unix(rowIn.StartDate.Date, 0).String())
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.Status)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, fmt.Sprintf("%s", rowIn.Subscriptions))
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.SupportTeam)
////
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.Temperature)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.TicketNumber)
////		rowOut = fmt.Sprintf("%s;%s", rowOut, rowIn.TicketSource)
////		fmt.Fprintln(writer, rowOut)
////		fileOut = fmt.Sprintf("%s\n%s", fileOut, rowOut)
////	}
////
////	writer.Flush()
////
////	return fileOut
////
////}
//
//type ByTicketNumberLastEmailDate []HitchIncident
//
//func (incidents ByTicketNumberLastEmailDate) Len() int {
//	return len(incidents)
//}
//
//func (incidents ByTicketNumberLastEmailDate) Swap(i, j int) {
//	incidents[i], incidents[j] = incidents[j], incidents[i]
//}
//
//func (incidents ByTicketNumberLastEmailDate) Less(i, j int) bool {
//	if strings.Compare(incidents[i].TicketNumber, incidents[j].TicketNumber) == -1 {
//		return true
//	}
//	if strings.Compare(incidents[i].TicketNumber, incidents[j].TicketNumber) == 1 {
//		return false
//	}
//
//	return incidents[i].LastEmail.Date < incidents[j].LastEmail.Date
//
//}
//
//type HitchResponse struct {
//	Count int `json:"count" csv:"count"`
//
//	Query struct {
//		CreatedAt struct {
//			Gte string `json:"$gte" csv:"$gte"`
//			Lte string `json:"$lte" csv:"$lte"`
//		} `json:"created_at" csv:"created_at"`
//	} `json:"query" csv:"query"`
//	Results []HitchIncident
//}
//
//type HitchIncident struct {
//	ID struct {
//		Oid string `json:"$oid" csv:"$oid"`
//	} `json:"_id" csv:"_id"`
//	AccountLeadTech      string `json:"account_lead_tech" csv:"account_lead_tech"`
//	AccountLeadTechEmail string `json:"account_lead_tech_email" csv:"account_lead_tech_email"`
//	AccountLeadTechName  string `json:"account_lead_tech_name" csv:"account_lead_tech_name"`
//	AccountManager       string `json:"account_manager" csv:"account_manager"`
//	AccountManagerEmail  string `json:"account_manager_email" csv:"account_manager_email"`
//	AccountManagerName   string `json:"account_manager_name" csv:"account_manager_name"`
//	AccountName          string `json:"account_name" csv:"account_name"`
//	AccountNumber        string `json:"account_number" csv:"account_number"`
//	AccountSegment       string `json:"account_segment" csv:"account_segment"`
//	AccountServiceLevel  string `json:"account_service_level" csv:"account_service_level"`
//	AdditionalInfo       string `json:"additional_info" csv:"additional_info"`
//	AmSupervisorEmail    string `json:"am_supervisor_email" csv:"am_supervisor_email"`
//	CreatedAt            struct {
//		Date int64 `json:"$date" csv:"$date"`
//	} `json:"created_at" csv:"created_at"`
//	CreatedBy      string   `json:"created_by" csv:"created_by"`
//	CreatedByEmail string   `json:"created_by_email" csv:"created_by_email"`
//	CreatedByName  string   `json:"created_by_name" csv:"created_by_name"`
//	CustomerCaused bool     `json:"customer_caused" csv:"customer_caused"`
//	DataCenter     string   `json:"data_center" csv:"data_center"`
//	Departments    []string `json:"departments" csv:"Departments"`
//	Devices        []struct {
//		DeviceID     json.Number `json:"device_id" csv:"device_id"`
//		DeviceLink   string      `json:"device_link" csv:"device_link"`
//		DeviceName   string      `json:"device_name" csv:"device_name"`
//		DeviceStatus string      `json:"device_status" csv:"device_status"`
//	} `json:"devices" csv:"devices"`
//	EndDate struct {
//		Date int64 `json:"$date" csv:"$date"`
//	} `json:"end_date" csv:"$end_date"`
//	Escalated bool `json:"escalated" csv:"escalated"`
//	LastEmail struct {
//		Date int64 `json:"$date" csv:"$date"`
//	} `json:"last_email" csv:"last_email"`
//	Log              []string `json:"log" csv:"log"`
//	IncidentCause    string   `json:"outage_cause" csv:"outage_cause"`
//	IncidentImpact   string   `json:"outage_impact" csv:"outage_impact"`
//	IncidentStatus   string   `json:"outage_status" csv:"outage_status"`
//	IncidentType     string   `json:"outage_type" csv:"outage_type"`
//	RackspaceCaused  bool     `json:"rackspace_caused" csv:"rackspace_caused"`
//	Segment          string   `json:"segment" csv:"segment"`
//	ShortDescription string   `json:"short_description" csv:"short_description"`
//	StartDate        struct {
//		Date int64 `json:"$date" csv:"$date"`
//	} `json:"start_date" csv:"start_date"`
//	Status            string      `json:"status" csv:"status"`
//	Subscriptions     []string    `json:"subscriptions" csv:"subscriptions"`
//	SupportTeam       string      `json:"support_team" csv:"support_team"`
//	TechEventMgr      interface{} `json:"tech_event_mgr" csv:"tech_event_mgr"`
//	TechEventMgrEmail interface{} `json:"tech_event_mgr_email" csv:"tech_event_mgr_email"`
//	TechEventMgrName  interface{} `json:"tech_event_mgr_name" csv:"tech_event_mgr_name"`
//	TechExpert        interface{} `json:"tech_expert" csv:"tech_expert"`
//	TechExpertEmail   interface{} `json:"tech_expert_email" csv:"tech_expert_email"`
//	TechExpertName    interface{} `json:"tech_expert_name" csv:"tech_expert_name"`
//	TechMajorMgr      interface{} `json:"tech_major_mgr" csv:"tech_major_mgr"`
//	TechMajorMgrEmail interface{} `json:"tech_major_mgr_email" csv:"tech_major_mgr_email"`
//	TechMajorMgrName  interface{} `json:"tech_major_mgr_name" csv:"tech_major_mgr_name"`
//	Temperature       string      `json:"temperature" csv:"temperature"`
//	TicketNumber      string      `json:"ticket_number" csv:"ticket_number"`
//	TicketSource      string      `json:"ticket_source" csv:"ticket_source"`
//}
//
//type QueryWithEndDate struct {
//	CreatedAt struct {
//		Gte struct {
//			Date int64 `json:"$date"`
//		} `json:"$gte"`
//		Lte struct {
//			Date int64 `json:"$date"`
//		} `json:"$lte"`
//	} `json:"created_at"`
//	RackspaceCaused bool   `json:"rackspace_caused"`
//	IncidentCause   string `json:"incident_cause"`
//	Segment         string `json:"segment"`
//	Status          string `json:"status"`
//
//	Departments struct {
//		In []string `json:"$in"`
//	} `json:"departments"`
//}
//
//type QueryWithOnlyStartDate struct {
//	CreatedAt struct {
//		Gte struct {
//			Date int64 `json:"$date"`
//		} `json:"$gte"`
//	} `json:"created_at"`
//	RackspaceCaused bool   `json:"rackspace_caused"`
//	IncidentCause   string `json:"incident_cause"`
//	Segment         string `json:"segment"`
//	Status          string `json:"status"`
//	Departments     struct {
//		In []string `json:"$in"`
//	} `json:"departments"`
//}
