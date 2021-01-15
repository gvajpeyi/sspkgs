// coreods provides access to the EBI Core_ODS database without having to know the schema.
package coreods

import (
	"database/sql"
	"strings"

	// blank import as required
	_ "database/sql"
	"fmt"
	"strconv"
	"time"

	// blank import as required
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/jmoiron/sqlx"
)

// ODSConfig contains the needed info to create a connection to the CORE_ODS database
type ODSConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	User             string `json:"user"`
	Password         string `json:"password"`
	CoreProdPassword string `json:"CoreProdPassword,omitempty"`
	DBName           string `json:"db_name"`
}

// ConnectionInfo returns the connection string generated from the ODSConfig data
func (c ODSConfig) ConnectionInfo() string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DBName)

}

// ODSService is an interface accessing the ODSDB interface
type ODSService interface {
	ODSDB
}

// ODSDB interface provides the needed functions to work with the CoreODS database
type ODSDB interface {
	Ping() (bool, error)
	DataCenterServerCount(start, end time.Time, dc string) (DCServerCount, error)
	ExchangeRate(convertFromCurrency string, month int, year int) (float64, error)
	DeviceDetails(dt DeviceType, ds DeviceStatus, deviceList *string) (*[]DeviceDetail, error)
	LoadAllDevicesAndAccountNums() (*[]DevAcct, error)
	GetAllActiveDevices() (*[]DeviceDetail, error)
}

// NewODSService creates a new service that implements the needed interfaces and provides access to the CoreODS database
func NewODSService(config ODSConfig) (ODSService, error) {

	connStr := fmt.Sprintf("odbc:server=%s; port=%d; user id=%s;password=%s; database=%s;log=3;encrypt=false;TrustServerCertificate=true", config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sqlx.Open("mssql", connStr)

	if err != nil {

		return nil, err
	}

	err = db.Ping()
	if err != nil {

		return nil, err
	}

	return &odsService{
		ODSDB: &odsDB{db},
	}, nil
}

type odsService struct {
	ODSDB
}

var _ ODSDB = &odsDB{}

type odsDB struct {
	db *sqlx.DB
}

func (ods *odsDB) Ping() (bool, error) {

	err := ods.db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ods *odsDB) ExchangeRate(convertFromCurrency string, month int, year int) (float64, error) {
	// exchange rate
	type ExchangeRate struct {
		Rate string
	}

	qry := `Select Max(e.ExchangeRate) as rate
	from Corporate_DMART.dbo.vw_BI_Summa_ExchangeRate e
	where e.[Currency Selection] = ?
and e.[Credit Memo Month] = ?
	and e.[Credit Memo Year] = ?`

	erate := ExchangeRate{}
	err := ods.db.QueryRowx(qry, convertFromCurrency, month, year).StructScan(&erate)

	if err != nil {
		return 0.0, err
	}
	erateF := castStringToFloat64(erate.Rate)

	return erateF, nil

}

func castStringToFloat64(strNumber string) float64 {
	erateF, err := strconv.ParseFloat(strNumber, 64)
	if err != nil {
		return 0.0
	}
	return erateF
}

type DeviceType int

const (
	VMOnly DeviceType = iota
	NonVMOnly
	VMandNonVM
)

func (dt DeviceType) String() string {
	return [...]string{"VMOnly", "NotVMOnly", "VMandNonVM"}[dt]
}

type DeviceStatus int

const (
	ExcludeOffline = iota
	AllStatuses
)

func (ds DeviceStatus) String() string {
	return [...]string{"ExcludeOffline", "AllStatuses"}[ds]
}

// DeviceDetail provides device specific information relating to a device as returned from CoreODS database.
type DeviceDetail struct {
	Status             string  `json:"status" db:"status"`
	ProductDescription string  `json:"product_description" db:"product_description"`
	Datacenter         string  `json:"datacenter" db:"datacenter_abbr"`
	ComputerNumber     string  `json:"computer_number" db:"computer_number"`
	CustomerNumber     string  `json:"customer_number" db:"customer_number"`
	ServerName         *string `json:"server_name" db:"server_name"`
	IPAddress          *string `json:"ip_address" db:"ip_address"`
}

func addPrefix(firstFilter bool, sqlFilters string) string {

	if firstFilter {
		sqlFilters += "Where "

	} else {
		sqlFilters += " And "
	}

	return sqlFilters
}

func generateSqlQuery(dt DeviceType, ds DeviceStatus, deviceList *string) string {

	qry := `select so.status, sku.product_description, dc.datacenter_abbr, s.computer_number, s.customer_number, s.server_name, ip.ip_address from CORE_Prod.server as s inner join CORE_Prod.datacenter as dc on s.datacenter_number = dc.datacenter_number inner join CORE_Prod.status_options as so on s.status_number = so.status_number inner join CORE_Prod.sku as sku on s.product_id = sku.product_sku`

	qry = `select so.status,
       sku.product_description,
       dc.datacenter_abbr,
       s.computer_number,
       s.customer_number,
       s.server_name,
       ip.IPAddress as 'ip_address'
from CORE_Prod.server as s
         inner join CORE_Prod.datacenter as dc on s.datacenter_number = dc.datacenter_number
         inner join CORE_Prod.status_options as so on s.status_number = so.status_number
         inner join CORE_Prod.sku as sku on s.product_id = sku.product_sku
        inner join CORE_Prod.IPSP_cache_IPAssignment as ip on s.computer_number = ip.computer_number
`

	firstFilter := true
	sqlFilters := ""
	var dl string
	if deviceList != nil {
		dl = *deviceList
	} else {
		dl = ""
	}
	if dl != "" {
		sqlFilters = fmt.Sprintf("%s %s ", addPrefix(firstFilter, sqlFilters), "s.computer_number in (?)")
		firstFilter = false
	}

	switch dt {
	case VMOnly:
		sqlFilters = fmt.Sprintf("%s %s", addPrefix(firstFilter, sqlFilters), " (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%')) ")
		firstFilter = false
	case NonVMOnly:
		sqlFilters = fmt.Sprintf("%s %s", addPrefix(firstFilter, sqlFilters), " (sku.product_description not like ('%Virtual%') AND sku.product_description not like ('%VM%')) ")
		firstFilter = false
	case VMandNonVM:
		// no filters needed here.  default behavior for the query is to not filter
		// sqlFilters = fmt.Sprintf("%s", sqlFilters)
		//firstFilter = false
	}

	switch ds {
	case ExcludeOffline:

		sqlFilters = fmt.Sprintf("%s %s", addPrefix(firstFilter, sqlFilters), "so.status_number != -1")

		firstFilter = false
	}

	sqlQry := fmt.Sprintf("%s %s order by s.computer_number", qry, sqlFilters)

	return sqlQry

}

func (ods *odsDB) DeviceDetails(dt DeviceType, ds DeviceStatus, deviceList *string) (*[]DeviceDetail, error) {

	qry := generateSqlQuery(dt, ds, deviceList)
	devDetail := DeviceDetail{}
	devDetails := []DeviceDetail{}
	var err error
	var rows *sqlx.Rows
	if deviceList == nil {
		rows, err = ods.db.Queryx(qry)

	} else {
		rows, err = ods.db.Queryx(qry, *deviceList)
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {

		err := rows.StructScan(&devDetail)

		if err != nil {
			return nil, err
		}
		devDetails = append(devDetails, devDetail)
	}

	//sort by device number
	//sort.Slice(devDetails, func(i, j int) bool { return devDetails[i].ComputerNumber < devDetails[j].ComputerNumber })
	// a device with multiple IPaddresses will each be on a unique row.
	// this block will combined multiple ips into a csv string on a single row

	if len(devDetails) < 1 {
		return nil, fmt.Errorf("No devices found")
	}

	currentDeviceNumber := devDetails[len(devDetails)-1].ComputerNumber
	for i := len(devDetails) - 1; i <= len(devDetails); i-- {
		if i == 0 {
			break
		}
		nextDeviceNumber := devDetails[i-1].ComputerNumber
		if currentDeviceNumber == nextDeviceNumber {
			if devDetails[i].IPAddress != nil && *devDetails[i].IPAddress != "" {
				ips := fmt.Sprintf("%s, %s", *devDetails[i].IPAddress, *devDetails[i-1].IPAddress)
				devDetails[i].IPAddress = &ips
			}
			devDetails = append(devDetails[:i-1], devDetails[i:]...)
		} else {
			currentDeviceNumber = nextDeviceNumber

		}
	}
	return &devDetails, err

}

func (ods *odsDB) GetAllActiveDevices() (*[]DeviceDetail, error) {

	// setup individual queries for each table.  much faster grabbing all the data and them
	// throwing together here versus letting db server do it.
	serverQry := `select s.computer_number
                    ,s.customer_number
                    ,s.server_name
                    ,s.datacenter_number
					,s.status_number
					,s.product_id
				from CORE_Prod.server as s WITH (NOLOCK)
				where s.status_number != -1;`

	statusOptionsQry := `select so.status
                            ,so.status_number
						from CORE_Prod.status_options as so WITH (NOLOCK)
						where so.status_number != -1;`

	skuQry := `select sku.product_description
				,sku.product_sku
			  from CORE_Prod.sku as sku;`

	dataCenterQry := `select  dc.datacenter_abbr
                        ,dc.datacenter_number
                      from CORE_Prod.datacenter as dc WITH (NOLOCK);`

	ipaddressQry := `select ip.IPAddress as 'ip_address',
                        ip.computer_number
					from  CORE_Prod.IPSP_cache_IPAssignment as ip WITH (NOLOCK);`

	// building out the structs and maps to send of for processing.
	// using mutexes per map to avoid race
	var srvr server
	servers := make(map[int]server)

	var ipa ipAddr
	ipAddresses := make(map[int][]string)

	var dc dataCenter
	dataCenters := make(map[int]dataCenter)

	var s sku
	skus := make(map[int]sku)

	var so statusOption
	statusOptions := make(map[int]statusOption)

	// build out a map of queries and using the returned data type name as the key.  will come
	// in handy when needing to scan row
	queries := make(map[string]string)
	queries["ipAddr"] = ipaddressQry
	queries["dataCenter"] = dataCenterQry
	queries["sku"] = skuQry
	queries["statusOption"] = statusOptionsQry
	queries["server"] = serverQry

	for qType, query := range queries {
		qType := qType

		//fmt.Printf("%s:  Querying Table: %s \n", time.Now().String(), qType)

		rows, err := ods.db.Queryx(query)
		if err != nil {
			//fmt.Println(fmt.Errorf("ods.GetAllActiveDevices(): %w", err))
		}

		for rows.Next() {

			switch qType {

			case "server":

				err := rows.StructScan(&srvr)
				if err != nil {
					return &[]DeviceDetail{}, fmt.Errorf("ods.GetAllActiveDevices().StructScan().srvr: %w", err)
				}

				srvrNum, err := strconv.Atoi(srvr.ComputerNumber)
				if err != nil {
					return &[]DeviceDetail{}, fmt.Errorf("invalid computer number   ods.GetAllActiveDevices(): %w", err)
				}

				servers[srvrNum] = srvr

			case "ipAddr":

				err := rows.StructScan(&ipa)
				if err != nil {
					return &[]DeviceDetail{}, fmt.Errorf("ods.GetAllActiveDevices().StructScan().ip: %w", err)
				}

				//fmt.Println("ip address: ", ipa)
				// ipAddresses[ipa.Addr] = append(ipAddresses[ipa.Addr],  ipa.Addr)

				ipAddresses[ipa.ComputerNumber] = append(ipAddresses[ipa.ComputerNumber], ipa.Addr)

			case "dataCenter":

				err := rows.StructScan(&dc)
				if err != nil {
					return &[]DeviceDetail{}, fmt.Errorf("ods.GetAllActiveDevices().StructScan().dc: %w", err)
				}
				dataCenters[dc.Number] = dc
			case "sku":

				err := rows.StructScan(&s)

				if err != nil {
					return &[]DeviceDetail{}, fmt.Errorf("ods.GetAllActiveDevices().StructScan().sku: %w", err)
				}

				skus[s.ProductSku] = s
			case "statusOption":

				err := rows.StructScan(&so)

				if err != nil {
					return &[]DeviceDetail{}, fmt.Errorf("ods.GetAllActiveDevices().StructScan().so: %w", err)
				}
				statusOptions[so.StatusNumber] = so

			default:

			}

		}
		if err != nil {
			//fmt.Println(fmt.Errorf("ods.GetAllActiveDevices(): %w", err))
		}

	} // query for loop

	//fmt.Printf("status options: %d\nskus: %d\ndcs: %d\nips: %d\nserver: %d\n", len(statusOptions), len(skus), len(dataCenters), len(ipAddresses), len(servers))

	var devDetails []DeviceDetail

	//fmt.Printf("\n%s:  combining query results: ", time.Now().String())

	for _, v := range servers {
		cn, err := strconv.Atoi(v.ComputerNumber)
		if err != nil {
			return nil, err
		}

		ipString := strings.Join(ipAddresses[cn], ",")

		//fmt.Printf("ipAddresses map: %+v\n",ipString)

		//ipAddrString := strings.Join(ipAddrs, ",")
		//	fmt.Println("ipaddrstring: ", ipAddrString)
		//ipAddrString := ipAddresses[v.ComputerNumber].Addr

		dd := DeviceDetail{
			Status:         statusOptions[v.StatusNumber].Status,
			Datacenter:     dataCenters[v.DatacenterNumber].Abbr,
			ComputerNumber: v.ComputerNumber,
			CustomerNumber: v.CustomerNumber,
			ServerName:     v.ServerName,
			IPAddress:      &ipString,
		}

		if v.ProductID.Valid {
			dd.ProductDescription = skus[int(v.ProductID.Int32)].ProductDescription
		}

		devDetails = append(devDetails, dd)

	}
	//fmt.Printf(" --  completed in %s\n", t.Sub(startC).String())

	//fmt.Printf("Devices found: %d\n", len(devDetails))
	return &devDetails, nil

}

type server struct {
	StatusNumber     int           `json:"status_number" db:"status_number"`
	DatacenterNumber int           `json:"datacenter_number" db:"datacenter_number"`
	ComputerNumber   string        `json:"computer_number" db:"computer_number"`
	CustomerNumber   string        `json:"customer_number" db:"customer_number"`
	ServerName       *string       `json:"server_name" db:"server_name"`
	ProductID        sql.NullInt32 `db:"product_id"`
}

type statusOption struct {
	Status       string `db:"status"`
	StatusNumber int    `db:"status_number"`
}

type sku struct {
	ProductSku         int    `db:"product_sku"`
	ProductDescription string `db:"product_description"`
}

type dataCenter struct {
	Number int    `db:"datacenter_number"`
	Abbr   string `db:"datacenter_abbr"`
}

type ipAddr struct {
	Addr           string `db:"ip_address"`
	ComputerNumber int    `db:"computer_number"`
}
