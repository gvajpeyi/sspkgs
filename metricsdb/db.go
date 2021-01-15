package metricsdb

//
// import (
// 	"fmt"
//
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mysql"
// )
//
// type MetricsDatabase struct {
// 	Port     int    `json:"port"`
// 	Host     string `json:"ip_address"`
// 	Name     string `json:"database_name"`
// 	User     string `json:"database_user"`
// 	Password string `json:"database_password"`
// 	CS       string `json:"connection_string"`
// }
//
// func (mdb *MetricsDatabase) OpenDB() (*gorm.DB, error) {
//
// 	mdb.CS = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mdb.User, mdb.Password, mdb.Host, mdb.Port, mdb.Name)
//
//
//
// 	db, err := gorm.Open("mysql", mdb.CS)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()
// 	db.LogMode(true)
// 	//db.AutoMigrate(&models.HitchIncident{}, &models.AVCount{}, &models.CycleCount{}, &models.CreditMemo{})
//
// 	return db, nil
// }
