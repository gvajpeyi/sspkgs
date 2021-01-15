package helpanto

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// DateStringToUnixMicro takes a RFC3339 formatted string and converts it
// to a int64 unix epoch time
func DateStringToUnixMicro(dateString string) (int64, error) {

	time1, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return 0, err
	}
	return (time1.Unix() * 1000), nil

}

// UnixMicroToDateString takes a int64 unix epoch time and converts it to a
// RFC3339 formatted string
func UnixMicroToDateString(unixMicro int64) string {

	stringTime := time.Unix(unixMicro/1000, 0).Format(time.RFC3339)

	return stringTime
}

func ConvertDateForMariaDBInsert(dateString string) (string, error) {
	const createdFormat = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"

	timeString, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return "", err
	}
	timeString.Format(createdFormat)
	timeString = timeString.UTC()
	return strings.Split(timeString.String(), " +")[0], nil
}

// LastMont returns true if time t has month prior to time.now().UTC().Date()

func LastMonth(t time.Time) bool {
	lastMonthDate := time.Now().UTC().AddDate(0, -1, 0)

	return  lastMonthDate.Month() == t.Month()

}

// SameYear returns true if time t has exact same year as  time.now().UTC().Date()

func SameYear(t time.Time) bool {
	tYear, _, _ := t.Date()
	nowYear, _, _ := time.Now().UTC().Date()

	//
	return tYear-nowYear == 0

}
// LastYear returns true if time t has year prior to time.now().UTC().Date()
func LastYear(t time.Time) bool {
	tYear, _, _ := t.Date()
	nowYear, _, _ := time.Now().UTC().Date()

	return (tYear + 1) -nowYear == 0

}

// SameMonthSameYear returns true if time t has same month and year as current date (time.Now().UTC())
func SameMonthSameYear(t time.Time) bool {
	tYear, tMonth, _ := t.Date()
	nowYear, nowMonth, _ := time.Now().UTC().Date()
	return  tYear - nowYear == 0 && tMonth - nowMonth == 0

}

// Quarter returns back an int representation of the quarter year  that time t is in. (1, 2, 3, or 4)
func PreviousQuarter(q, y int) (int, int){
	prevQtr :=0
	prevQtrYear := y
	if q == 1 {
		prevQtr = 4
		prevQtrYear = y - 1
	} else if q == 2 {
		prevQtr = 1
		prevQtrYear = y
	} else if q == 3 {
		prevQtr = 2
		prevQtrYear = y
	} else if q == 4 {
		prevQtr = 3
		prevQtrYear = y
	}
	return  prevQtr, prevQtrYear
	
}

func Quarter(m time.Month) (int) {
	currQtr := int(math.Ceil((float64(int(m)) / 3)))
	return currQtr
	
}

// InSameQuarter check to see if time a and time b exist in the exact same quarter
func InSameQuarter(a, b time.Time) bool{

		_,am,_ := a.Date()
		_,bm,_ := b.Date()
		aq:= Quarter(am)
		bq := Quarter(bm)

		return a.Year() == b.Year() &&  aq == bq
}

func SecondsInMonth(y int, m time.Month) int64 {
	start := time.Date(y, m, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
	return int64(end.Sub(start).Seconds())
}

// func QtrStartEndDates(q,y int) (time.Time, time.Time){
//
// 	var startDate, endDate time.Time
// 	if q == 1 {
// 		startDate = time.Date(y, time.Month(1), 01, 0, 0, 0, 0, time.UTC)
// 		endDate = time.Date(y, time.Month(3), 31, 23, 59, 59, 0, time.UTC)
//
// 	} else if q == 2{
// 		startDate = time.Date(y, time.Month(4), 01, 0, 0, 0, 0, time.UTC)
// 		endDate = time.Date(y, time.Month(6), 30, 23, 59, 59, 0, time.UTC)
//
// }else if q == 3{
// 		startDate = time.Date(y, time.Month(7), 01, 0, 0, 0, 0, time.UTC)
// 		endDate = time.Date(y, time.Month(9), 30, 23, 59, 59, 0, time.UTC)
//
// 	}else if q == 4{
// 		startDate = time.Date(y, time.Month(10), 01, 0, 0, 0, 0, time.UTC)
// 		endDate = time.Date(y, time.Month(12), 31, 23, 59, 59, 0, time.UTC)
//
// 	}
//
// 	return startDate, endDate
// }

const SecondsInDay int = 60*60*24
func SecondsInPeriod(a, b time.Time) int64 {
	
	periodDuration := b.Sub(a).Seconds()
	
	return int64(periodDuration)

}

func DaysInMonth(y int, m time.Month) int{
	
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
	
}

func DaysInYear(y int) int {
	
	return time.Date(y+1, 0, 0, 0, 0, 0, 0, time.UTC).Day()
	
}

func QuarterDates(y int, q int) (time.Time, time.Time){
	
	var t time.Time
	var e time.Time
	
	if q == 1 {
		t = time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
		e = time.Date(y, 3, 31, 11, 59, 59, 0, time.UTC)
		
	} else if q == 2 {
		t = time.Date(y, 4, 1, 0, 0, 0, 0, time.UTC)
		e = time.Date(y, 6, 30, 11, 59, 59, 0, time.UTC)
		
	} else if q == 3 {
		t = time.Date(y, 7, 1, 0,0, 0, 0, time.UTC)
		e = time.Date(y, 9, 30, 11, 59, 59, 0, time.UTC)
		
	} else if q == 4 {
		t = time.Date(y, 10, 1, 0, 0, 0, 0, time.UTC)
		e = time.Date(y, 12, 31, 11, 59, 59, 0, time.UTC)
		
	}

	
	return t, e
}




// PriorQuarter checks to see if date b is in the quarter prior to date a
func PriorQuarter(a, b time.Time) bool {

	aY, aM, _ := a.Date()
	aQtr:= Quarter(aM)

	bY, bM,  _ := b.Date()
	bQtr:= Quarter(bM)


	if aY - bY == 1 && aQtr == 1 && bQtr == 4{
		return true
	}
	if aY == bY && aQtr - 1 == bQtr{
		return true
	}
	return false


}



func MonthFromString(m string) (time.Month, error){
				var mth time.Month
					    switch m {
				    case "Jan":
					    mth = time.Month(1)
				
				    case "Feb":
					    mth = time.Month(2)
				    case "Mar":
					    mth = time.Month(3)
				
				    case "Apr":
					    mth = time.Month(4)
				    case "May":
					    mth = time.Month(5)
				
				    case "Jun":
					    mth = time.Month(6)
				    case "Jul":
					    mth = time.Month(7)
				
				    case "Aug":
					    mth = time.Month(8)
				    case "Sep":
					    mth = time.Month(9)
				    case "Oct":
					    mth = time.Month(10)
				
				    case "Nov":
					    mth = time.Month(11)
				
				    case "Dec":
					    mth = time.Month(12)
					
					    default:
					    	return 0, fmt.Errorf("invalid Month: %s", m)
				    }
					    
					    return mth, nil
}

