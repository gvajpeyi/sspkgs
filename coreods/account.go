package coreods

import (
	"fmt"
)

type DevAcct struct{
		ComputerNumber int `json:"computer_number" db:"computer_number"`
	CustomerNumber int `json:"customer_number" db:"customer_number"`
}


// LoadAllDevicesAndAccountNums returns all device numbers and associated account numbers
// Data is retruned as a slice of DevAcct structs.
// If a non integer account number or device number is returned it is ignored and the next record is processed.

	func (ods *odsDB) LoadAllDevicesAndAccountNums() (*[]DevAcct, error) {
	var devDetails  []DevAcct
	devices := DevAcct{}
		qry := `Select s.computer_number, s.customer_number from CORE_Prod.server as s `
		
		rows, err := ods.db.Queryx(qry)
			if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
		
		for rows.Next() {
			
			err := rows.StructScan(&devices)
				if err != nil {
		return nil, fmt.Errorf("row scan failed: %v", err)
	}
			devDetails = append(devDetails, devices)
			
		}


	return &devDetails, err

}

