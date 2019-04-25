package coreods

import (
	_ "database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/jmoiron/sqlx"
)

type ODSConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	User             string `json:"user"`
	Password         string `json:"password"`
	CoreProdPassword string `json:"CoreProdPassword,omitempty"`

	DBName string `json:"db_name"`
}

func (c ODSConfig) ConnectionInfo() string {
	if c.Password == "" {

		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DBName)

	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DBName)

}

type ODSService interface {
	ODSDB
}

type ODSDB interface {
	Ping() (bool, error)
	DataCenterServerCount(start, end time.Time, dc string) (DCServerCount, error)
	ExchangeRate(convertFromCurrency string, month int, year int) (float64, error)
	DeviceDetails(deviceList string) (*[]DeviceDetail, error)
}

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
	erateF, err := strconv.ParseFloat(string(erate.Rate), 64)

	if err != nil {
		return 0.0, err
	}
	return erateF, nil

}

type DeviceDetail struct {
	Status             string
	ProductDescription string
	Datacenter         string
	ComputerNumber     string
	CustomerNumber     string
	ServerName         string
}

func (ods *odsDB) DeviceDetails(deviceList string) (*[]DeviceDetail, error) {
	// exchange rate

	qry := `select so.status, sku.product_description, dc.datacenter_abbr, s.computer_number, s.customer_number, s.server_name
from CORE_Prod.server as s
       inner join CORE_Prod.datacenter as dc
                  on s.datacenter_number = dc.datacenter_number
       inner join CORE_Prod.status_options as so
                  on s.status_number = so.status_number
       inner join CORE_Prod.sku as sku
                  on s.product_id = sku.product_sku
Where sku.product_description not like ('%Virtual%')
		and sku.product_description not like ('%VM%')
		and s.computer_number in (?)`

	devDetails := new([]DeviceDetail)
	err := ods.db.QueryRowx(qry, deviceList).StructScan(&devDetails)

	if err != nil {
		return nil, err
	}

	return devDetails, err

}
