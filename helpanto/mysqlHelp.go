package helpanto

import (
	"database/sql"
)

func GetDB(dbDSN string) (*sql.DB, error) {
	
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		
		db.Close()
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		
		db.Close()
		return nil, err
	}
	return db, nil
	
}
