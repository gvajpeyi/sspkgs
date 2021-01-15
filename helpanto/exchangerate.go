package helpanto

//
//
// 	"github.rackspace.com/SegmentSupport/raxss/pkg/coreods"
// )
//
//
// var ODS, err = coreods.NewODSService(config)
// const (
// 	USD Currency = "USD"
// 	GBP Currency = "GBP"
// 	EUR Currency = "EUR"
// 	AUD Currency = "AUD"
// 	CNY Currency = "CNY"
// 	HKD Currency = "HKD"
// 	RUB Currency = "RUB"
// 	SGD Currency = "SGD"
// )
//
// func (c Currency) ToString() string {
//
// 	if c == USD {
// 		return "USD"
// 	}
// 	if c == USD {
// 		return "GBP"
// 	}
// 	if c == USD {
// 		return "EUR"
// 	}
// 	if c == USD {
// 		return "AUD}"
// 	}
// 	if c == USD {
// 		return "CNY"
// 	}
// 	if c == USD {
// 		return "HKD"
// 	}
// 	if c == USD {
// 		return "RUB"
// 	}
// 	if c == USD {
// 		return "SGD"
// 	}
// 	return ""
// }
//
// func ConvertToUSD(ebiConnectionString string, convertFromCurrency string, month int, year int, valueToConvert float64) (*float64, error) {
//
//
// 	exchangeRate, err := queryEBI(ebiConnectionString, convertFromCurrency, month, year)
//
// 	if err != nil {
// 		return nil, err
// 	}
// 	convertedCurrency := (1 / *exchangeRate) * valueToConvert
//
// 	return &convertedCurrency, nil
// }
//
// func queryEBI(ebiConnectionString string, convertFromCurrency string, month int, year int) (*float64, error) {
//
// 	// var exchangeRate float64
// 	// var erate []byte
//
//
// 	//
// 	// rows, err := dbs.ExchangeRate(convertFromCurrency, month, year)
// 	//
// 	//
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// var rowCounter = 0
// 	// if rows == nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	for rows.Next() {
// 	// 		rowCounter++
// 	// 		err = rows.Scan(&erate)
// 	//
// 	// 		if err != nil {
// 	// 			return nil, err
// 	// 		}
// 	//
// 	// 		exchangeRate, err = strconv.ParseFloat(string(erate), 64)
// 	// 		if err != nil {
// 	// 			return nil, err
// 	// 		}
// 	//
// 	// 	}
// 	// }
//
// 	//return &exchangeRate, nil
// 	return nil, nil
//
// }
// import (
// 	_ "github.com/denisenkom/go-mssqldb"
