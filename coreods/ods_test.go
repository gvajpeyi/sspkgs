package coreods

import (
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/DATA-DOG/go-sqlmock"
	// _ "github.com/denisenkom/go-mssqldb"
	// _ "github.com/denisenkom/go-mssqldb"
	// _ "github.com/denisenkom/go-mssqldb"
	// _ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

func getNewMockDB() Mocker {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return Mocker{}
	}
	return Mocker{
		db:   mockDB,
		mock: mock,
	}
}

func gimmeAPointer(s string) *string {
	return &s
}
func setupODSTests() (ODSService, error) {
	odsUsername := os.Getenv("ODSID")
	odsPassword := os.Getenv("ODSPW")
	odsHost := os.Getenv("ODSHOST")
	odsPortS := os.Getenv("ODSPORT")
	//os.Getenv("CoreProdPassword")
	//comment

	odsPort, err := strconv.Atoi(odsPortS)
	if err != nil {
		return nil, err
	}

	if odsUsername == "" || odsPassword == "" || odsPortS == "" {
		return nil, fmt.Errorf("database env variables not configured")
	}
	dbConfig := ODSConfig{
		Host:     odsHost,
		Port:     odsPort,
		User:     odsUsername,
		Password: odsPassword,
		DBName:   "Operational_reporting_CORE",
	}

	ods, err := NewODSService(dbConfig)

	return ods, fmt.Errorf("setupODSTests: %w", err)

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
		Host:     odsHost,
		Port:     odsPort,
		User:     odsUsername,
		Password: odsPassword,
		DBName:   "Corporate_DMART",
	}

	ods, err := NewODSService(dbConfig)
	return ods, err

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
		name    string
		T       TestArgs
		want    float64
		wantErr bool
	}{

		{name: "Valid AUD", T: TestArgs{"AUD", 1, 2019}, want: 1.39971757985588, wantErr: false},
		{name: "Valid EUR", T: TestArgs{"EUR", 1, 2019}, want: 0.879225, wantErr: false},
		{name: "Valid USD", T: TestArgs{"USD", 1, 2019}, want: 1, wantErr: false},
		{name: "Valid HKD", T: TestArgs{"HKD", 1, 2019}, want: 7.84073551516471},
		{name: "Valid GBP", T: TestArgs{"GBP", 8, 2017}, want: 0.771124120667742, wantErr: false},
		{name: "Invalid Month 18", T: TestArgs{"GBP", 18, 7}, want: 0.771124120667742, wantErr: true},
		{name: "Invalid Country Code FFF", T: TestArgs{"FFF", 8, 2017}, want: 0.771124120667742, wantErr: true},
	}

	for _, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {
			var err error
			got, err := dbs.ExchangeRate(tt.T.convertFromCurrency, tt.T.month, tt.T.year)
			if (err != nil) && tt.wantErr {
				return
			}

			if (err != nil) != tt.wantErr {
				t.Logf("Look at me, I wanted an error")
				t.Errorf("Exchange Rate Error:  got  = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_generateSqlQuery(t *testing.T) {

	type args struct {
		dt         DeviceType
		ds         DeviceStatus
		deviceList *string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "VMOnly -Exclude Offline -No Device List",
			args: args{
				dt:         VMOnly,
				ds:         ExcludeOffline,
				deviceList: nil,
			},
			want: `select so.status,
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
		Where (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%')) and so.status_number != -1 order by s.computer_number`,
		},

		{
			name: "VMOnly -All Statuses -No Device List",
			args: args{
				dt:         VMOnly,
				ds:         AllStatuses,
				deviceList: nil,
			},
			want: `select so.status,
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
 		Where (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%'))
         order by s.computer_number`,
		},

		{
			name: "VMOnly -Exclude Offline -Has Device List",
			args: args{
				dt:         VMOnly,
				ds:         ExcludeOffline,
				deviceList: &[]string{"1234,5678,9011"}[0],
			},
			want: `select so.status,
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
 		Where s.computer_number in (?)
         and (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%'))
         and so.status_number != -1
         order by s.computer_number`,
		},
		{
			name: "VMOnly -All Statuses -Has device list",
			args: args{
				dt:         VMOnly,
				ds:         AllStatuses,
				deviceList: &[]string{"1234,5678,9011"}[0],
			},
			want: `select so.status,
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
 		Where s.computer_number in (?)
         and (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%'))
         order by s.computer_number`,
		},
		{
			name: "NonVMOnly -Exclude Offline -No Device List",
			args: args{
				dt:         NonVMOnly,
				ds:         ExcludeOffline,
				deviceList: nil,
			},
			want: `select so.status,
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
 		Where (sku.product_description not like ('%Virtual%')and sku.product_description not like ('%VM%'))
         and so.status_number != -1
         order by s.computer_number`,
		},
		{
			name: "NonVMOnly -Exclude Offline -Has Device List",
			args: args{
				dt:         NonVMOnly,
				ds:         ExcludeOffline,
				deviceList: &[]string{"1234,5678,9011"}[0],
			},
			want: `select so.status,
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
 		Where s.computer_number in (?)
         and (sku.product_description not like ('%Virtual%')and sku.product_description not like ('%VM%'))
         and so.status_number != -1
         order by s.computer_number`,
		},

		{
			name: "NonVMOnly -All Statuses -Has device list",
			args: args{
				dt:         NonVMOnly,
				ds:         AllStatuses,
				deviceList: &[]string{"1234,5678,9011"}[0],
			},
			want: `select so.status,
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
 		Where s.computer_number in (?)
         and (sku.product_description not like ('%Virtual%')and sku.product_description notlike ('%VM%'))
         order by s.computer_number`,
		},
		{
			name: "NonVMOnly -All Statuses -No device list",
			args: args{
				dt:         NonVMOnly,
				ds:         AllStatuses,
				deviceList: nil,
			},
			want: `select so.status,
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
 		Where  (sku.product_description not like ('%Virtual%')and sku.product_description notlike ('%VM%'))
         order by s.computer_number`,
		},
		//
		{
			name: "VMandNonVM -All Statuses -No device list",
			args: args{
				dt:         VMandNonVM,
				ds:         AllStatuses,
				deviceList: nil,
			},
			want: `select so.status,
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
         order by s.computer_number`,
		},
		{
			name: "VMandNonVM -All Statuses -Has device list",
			args: args{
				dt:         VMandNonVM,
				ds:         AllStatuses,
				deviceList: &[]string{"1234,5678,9011"}[0],
			},
			want: `select so.status,
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
 		Where s.computer_number in (?)
         order by s.computer_number`,
		},
		//
		{
			name: "VMandNonVM -Exclude Offline -No device list",
			args: args{
				dt:         VMandNonVM,
				ds:         ExcludeOffline,
				deviceList: nil,
			},
			want: `select so.status,
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
 		Where so.status_number != -1 
         order by s.computer_number`,
		},
		{
			name: "VMandNonVM -Exclude Offline -Has device list",
			args: args{
				dt:         VMandNonVM,
				ds:         ExcludeOffline,
				deviceList: &[]string{"1234,5678,9011"}[0],
			},
			want: `select so.status,
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
 		Where s.computer_number in (?)
         and so.status_number != -1
         order by s.computer_number`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateSqlQuery(tt.args.dt, tt.args.ds, tt.args.deviceList)

			if got == "" {
				t.Errorf("No query returend %v, want %v", got, tt.want)
			}
			if stripAllWhitespace(got) != stripAllWhitespace(tt.want) {
				t.Errorf("Queries not equal returend \n  %v, \n   want %v", got, tt.want)

			}
		})
	}

}

func stripAllWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}

	}
	return strings.ToLower(b.String())
}
func setupInvalidODSDB() *sqlx.DB {

	odsUsername := os.Getenv("BadID")
	odsPassword := os.Getenv("BadPass")
	odsHost := os.Getenv("BadHost")
	odsPortS := os.Getenv("BadPort")
	os.Getenv("BaddPassTwo")
	odsPort, _ := strconv.Atoi(odsPortS)

	dbConfig := ODSConfig{
		Host:     odsHost,
		Port:     odsPort,
		User:     odsUsername,
		Password: odsPassword,
		DBName:   "BadDB",
	}

	connStr := fmt.Sprintf("odbc:server=%s; port=%d; user id=%s;password=%s; database=%s;log=3;encrypt=false;TrustServerCertificate=true", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	db, _ := sqlx.Open("mssql", connStr)
	return db

}

func setupValidODSDB() *sqlx.DB {

	odsUsername := os.Getenv("ODSID")
	odsPassword := os.Getenv("ODSPW")
	odsHost := os.Getenv("ODSHOST")
	odsPortS := os.Getenv("ODSPORT")
	os.Getenv("CoreProdPassword")
	odsPort, _ := strconv.Atoi(odsPortS)

	dbConfig := ODSConfig{
		Host:     odsHost,
		Port:     odsPort,
		User:     odsUsername,
		Password: odsPassword,
		DBName:   "Operational_reporting_CORE",
	}

	connStr := fmt.Sprintf("odbc:server=%s; port=%d; user id=%s;password=%s; database=%s;log=3;encrypt=false;TrustServerCertificate=true", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	db, _ := sqlx.Open("mssql", connStr)
	return db

}

func Test_odsDB_Ping(t *testing.T) {

	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Invalid DB Test",
			fields: fields{
				db: setupInvalidODSDB(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Valid DB Test",
			fields: fields{
				db: setupValidODSDB(),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ods := &odsDB{
				db: tt.fields.db,
			}
			got, err := ods.Ping()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_castStringToFloat64(t *testing.T) {
	type args struct {
		strNumber string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Float",
			args: args{
				strNumber: "44.4",
			},
			want: 44.4,
		},
		{
			name: "Integer",
			args: args{
				strNumber: "44",
			},
			want: 44,
		},
		{
			name: "Letters",
			args: args{
				strNumber: "ABC",
			},
			want: 0.0,
		},
		{
			name: "Nil",
			args: args{
				strNumber: "",
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := castStringToFloat64(tt.args.strNumber); got != tt.want {
				t.Errorf("castStringToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestODSConfig_ConnectionInfo(t *testing.T) {
	type fields struct {
		Host             string
		Port             int
		User             string
		Password         string
		CoreProdPassword string
		DBName           string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Valid Connection Details",
			fields: fields{
				Host:             "myhost",
				Port:             1111,
				User:             "myuser",
				Password:         "myupass",
				CoreProdPassword: "mycoreprodpass",
				DBName:           "mydbname",
			},
			want: "myuser:myupass@tcp(myhost:1111)/mydbname?parseTime=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ODSConfig{
				Host:             tt.fields.Host,
				Port:             tt.fields.Port,
				User:             tt.fields.User,
				Password:         tt.fields.Password,
				CoreProdPassword: tt.fields.CoreProdPassword,
				DBName:           tt.fields.DBName,
			}
			if got := c.ConnectionInfo(); got != tt.want {
				t.Errorf("ConnectionInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewODSService(t *testing.T) {

	validPort, _ := strconv.Atoi(os.Getenv("ODSPORT"))
	type args struct {
		config ODSConfig
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{{
		name: "Valid ODS Service Creation",
		args: args{
			config: ODSConfig{
				User:     os.Getenv("ODSID"),
				Password: os.Getenv("ODSPW"),
				Host:     os.Getenv("ODSHOST"),
				Port:     validPort,
				DBName:   os.Getenv("CoreProdPassword"),
			},
		},
		want:    true,
		wantErr: false,
	},
		{
			name: "Invalid ODS Service Creation",
			args: args{
				config: ODSConfig{
					Host:     "myhost",
					Port:     1111,
					User:     "myuser",
					Password: "mypass",
					DBName:   "mydb",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Nil  ODS Service Creation",
			args: args{
				config: ODSConfig{
					Host:     "",
					Port:     1111,
					User:     "",
					Password: "",
					DBName:   "",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewODSService(tt.args.config)
			if (err != nil) && !tt.wantErr && got == nil {
				t.Errorf("NewODSService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				resp, _ := got.Ping()

				if !reflect.DeepEqual(resp, tt.want) {
					t.Errorf("NewODSService() got = %v, want %v", resp, tt.want)
				}
			}

		})
	}
}

func TestDeviceType_String(t *testing.T) {
	tests := []struct {
		name string
		dt   DeviceType
		want string
	}{
		{
			name: "VMOnly",
			dt:   VMOnly,
			want: "VMOnly",
		},
		{
			name: "NotVMOnly",
			dt:   NonVMOnly,
			want: "NotVMOnly",
		}, {
			name: "VMandNonVM",
			dt:   VMandNonVM,
			want: "VMandNonVM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceStatus_String(t *testing.T) {
	tests := []struct {
		name string
		ds   DeviceStatus
		want string
	}{
		{
			name: "ExcludeOffline",
			ds:   ExcludeOffline,
			want: "ExcludeOffline",
		},
		{
			name: "AllStatuses",
			ds:   AllStatuses,
			want: "AllStatuses",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ds.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Mocker struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func Test_odsDB_DeviceDetails(t *testing.T) {

	type fields struct {
		mocker         Mocker
		columnHeadings []string
		expectedQuery  string
		expectedRows   []string
	}
	type args struct {
		dt         DeviceType
		ds         DeviceStatus
		deviceList string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]DeviceDetail
		wantErr bool
	}{
		{
			name: "VMandNonVM Exclude Offline All Devices",
			fields: fields{
				mocker:         getNewMockDB(),
				columnHeadings: []string{"status", "product_description", "datacenter_abbr", "computer_number", "customer_number", "server_name", "ip_address"},
				expectedQuery:  "select so.status, sku.product_description, dc.datacenter_abbr, s.computer_number, s.customer_number, s.server_name, ip.IPAddress as 'ip_address' from CORE_Prod.server as s inner join CORE_Prod.datacenter as dc on s.datacenter_number = dc.datacenter_number inner join CORE_Prod.status_options as so on s.status_number = so.status_number inner join CORE_Prod.sku as sku on s.product_id = sku.product_sku inner join CORE_Prod.IPSP_cache_IPAssignment as ip on s.computer_number = ip.computer_number Where so.status_number != -1 order by s.computer_number",
				expectedRows: []string{
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,74.205.2.141",
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,10.17.106.141",
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,10.132.149.153",
				},
			},
			args: args{
				dt:         VMandNonVM,
				ds:         ExcludeOffline,
				deviceList: "",
			},
			want: &[]DeviceDetail{
				{"Online/Complete",
					"Virtual Machine Linux Required",
					"IAD3",
					"681189",
					"1583737",
					gimmeAPointer("681190-prod.segmentsupport.rackspace.com"),
					gimmeAPointer("10.132.149.153, 10.17.106.141, 74.205.2.141"),
				}},
			wantErr: false,
		},
		{
			name: "VMandNonVM Exclude Offline All Devices - no returned ips",
			fields: fields{
				mocker:         getNewMockDB(),
				columnHeadings: []string{"status", "product_description", "datacenter_abbr", "computer_number", "customer_number", "server_name", "ip_address"},
				expectedQuery:  "select so.status, sku.product_description, dc.datacenter_abbr, s.computer_number, s.customer_number, s.server_name, ip.IPAddress as 'ip_address' from CORE_Prod.server as s inner join CORE_Prod.datacenter as dc on s.datacenter_number = dc.datacenter_number inner join CORE_Prod.status_options as so on s.status_number = so.status_number inner join CORE_Prod.sku as sku on s.product_id = sku.product_sku inner join CORE_Prod.IPSP_cache_IPAssignment as ip on s.computer_number = ip.computer_number Where so.status_number != -1 order by s.computer_number",
				expectedRows: []string{
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,",
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,",
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,",
				},
			},
			args: args{
				dt:         VMandNonVM,
				ds:         ExcludeOffline,
				deviceList: "",
			},
			want: &[]DeviceDetail{
				{"Online/Complete",
					"Virtual Machine Linux Required",
					"IAD3",
					"681189",
					"1583737",
					gimmeAPointer("681190-prod.segmentsupport.rackspace.com"),
					gimmeAPointer(""),
				}},
			wantErr: false,
		},

		{
			name: "VMOnly Exclude Offline All Devices",
			fields: fields{
				mocker:         getNewMockDB(),
				columnHeadings: []string{"status", "product_description", "datacenter_abbr", "computer_number", "customer_number", "server_name", "ip_address"},
				//expectedQuery: "select so.status, sku.product_description, dc.datacenter_abbr, s.computer_number, s.customer_number, s.server_name, ip.IPAddress as 'ip_address' from CORE_Prod.server as s inner join CORE_Prod.datacenter as dc on s.datacenter_number = dc.datacenter_number inner join CORE_Prod.status_options as so on s.status_number = so.status_number inner join CORE_Prod.sku as sku on s.product_id = sku.product_sku inner join CORE_Prod.IPSP_cache_IPAssignment as ip on s.computer_number = ip.computer_number Where (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%')) And so.status_number != -1 order by s.computer_number",
				expectedQuery: "select (.+) from CORE(.+) Where \\(sku.product_description like \\(\\'%Virtual%\\'\\) OR sku.product_description like \\(\\'%VM%\\'\\)\\) And so.status_number != -1 order by s.computer_number$",
				expectedRows: []string{
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,",
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,",
					"Online/Complete,Virtual Machine Linux Required,IAD3,681189,1583737,681190-prod.segmentsupport.rackspace.com,",
				},
			},
			args: args{
				dt:         VMOnly,
				ds:         ExcludeOffline,
				deviceList: "",
			},

			want: &[]DeviceDetail{
				{"Online/Complete",
					"Virtual Machine Linux Required",
					"IAD3",
					"681189",
					"1583737",
					gimmeAPointer("681190-prod.segmentsupport.rackspace.com"),
					gimmeAPointer(""),
				}},
			wantErr: false,
		},

		{
			name: "No Devices Found - Return Error",
			fields: fields{
				mocker:         getNewMockDB(),
				columnHeadings: []string{"status", "product_description", "datacenter_abbr", "computer_number", "customer_number", "server_name", "ip_address"},
				//expectedQuery: "select so.status, sku.product_description, dc.datacenter_abbr, s.computer_number, s.customer_number, s.server_name, ip.IPAddress as 'ip_address' from CORE_Prod.server as s inner join CORE_Prod.datacenter as dc on s.datacenter_number = dc.datacenter_number inner join CORE_Prod.status_options as so on s.status_number = so.status_number inner join CORE_Prod.sku as sku on s.product_id = sku.product_sku inner join CORE_Prod.IPSP_cache_IPAssignment as ip on s.computer_number = ip.computer_number Where (sku.product_description like ('%Virtual%') OR sku.product_description like ('%VM%')) And so.status_number != -1 order by s.computer_number",
				expectedQuery: "select (.+) from CORE(.+) Where \\(sku.product_description like \\(\\'%Virtual%\\'\\) OR sku.product_description like \\(\\'%VM%\\'\\)\\) And so.status_number != -1 order by s.computer_number$",
				expectedRows:  []string{},
			},
			args: args{
				dt:         VMOnly,
				ds:         ExcludeOffline,
				deviceList: "",
			},

			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.fields.mocker.db.Close()
			sqlxDB := sqlx.NewDb(tt.fields.mocker.db, "sqlmock")
			odsDB := &odsDB{sqlxDB}
			ODSDB := odsDB
			odssrv := odsService{ODSDB}

			rows := sqlmock.NewRows(tt.fields.columnHeadings)
			for _, r := range tt.fields.expectedRows {
				rows.FromCSVString(r)
			}

			tt.fields.mocker.mock.ExpectQuery(tt.fields.expectedQuery).WillReturnRows(rows)
			got, err := odssrv.DeviceDetails(tt.args.dt, tt.args.ds, &tt.args.deviceList)
			if (err != nil) != tt.wantErr {

				qs := strings.Split(err.Error(), " with expected regexp")

				var qrs []string
				for _, q := range qs {
					t.Error(" ")
					q = strings.Trim(q, "Query: could not match actual sql: ")
					qrs = append(qrs, q)
				}
				t.Errorf("queries equal?: %v", qrs[0] == qrs[1])
				t.Errorf("queries deep equal?: %v", reflect.DeepEqual(qrs[0], qrs[1]))

				t.Errorf(err.Error())
				// t.Errorf("DeviceDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				if tt.want != nil {
					t.Errorf("Got nil, want %v", tt.want)
					return
				} else {
					return
				}
			}

			t.Log("want: ", tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeviceDetails() got = %v, want %v", got, tt.want)
			}
			sqlxDB = nil
			odsDB = nil
			rows = nil
		})
	}
}

func Test_odsDB_DeviceDetails_Split_Queries(t *testing.T) {

	odsService, _ := setupODSTests()

	_, err := odsService.Ping()
	if err != nil {
		t.Logf("%v\n%v\n", err, errors.Unwrap(err))
	}
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		dt         DeviceType
		ds         DeviceStatus
		deviceList *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]DeviceDetail
		wantErr bool
	}{
		{
			name:   "give it to me all now quickly",
			fields: fields{db: setupValidODSDB()},
			args: args{
				dt:         VMandNonVM,
				ds:         ExcludeOffline,
				deviceList: nil,
			},
			want:    &[]DeviceDetail{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ods := &odsDB{
				db: tt.fields.db,
			}
			got, err := ods.GetAllActiveDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("DeviceDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("DeviceDetails()  want %v", tt.want)
			}

			for _, i := range *got {
				t.Logf("device: %s:   %s", i.ComputerNumber, *i.IPAddress)

			}
		})
	}
}
