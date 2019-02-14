package coreods

import (
	_ "database/sql"
	"fmt"
	"time"
)

type DCServerCount struct {
	DcName      string `db:"dc_name"`
	ServerCount int    `db:"server_count"`
}

//error, scanning:  missing destination name device_datacenter_abbr in *coreods.DCServerCount

func (ods *odsDB) DataCenterServerCount(start, end time.Time, dc string) (DCServerCount, error) {

	sc := DCServerCount{}

	sY, sM, sD := start.Date()
	eY, eM, eD := end.Date()
	startTimeKey := fmt.Sprintf("%d%02d%02d", sY, int(sM), sD)
	endTimeKey := fmt.Sprintf("%d%02d%02d", eY, int(eM), eD)
	qry := "SELECT device_datacenter_abbr as 'dc_name',COUNT(UK_DMART.DCM.t_DC_Total_Servers_Trend.Device_Number) as server_count  FROM UK_DMART.DCM.t_DC_Total_Servers_Trend WHERE device_datacenter_abbr = ? and time_Key >= ? and time_Key <= ?    GROUP BY time_month_key ,device_datacenter_abbr  order by 1 desc;"

	rowErr := ods.db.QueryRowx(qry, dc, startTimeKey, endTimeKey).StructScan(&sc)

	if rowErr != nil {

		return DCServerCount{}, rowErr
	}

	return sc, nil

}

/*

--get distinct datacenters
SELECT distinct [device_datacenter_abbr]
FROM [UK_DMART].[DCM].[t_DC_Total_Servers_Trend]
  Where time_month_key = CONVERT(CHAR(6), GETDATE(), 112);



*/
