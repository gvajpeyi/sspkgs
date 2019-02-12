package creditmemo
//
//import (
//	"fmt"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/metricsdb"
//
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/jmoiron/sqlx"
//	"io/ioutil"
//	"strconv"
//	"strings"
//)
//
//var (
//	dbName string
//	dbUser string
//	dbIP   string
//	dbPort int
//	dbPass string
//	dbDSN  string
//)
//type CreditMemoDB metricsdb.MetricsDatabase
//
//func NewCreditMemoDatabase(mdb *metricsdb.MetricsDatabase) *CreditMemoDB {
//	var creditMemoDB CreditMemoDB
//	cstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mdb.DatabaseUser, mdb.DatabasePassword, mdb.IPAddress, mdb.Port, mdb.DatabaseName)
//
//	creditMemoDB.Port = mdb.Port
//	creditMemoDB.IPAddress = mdb.IPAddress
//	creditMemoDB.DatabaseName = mdb.DatabaseName
//	creditMemoDB.DatabaseUser = mdb.DatabaseUser
//	creditMemoDB.DatabasePassword = mdb.DatabasePassword
//	creditMemoDB.ConnectionString = cstring
//
//	fmt.Println(cstring)
//	return &creditMemoDB
//}
//func getDBConfig() ( string, error) {
//	content, err := ioutil.ReadFile("../config/env.txt")
//	if err != nil {
//		return "", err
//	}
//
//	fmt.Println("shouldn't see this")
//	lines := strings.Split(string(content), "\n")
//
//	for _, line := range lines {
//
//		envVariable := strings.Split(string(line), "=")
//
//		switch envVariable[0] {
//
//		case "AvMetricsDBIP":
//			dbIP = envVariable[1]
//
//		case "AvMetricsDBPort":
//			dbPort, _ = strconv.Atoi(envVariable[1])
//
//		case "AvMetricsDBUser":
//			dbUser = envVariable[1]
//
//		case "AvMetricsDBPass":
//			dbPass = envVariable[1]
//		case "AvMetricsDBName":
//			dbName = envVariable[1]
//		}
//	}
//	dbDSN = fmt.Sprintf("%v:%v@tcp(%v:%d)/%v", dbUser, dbPass, dbIP, dbPort, dbName)
//
//
//	return dbDSN, nil
//}
//
//
//func InsertMetricIntoDatabase(cmData *CreditMemoResponse) {
//	dbDSN, err := getDBConfig()
//	db, err := sqlx.Open("mysql", dbDSN)
//	if err != nil {
//		fmt.Printf("Error connection: %v\n", err)
//		return
//	}
//	defer db.Close()
//	err = db.Ping()
//	if err != nil {
//
//		fmt.Printf("Ping error: %v\n", err)
//		return
//	}
//
//	tx := db.MustBegin()
//
//	for _, cm := range cmData.Requests {
//		cdevList := ""
//
//		for _, dev := range cm.Request.Devices {
//			cdevList += fmt.Sprintf("%s, ", dev.ID)
//		}
//		lastComma := strings.LastIndex(cdevList, ",")
//		var devicesList string
//		if lastComma >= 0 {
//			devicesList = cdevList[:lastComma]
//		} else {
//			devicesList = ""
//		}
//
//		earnedRevenue := 1
//		if cm.Request.EarnedRevenue == false {
//			earnedRevenue = 0
//		}
//
//		isTicketPrivate := 1
//		if cm.Request.IsTicketPrivate == false {
//			isTicketPrivate = 0
//		}
//		sqlStr := `Insert into dc_metrics_credit_memo (
//            CMID                ,
//			Amount          ,
//			` + "`Type`" + `             ,
//			Team               ,
//			Calculation        ,
//			CurrencyCode  ,
//			InitiatorSso     ,
//			AccountManagerSso       ,
//			CustomerName            ,
//			Description            ,
//			DurationDays           ,
//			IncidentDate            ,
//			IncidentEventID         ,
//			Prevention              ,
//			RackspaceAccountNumber  ,
//			ReasonName             ,
//			AdditionalAmount      ,
//			ResponsibleDepartmentName  ,
//			ServiceFailureName      ,
//			Status           ,
//			WorkflowState    ,
//			CreatedDatetime ,
//			UpdatedDatetime,
//			ClosedDate    ,
//			LastApprovedBySso         ,
//			QaApprovedBySso         ,
//			IsTicketPrivate,
//			AccountTypeName ,
//			Segment    ,
//			SubSegment    ,
//			BusinessUnit  ,
//			EarnedRevenue,
//			Devices) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
//
//		// exec the schema or fail; multi-statement Exec behavior varies between
//		// database drivers;  pq will exec them all, sqlite3 won't, ymmv
//
//		db.Exec(sqlStr,
//			cm.Request.ID,
//			cm.Request.Amount,
//			cm.Request.Type,
//			cm.Request.Team,
//			cm.Request.Calculation,
//			cm.Request.Currency.Code,
//			cm.Request.Initiator.Sso,
//			cm.Request.AccountManager.Sso,
//			cm.Request.CustomerName,
//			cm.Request.Description,
//			cm.Request.DurationDays,
//			cm.Request.IncidentDate,
//			cm.Request.IncidentEventID,
//			cm.Request.Prevention,
//			cm.Request.RackspaceAccountNumber,
//			cm.Request.ReasonCode.Name,
//			cm.Request.AdditionalAmount,
//			cm.Request.ResponsibleDepartment.Name,
//			cm.Request.ServiceFailure.Name,
//			cm.Request.Status,
//			cm.Request.WorkflowState,
//			cm.Request.CreatedDatetime,
//			cm.Request.UpdatedDatetime,
//			cm.Request.ClosedDate,
//			cm.Request.LastApprovedBy.Sso,
//			cm.Request.QaApprovedBy.Sso,
//			isTicketPrivate,
//			cm.Request.AccountType.Name,
//			cm.Request.Segment,
//			cm.Request.SubSegment,
//			cm.Request.BusinessUnit,
//			earnedRevenue,
//			devicesList)
//
//
//	}
//	tx.Commit()
//
//}
//
//
//type CreditMemoTotalResult struct{
//	Department string `json:"department"`
//	FromDate string `json:"from_date"`
//	ToDate string `json:"to_date"`
//	CreditMemoCountCurrentMonth int `json:"credit_memo_count_cm"`
//	CreditMemoTotalCurrentMonth float64 `json:"credit_memo_total_cm"`
//
//	CreditMemoCountLastMonth int `json:"credit_memo_count_lm"`
//	CreditMemoTotalLastMonth float64 `json:"credit_memo_total_lm"`
//
//	CreditMemoCountLastQtr int `json:"credit_memo_count_lq"`
//	CreditMemoTotalLastQtr float64 `json:"credit_memo_total_lq"`
//
//	CreditMemoCountYTD int `json:"credit_memo_count_ytd"`
//	CreditMemoTotalYTD float64 `json:"credit_memo_total_ytd"`
//}
//
//type CreditMemoCurrentCounts struct{
//	Datacenter	string
//	CurrentMonthCount int
//	CurrentMonthAmount float64
//	Last30DayCount int
//	Last30DayAmount float64
//	Last60DayCount int
//	Last60DayAmount float64
//	Last90DayCount int
//	Last90DayAmount float64
//}
//
//
//
//func (mdb *CreditMemoDB)GetCurrentMonthCounts(dept string)(*CreditMemoCurrentCounts, error) {
//	var cmCounts CreditMemoCurrentCounts
//
//	var (
//
//		CurrentMonthCount []byte
//		CurrentMonthAmount []byte
//		Last30DayCount []byte
//		Last30DayAmount []byte
//		Last60DayCount []byte
//		Last60DayAmount []byte
//		Last90DayCount []byte
//		Last90DayAmount []byte
//
//
//	)
//	//dbDSN, err := getDBConfig()
//	db, err := sqlx.Open("mysql", mdb.ConnectionString)
//	if err != nil {
//		fmt.Printf("Error connection: %v\n", err)
//		return nil, err
//	}
//
//	defer db.Close()
//	err = db.Ping()
//	if err != nil {
//
//		fmt.Printf("Ping error: %v\n", err)
//		return nil, err
//	}
//
//
//	//
//	//sqlStatement := `Select cmid,  from dc_metrics_credit_memo
//	//				 where ClosedDate >= ?
//	//				 and ClosedDate <= ?
//	//				 and ResponsibleDepartmentName = ?`
//
//
//sqlStatement :=`
//Select Distinct (
//                Select count(cmid)
//	from dc_metrics_credit_memo
//	Where month(ClosedDate) = month(now()
//	and ResponsibleDepartmentName like (?))
//	) as current_month_count,
//
//                (Select sum(Amount)
//	from dc_metrics_credit_memo
//	Where month(ClosedDate) = month(now()
//	and ResponsibleDepartmentName like (?))
//	) as current_month_amount,
//
//		(SELECT count(cmid)
//	from dc_metrics_credit_memo
//	where ClosedDate between (now() + INTERVAL -30 DAY) and now()
//	and ResponsibleDepartmentName like (?)) as last_month_count,
//
//
//
//                (Select sum(Amount)
//	from dc_metrics_credit_memo
//	where ClosedDate between (now() + INTERVAL -30 DAY) and now()
//	and ResponsibleDepartmentName like (?)) as last_month_amount,
//
//
//
//		(SELECT count(cmid)
//	from dc_metrics_credit_memo
//	where DATE_SUB(ClosedDate, INTERVAL 1 QUARTER)
//and ResponsibleDepartmentName like (?))as qtd_count,
//
//
//                (Select sum(Amount)
//	from dc_metrics_credit_memo
//	where DATE_SUB(ClosedDate, INTERVAL 1 QUARTER)
//and ResponsibleDepartmentName like (?))as qtd_amount,
//
//
//
//		(SELECT count(cmid)
//	from dc_metrics_credit_memo
//	where year(ClosedDate) =  year(now())
//	and ResponsibleDepartmentName like (?)) as ytd_count,
//
//
//                (Select sum(Amount)
//	from dc_metrics_credit_memo
//	where year(ClosedDate) =  year(now())
//	and ResponsibleDepartmentName like (?)) as ytd_amount
//
//
//	from dc_metrics_credit_memo
//	where ResponsibleDepartmentName like (?)
//
//`
//
//deptWildcard := fmt.Sprintf("%%%s", dept)
//
//fmt.Println(sqlStatement)
//fmt.Println(dept)
//		rows, err := db.Query(sqlStatement, deptWildcard,deptWildcard,deptWildcard,deptWildcard,deptWildcard,deptWildcard,deptWildcard,deptWildcard,deptWildcard)
//
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var rowCounter= 0
//	if rows == nil {
//		return nil, err
//	} else {
//		for rows.Next() {
//			rowCounter++
//			err = rows.Scan(
//				&CurrentMonthCount,
//			&CurrentMonthAmount,
//			&Last30DayCount,
//			&Last30DayAmount ,
//			&Last60DayCount ,
//			&Last60DayAmount,
//			&Last90DayCount,
//			&Last90DayAmount,
//			)
//
//			currentAmount, _ := strconv.ParseFloat(string(CurrentMonthAmount), 64)
//			currentCount, _ := strconv.Atoi(string(CurrentMonthCount))
//
//			last30Amount, _ := strconv.ParseFloat(string(Last30DayAmount), 64)
//			last30Count, _ := strconv.Atoi(string(Last30DayCount))
//
//
//			last60Amount, _ := strconv.ParseFloat(string(Last60DayAmount), 64)
//			last60Count, _ := strconv.Atoi(string(Last60DayCount))
//
//
//			last90Amount, _ := strconv.ParseFloat(string(Last90DayAmount), 64)
//			last90Count, _ := strconv.Atoi(string(Last90DayCount))
//
//			cmCounts = CreditMemoCurrentCounts{
//				Datacenter: dept,
//				CurrentMonthAmount: currentAmount,
//				CurrentMonthCount: currentCount,
//				Last30DayAmount: last30Amount,
//				Last30DayCount: last30Count,
//				Last60DayCount: last60Count,
//				Last60DayAmount: last60Amount,
//				Last90DayCount: last90Count,
//				Last90DayAmount: last90Amount,
//			}
//			if err != nil {
//				return nil, err
//			}
//
//
//		}
//		return &cmCounts, nil
//
//	}
//}