package hitch
//
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/metricsdb"
//	"sort"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type HitchReport struct {
//	TicketNumber        string
//	RecordInserted      string
//	CustomerTemperature string
//	OutageStatus        string
//	StartDate           string
//	ShortDescription    string
//	Segment             string
//	RackspaceCaused     int
//	OutageType          string
//	OutageCause         string
//	OutageImpact        string
//	LastEmailSent       string
//	Escalated           int
//	EndDate             string
//	Departments         string
//	DataCenter          string
//	CustomerCaused      int
//	CreatedByName       string
//	CreatedByEmail      string
//	Time                string
//	AccountNumber       string
//	AccountName         string
//	AccountSegment      string
//	SupportTeam         string
//}
//
//type HitchDB metricsdb.MetricsDatabase
//
//type MetricCurrentCount struct {
//	Datacenter             string
//	Metric                 string
//	CurrentMonthCount      int
//	CurrentMonthChangeRate float32
//	Last30DayChangeRate    float32
//	Last60DayChangeRate    float32
//	Last90DayChangeRate    float32
//}
//
//func NewHitchDatabase(mdb *metricsdb.MetricsDatabase) *HitchDB {
//	var hitchDB HitchDB
//	cstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mdb.DatabaseUser, mdb.DatabasePassword, mdb.IPAddress, mdb.Port, mdb.DatabaseName)
//
//	hitchDB.Port = mdb.Port
//	hitchDB.IPAddress = mdb.IPAddress
//	hitchDB.DatabaseName = mdb.DatabaseName
//	hitchDB.DatabaseUser = mdb.DatabaseUser
//	hitchDB.DatabasePassword = mdb.DatabasePassword
//	hitchDB.ConnectionString = cstring
//
//	return &hitchDB
//}
//
//func hitchMetricCountQuery(connectionString string, selectStatement string, dc string) (*HitchCounts, error) {
//	fmt.Println("Got: ", selectStatement)
//	var hitchCounts *HitchCounts
//	if connectionString == "" {
//		return nil, fmt.Errorf("Connection String is empty")
//
//	}
//	fmt.Println("Conn string: ", connectionString)
//
//	db, err := sql.Open("mysql", connectionString)
//	if err != nil {
//		fmt.Println("openDB Error with database ", err)
//		return nil, err
//	}
//	defer db.Close()
//	fmt.Println("db open not pinged yet")
//	err = db.Ping()
//	if err != nil {
//		fmt.Println("db ping error: ", err.Error())
//		return nil, err
//	} else {
//		fmt.Println("Ping good")
//	}
//
//	var (
//		CurrentMonth []byte
//		LastMonth    []byte
//		LastQuarter  []byte
//		LastYear     []byte
//	)
//
//	var results *sql.Rows
//
//	results, err = db.Query(selectStatement, dc, dc, dc, dc, dc)
//
//	if err != nil {
//		//fmt.Println("Error: ", dbError)
//		return nil, fmt.Errorf("%+v:\n %v:\n%s\n", selectStatement, err, err.Error())
//	}
//	fmt.Println("no db errors")
//	defer results.Close()
//
//	for results.Next() {
//		err = results.Scan(
//			&CurrentMonth,
//			&LastMonth,
//			&LastQuarter,
//			&LastYear,
//		)
//		if err != nil {
//			errf := fmt.Errorf("DB: Row Scan Failed: %s", err.Error())
//			return nil, errf
//		}
//
//		currentMonth, _ := strconv.Atoi(string(CurrentMonth))
//		lastMonth, _ := strconv.Atoi(string(LastMonth))
//		lastQuarter, _ := strconv.Atoi(string(LastQuarter))
//		yearToDate, _ := strconv.Atoi(string(LastYear))
//
//		hitchCounts = &HitchCounts{
//			CurrentMonth: currentMonth,
//			LastMonth:    lastMonth,
//			LastQuarter:  lastQuarter,
//			YearToDate:   yearToDate,
//		}
//
//	}
//	if hitchCounts != nil {
//		return hitchCounts, nil
//	}
//	return nil, fmt.Errorf("No count data found")
//
//}
//
//type HitchCounts struct {
//	CurrentMonth int
//	LastMonth    int
//	LastQuarter  int
//	YearToDate   int
//}
//
//// HitchReportQuery executes a select query against the metrics db
//// hitch table and returns a slice of hitch reports with all fields.
//func hitchReportQuery(connectionString string, selectStatement string, dc string, ticketNumber *string) ([]*HitchReport, error) {
//	hitchReports := []*HitchReport{}
//
//	fmt.Println("connection string: ", connectionString)
//
//	if connectionString == "" {
//		return nil, fmt.Errorf("Connection String is empty")
//
//	}
//	db, err := sql.Open("mysql", connectionString)
//	if err != nil {
//		fmt.Println("Error: ", err)
//		return nil, err
//	}
//
//	err = db.Ping()
//	if err != nil {
//		fmt.Println("db ping error: ", err.Error())
//		return nil, err
//	} else {
//		fmt.Println("Ping good")
//	}
//	defer db.Close()
//	var (
//		TicketNumber        []byte
//		RecordInserted      []byte
//		CustomerTemperature []byte
//		OutageStatus        []byte
//		StartDate           []byte
//		ShortDescription    []byte
//		Segment             []byte
//		RackspaceCaused     []byte
//		OutageType          []byte
//		OutageCause         []byte
//		OutageImpact        []byte
//		LastEmailSent       []byte
//		Escalated           []byte
//		EndDate             []byte
//		Departments         []byte
//		DataCenter          []byte
//		CustomerCaused      []byte
//		CreatedByName       []byte
//		CreatedByEmail      []byte
//		Time                []byte
//		AccountNumber       []byte
//		AccountName         []byte
//		AccountSegment      []byte
//		SupportTeam         []byte
//	)
//
//	var results *sql.Rows
//	var dbError error
//
//	if ticketNumber != nil {
//		results, dbError = db.Query(selectStatement, ticketNumber, dc)
//	} else {
//		results, dbError = db.Query(selectStatement, dc)
//
//	}
//	if dbError != nil {
//		fmt.Println("Error: ", dbError)
//		return nil, fmt.Errorf("%+v:\n %v:\n%s\n", selectStatement, dbError, err.Error())
//	}
//
//	defer results.Close()
//
//	for results.Next() {
//		err = results.Scan(
//			&TicketNumber,
//			&RecordInserted,
//			&CustomerTemperature,
//			&OutageStatus,
//			&StartDate,
//			&ShortDescription,
//			&Segment,
//			&RackspaceCaused,
//			&OutageType,
//			&OutageCause,
//			&OutageImpact,
//			&LastEmailSent,
//			&Escalated,
//			&EndDate,
//			&Departments,
//			&DataCenter,
//			&CustomerCaused,
//			&CreatedByName,
//			&CreatedByEmail,
//			&Time,
//			&AccountNumber,
//			&AccountName,
//			&AccountSegment,
//			&SupportTeam,
//		)
//
//		if err != nil {
//			errf := fmt.Errorf("DB: Row Scan Failed: %s", err.Error())
//			return nil, errf
//		}
//
//		raxCaused, _ := strconv.Atoi(string(RackspaceCaused))
//		escalated, _ := strconv.Atoi(string(Escalated))
//		ccaused, _ := strconv.Atoi(string(CustomerCaused))
//
//		hitchReport := HitchReport{
//			TicketNumber:        string(TicketNumber),
//			RecordInserted:      string(RecordInserted),
//			CustomerTemperature: string(CustomerTemperature),
//			OutageStatus:        string(OutageStatus),
//			StartDate:           string(StartDate),
//			ShortDescription:    string(ShortDescription),
//			Segment:             string(Segment),
//			RackspaceCaused:     raxCaused,
//			OutageType:          string(OutageType),
//			OutageCause:         string(OutageCause),
//			OutageImpact:        string(OutageImpact),
//			LastEmailSent:       string(LastEmailSent),
//			Escalated:           escalated,
//			EndDate:             string(EndDate),
//			Departments:         string(Departments),
//			DataCenter:          string(DataCenter),
//			CustomerCaused:      ccaused,
//			CreatedByName:       string(CreatedByName),
//			CreatedByEmail:      string(CreatedByEmail),
//			Time:                string(Time),
//			AccountNumber:       string(AccountNumber),
//			AccountName:         string(AccountName),
//			AccountSegment:      string(AccountSegment),
//			SupportTeam:         string(SupportTeam),
//		}
//
//		hitchReports = append(hitchReports, &hitchReport)
//
//	}
//	return hitchReports, nil
//
//}
//
//func (hdb *HitchDB) GetCurrentMonthCounts(dc string) (*HitchCounts, error) {
//
//	sqlQuery := `Select Distinct (Select count(ticket_number)
//	from dc_metrics_hitch
//	Where month(end_date) = month(now()
//	and data_center = ?)
//	) as current_month,
//
//		(SELECT count(ticket_number)
//	from dc_metrics_hitch
//	where end_date between (now() + INTERVAL -30 DAY) and now()
//	and data_center = ?) as lastmonth,
//
//		(SELECT count(ticket_number)
//	from dc_metrics_hitch
//	where DATE_SUB(end_date, INTERVAL 1 QUARTER)
//and data_center = ?)as qtd,
//
//		(SELECT count(ticket_number)
//	from dc_metrics_hitch
//	where year(end_date) =  year(now())
//	and data_center = ?) as ytd
//
//	from dc_metrics_hitch
//	where data_center = ?
//`
//	fmt.Println("Sending: ", sqlQuery)
//
//	counts, err := hitchMetricCountQuery(hdb.ConnectionString, sqlQuery, dc)
//	if err != nil {
//		return nil, fmt.Errorf(fmt.Sprintf("%s;;\nError returned from hitchMetricCountQuery()", err.Error()))
//	}
//
//	return counts, nil
//
//}
//
//func (hdb *HitchDB) GetHitchIncidents(dc string, ticketNumber *string) ([]*HitchReport, *Chart, error) {
//
//	var sqlQuery string
//	var hitchResults []*HitchReport
//	var err error
//
//	if ticketNumber == nil {
//		sqlQuery = `Select
//ticket_number,
//record_inserted,
//temperature,
//outage_status,
//start_date,
//short_description,
//segment,
//rackspace_caused,
//outage_type,
//outage_cause,
//outage_impact,
//last_email_sent,
//escalated,
//end_date,
//departments,
//data_center,
//customer_caused,
//created_by_name,
//created_by_email,
//time,
//account_number,
//account_name,
//account_segment,
//support_team
//from dc_metrics_hitch
//where data_center = ?`
//
//		fmt.Printf("sql: %s\n\n", sqlQuery)
//		hitchResults, err = hitchReportQuery(hdb.ConnectionString, sqlQuery, dc, nil)
//		if err != nil {
//			return nil, nil, fmt.Errorf(fmt.Sprintf("%s;;\nError returned from hitchReportQuery()", err.Error()))
//		}
//
//	} else {
//
//		sqlQuery = `Select
//ticket_number,
//record_inserted,
//temperature,
//outage_status,
//start_date,
//short_description,
//segment,
//rackspace_caused,
//outage_type,
//outage_cause,
//outage_impact,
//last_email_sent,
//escalated,
//end_date,
//departments,
//data_center,
//customer_caused,
//created_by_name,
//created_by_email,
//time,
//account_number,
//account_name,
//account_segment,
//support_team
//from dc_metrics_hitch
//where ticket_number = ?
//and data_center = ?`
//
//		hitchResults, err = hitchReportQuery(hdb.ConnectionString, sqlQuery, dc, ticketNumber)
//
//		if err != nil {
//			fmt.Println("Returning Error")
//
//			return nil, nil, fmt.Errorf(fmt.Sprintf("%s;;\nError returned from hitchReportQuery()", err.Error()))
//		}
//	}
//
//	chartData := formatChartData(hitchResults, "bar", "Racker Error")
//
//	return hitchResults, chartData, nil
//
//}
//
//type ByEndDate []*HitchReport
//
//func (s ByEndDate) Len() int {
//	return len(s)
//
//}
//func (s ByEndDate) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}
//func (s ByEndDate) Less(i, j int) bool {
//	return len(s[i].EndDate) < len(s[j].EndDate)
//}
//func formatChartData(results []*HitchReport, chartType string, title string) *Chart {
//	sort.Sort(ByEndDate(results))
//
//	chart := Chart{}
//	data := Data{}
//
//	chart.Type = chartType
//	chart.Options.Responsive = true
//	chart.Options.Title.Display = true
//	chart.Options.Title.Text = title
//
//	borderColor := []string{
//		"rgb(46, 204, 64, 1),rgb(0, 116, 217, 1),rgb(255, 220, 0, 1),rgb(177, 13, 201, 1),rgb(57, 204, 204, 1),rgb(61, 153, 112, 1),rgb(255, 133, 27, 1),rgb(0, 31, 63, 1),rgb(133, 20, 75, 1),rgb(255, 65, 54, 1),rgb(1, 255, 112, 1),rgb(127, 219, 255, 1)",
//	}
//
//	backgroundColor := []string{
//		"rgb(46,204,64,.2),rgb(0,116,217,.1),rgb(255,220,0,.1),rgb(177,13,201,.1),rgb(57,204,204,.1),rgb(61,153,112,.1),rgb(255,133,27,.1),rgb(0,31,63,.1),rgb(133,20,75,.1),rgb(255,65,54,.1),rgb(1,255,112,.1), rgb(127,219,255,.1)",
//	}
//
//	dataset := Dataset{
//		BackgroundColor: backgroundColor,
//		BorderColor:     borderColor,
//		Label:           "Hitch Reports",
//	}
//
//	currYear := time.Now().Year()
//
//	jan := 0
//	feb := 0
//	mar := 0
//	apr := 0
//	may := 0
//	jun := 0
//	jul := 0
//	aug := 0
//	sep := 0
//	oct := 0
//	nov := 0
//	dec := 0
//
//	for _, result := range results {
//
//		dbTime := strings.Replace(result.EndDate, " ", "T", 1)
//		dbTime = fmt.Sprintf("%s+00:00", dbTime)
//		tm, _ := time.Parse(time.RFC3339, dbTime)
//
//		if tm.Year() == currYear {
//
//			if tm.Month() == time.January {
//				jan++
//			}
//			if tm.Month() == time.February {
//				feb++
//			}
//			if tm.Month() == time.March {
//				mar++
//			}
//			if tm.Month() == time.April {
//				apr++
//			}
//			if tm.Month() == time.May {
//				may++
//			}
//			if tm.Month() == time.June {
//				jun++
//			}
//			if tm.Month() == time.July {
//				jul++
//			}
//			if tm.Month() == time.August {
//				aug++
//			}
//			if tm.Month() == time.September {
//				sep++
//			}
//			if tm.Month() == time.October {
//				oct++
//			}
//			if tm.Month() == time.November {
//				nov++
//			}
//			if tm.Month() == time.December {
//				dec++
//			}
//
//		}
//
//	}
//
//	dataset.Data = []int{jan, feb, mar, apr, may, jun, jul, aug, sep, oct, nov, dec}
//	data.Labels = []string{"J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"}
//
//	data.Datasets = append(data.Datasets, dataset)
//	chart.ChartData = data
//	return &chart
//
//}
//
//type Chart struct {
//	Type      string `json:"type"`
//	ChartData Data   `json:"data"`
//	Options   struct {
//		Responsive bool `json:"responsive"`
//		Title      struct {
//			Display bool   `json:"display"`
//			Text    string `json:"text"`
//		} `json:"title"`
//	} `json:"options"`
//}
//
//type Data struct {
//	Datasets []Dataset `json:"datasets"`
//	Labels   []string  `json:"labels"`
//}
//
//type Dataset struct {
//	Label           string   `json:"label"`
//	Data            []int    `json:"data"`
//	BackgroundColor []string `json:"backgroundColor"`
//	BorderColor     []string `json:"borderColor"`
//	Type            string   `json:"type,omitempty"`
//}
