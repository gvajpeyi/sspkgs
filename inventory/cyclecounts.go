package cyclecounts
//
//import (
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/metricsdb"
//	"github.rackspace.com/SegmentSupport/raxss/pkg/network"
//	"strconv"
//	_ "github.com/go-sql-driver/mysql"
//
//)
//
//
//
//
//type CycleCountDB metricsdb.MetricsDatabase
//func NewCycleCountDatabase(mdb *metricsdb.MetricsDatabase) *CycleCountDB {
//	var cycleCountDB CycleCountDB
//
//	cstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mdb.DatabaseUser, mdb.DatabasePassword, mdb.IPAddress, mdb.Port, mdb.DatabaseName)
//
//	cycleCountDB.Port = mdb.Port
//	cycleCountDB.IPAddress = mdb.IPAddress
//	cycleCountDB.DatabaseName = mdb.DatabaseName
//	cycleCountDB.DatabaseUser = mdb.DatabaseUser
//	cycleCountDB.DatabasePassword = mdb.DatabasePassword
//	cycleCountDB.ConnectionString = cstring
//
//	fmt.Println(cstring)
//	return &cycleCountDB
//}
//
//
//func (cdb *CycleCountDB) GetCycleCounts(dc string)(*CycleCount, error){
// dc = network.GetRegionFromDC(dc)
//
//	//
//fmt.Println("Getting for: ", dc)
//
//var I2cnt []byte
//var Cfp []byte
//var Cf []byte
//var Srk []byte
//var Doc []byte
//var Dc []byte
//var FPY4WkAvg []byte
//var	Fnl4WkAvg []byte
//var Shrink4WkAvg []byte
//var cycleCount CycleCount
//
//
//	sqlStatement :=`
//Select c1.date_of_count, c1.datacenter, c1.items_to_count, c1.counted_first_pass, c1.counted_final, c1.shrink,
//		(
//	select sum(c2.counted_first_pass) / sum(c2.items_to_count)
//		from cycle_counts as c2
//		where datediff(c1.date_of_count, c2.date_of_count) Between 0 and 27
//		and c2.datacenter = ?
//		) as 'FPY4WeekAvg',
//
//		(
//		select (cast(sum(c2.counted_final) as decimal(12,4)) / cast(sum(c2.items_to_count) as decimal(12,4)))
//			from cycle_counts as c2
//			where datediff(c1.date_of_count, c2.date_of_count) Between 0 and 27
//			and c2.datacenter = ?
//			) as 'Final4WeekAvg'
//
//			,
//
//			(
//			select cast(sum(c2.shrink) as decimal(12,4))/cast(count(c2.shrink) as decimal(12,4))
//				from cycle_counts as c2
//				where datediff(c1.date_of_count, c2.date_of_count) Between 0 and 27
//				and c2.datacenter = ?
//				) as 'Shrink4WeekAvg'
//
//
//				From cycle_counts as c1
//				where c1.datacenter = ?
//order by c1.date_of_count desc LIMIT 1;
// `
//
//		//dbDSN, err := getDBConfig()
//		db, err := sqlx.Open("mysql", cdb.ConnectionString)
//		if err != nil {
//			fmt.Printf("Error connection: %v\n", err)
//			return nil, err
//		}
//
//		defer db.Close()
//		err = db.Ping()
//		if err != nil {
//
//			fmt.Printf("Ping error: %v\n", err)
//			return nil, err
//		}
//
//
//
//
//		fmt.Println(sqlStatement)
//		rows, err := db.Query(sqlStatement, dc, dc, dc, dc)
//
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//		var rowCounter= 0
//		if rows == nil {
//			return nil, err
//		} else {
//			for rows.Next() {
//				rowCounter++
//				err = rows.Scan(
//					&Doc,
//					&Dc,
//					&I2cnt,
//					&Cfp,
//					&Cf,
//					&Srk,
//					&FPY4WkAvg,
//					&Fnl4WkAvg,
//					&Shrink4WkAvg,
//
//				)
//				fmt.Println("dc: ", string(Dc))
//				cycleCount = CycleCount{
//					string(Dc),
//					string(Doc),
//					byte2int(I2cnt),
//					byte2int(Cfp),
//					byte2int(Cf),
//					byte2flt(Srk),
//					byte2flt(FPY4WkAvg),
//					byte2flt(Fnl4WkAvg),
//					byte2flt(Shrink4WkAvg),
//				}
//				if err != nil {
//					return nil, err
//				}
//
//
//			}
//			return &cycleCount, nil
//
//		}
//	}
//
//func byte2int(byteData []byte) (int ){
//	intData, _  := strconv.Atoi(string(byteData))
//	return intData
//}
//func byte2flt(byteData []byte) (float64){
//	fltdata, _ := strconv.ParseFloat(string(byteData), 64)
//	return fltdata
//
//
//
//}
//
//
//	type CycleCount struct{
//		DataCenter string `json:"data_center"`
//		CountDate string `json:"count_date"`
//		ItemsToCount int `json:"items_to_count"`
//		CountedFirstPass int `json:"counted_first_pass"`
//		CountedFinal int `json:"counted_final"`
//		Shrink float64 `json:"shrink"`
//		FPY4WAvg float64 `json:"fpy_4week_avg"`
//		Fnl4WAvg float64 `json:"final_4week_avg"`
//		Shrink4WkAvg float64 `json:"shrink_4week_avg"`
//
	//}//