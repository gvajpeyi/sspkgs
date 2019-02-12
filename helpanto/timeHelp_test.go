package helpanto

import (
	"testing"
	"time"
)

func TestDateStringToUnixMicro(t *testing.T) {

	dateString := "2018-06-01T00:00:00.000Z"
	var expectedUnixTime int64 = 1527811200000
	unixMicroTime, err := DateStringToUnixMicro(dateString)

	if err != nil {
		t.Fatalf("Failed to convert to unix time: %s", err)
	}

	if expectedUnixTime != unixMicroTime {
		t.Failed()
	}

}

func TestInvalidDateStringToUnixMicro(t *testing.T) {

	dateString := "2018-06-"
	_, err := DateStringToUnixMicro(dateString)

	if err == nil {
		t.Fatalf("Failed to convert to unix time: %s", err)
	}

}

func TestUnixMicroToDateString(t *testing.T) {
	expectedDateString := "2018-06-01T00:00:00.000Z"
	var unixTime int64 = 1527811200000

	dateString := UnixMicroToDateString(unixTime)

	if dateString != expectedDateString {
		t.Failed()
	}

}

func TestLastMonth(t *testing.T) {

	if !LastMonth(time.Now().UTC()) {
		t.Log("time was not last month")
	}

	currentDate := time.Now().UTC()
	lastMonth := time.Now().UTC().AddDate(0, -1, 0)
	lastQtr := time.Now().UTC().AddDate(0, -4, 0)
	lastYear := time.Now().UTC().AddDate(-1, 0, 0)
	
	
	testCases := []struct {
		arg  time.Time
		want bool
	}{
		{currentDate, false},
		{lastMonth,  true},
		{lastQtr,  false},
		{lastYear,  false},
	}

	for _, tc := range testCases {
		got := LastMonth(tc.arg)
		t.Logf("LastMonth(%q) = %v; want %v", tc.arg, got, tc.want)

		if got != tc.want {
			t.Errorf("LastMonth(%q) = %v; want %v", tc.arg, got, tc.want)
		}
	}



}


func TestLastYear(t *testing.T) {
	currentDate := time.Now().UTC()
	lastMonth := time.Now().UTC().AddDate(0, -1, 0)
	lastQtr := time.Now().UTC().AddDate(0, -4, 0)
	lastYear := time.Now().UTC().AddDate(-1, 0, 0)
	testCases := []struct {
		arg  time.Time
		want bool
	}{
				{currentDate, false},
		{lastMonth,  false},
		{lastQtr,  true},
		{lastYear,  true},

	}

	for _, tc := range testCases {
		got := LastYear(tc.arg)
		t.Logf("LastYear(%q) = %v; want %v", tc.arg, got, tc.want)

		if got != tc.want {
			t.Errorf("LastYear(%q) = %v; want %v", tc.arg, got, tc.want)
		}
	}

}

func TestSameYear(t *testing.T) {
	currentDate := time.Now().UTC()
	lastMonth := time.Now().UTC().AddDate(0, -1, 0)
	lastQtr := time.Now().UTC().AddDate(0, -4, 0)
	lastYear := time.Now().UTC().AddDate(-1, 0, 0)
	testCases := []struct {
		arg  time.Time
		want bool
	}{
				{currentDate, true},
		{lastMonth,  true},
		{lastQtr,  false},
		{lastYear,  false},


	}

	for _, tc := range testCases {
		got := SameYear(tc.arg)
		t.Logf("SameYear(%q) = %v; want %v", tc.arg, got, tc.want)

		if got != tc.want {
			t.Errorf("SameYear(%q) = %v; want %v", tc.arg, got, tc.want)
		}
	}





}

func TestSecondsInPeriod(t *testing.T) {
	
	testCases := []struct {
		arg  time.Time
		want int
	}{
		{time.Date(2019, time.Month(1), 21, 0, 0, 0, 0, time.UTC), 7732799},
		{time.Date(2019, time.Month(2), 21, 0, 0, 0, 0, time.UTC), 7732799},
		{time.Date(2019, time.Month(3), 11, 0, 0, 0, 0, time.UTC), 7732799},
		{time.Date(2018, time.Month(4), 21, 0, 0, 0, 0, time.UTC), 7819199},
		{time.Date(2018, time.Month(5), 11, 0, 0, 0, 0, time.UTC), 7819199},
		{time.Date(2019, time.Month(6), 11, 0, 0, 0, 0, time.UTC), 7819199},
		{time.Date(2019, time.Month(7), 21, 0, 0, 0, 0, time.UTC), 7905599},
		{time.Date(2019, time.Month(8), 11, 0, 0, 0, 0, time.UTC), 7905599},
		{time.Date(2019, time.Month(9), 21, 0, 0, 0, 0, time.UTC), 7905599},
		{time.Date(2019, time.Month(10), 21, 0, 0, 0, 0, time.UTC), 7905599},
		{time.Date(2019, time.Month(11), 21, 0, 0, 0, 0, time.UTC), 7905599},
		{time.Date(2019, time.Month(12), 11, 0, 0, 0, 0, time.UTC), 7905599},
	}
	
	for _, tc := range testCases {
		y, m, _ := tc.arg.Date()
		qtrStart, qtrEnd := QuarterDates(y, Quarter(m))
		curQtrSec := SecondsInPeriod(qtrStart, qtrEnd)
		
		got := curQtrSec
		t.Logf("want %v;  got %v", tc.want, got )
		
		if int64(got) != int64(tc.want) {
			t.Errorf("want %v;  got %v", tc.want, got )
		}
	}
	
}

func TestQuarter(t *testing.T) {

	testCases := []struct {
		arg  time.Time
		want int
	}{
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2019, time.Month(2), 1, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2019, time.Month(3), 1, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2019, time.Month(4), 1, 0, 0, 0, 0, time.UTC), 2},
		{time.Date(2019, time.Month(5), 1, 0, 0, 0, 0, time.UTC), 2},
		{time.Date(2019, time.Month(6), 1, 0, 0, 0, 0, time.UTC), 2},
		{time.Date(2019, time.Month(7), 1, 0, 0, 0, 0, time.UTC), 3},
		{time.Date(2019, time.Month(8), 1, 0, 0, 0, 0, time.UTC), 3},
		{time.Date(2019, time.Month(9), 1, 0, 0, 0, 0, time.UTC), 3},
		{time.Date(2019, time.Month(10), 1, 0, 0, 0, 0, time.UTC), 4},
		{time.Date(2019, time.Month(11), 1, 0, 0, 0, 0, time.UTC), 4},
		{time.Date(2019, time.Month(12), 1, 0, 0, 0, 0, time.UTC), 4},
	}


		for _, tc := range testCases{
			_,tm,_ := tc.arg.Date()
			got := Quarter(tm)
			t.Logf("Quarter(%q) = %v; want %v", tc.arg, got, tc.want)

			if got != tc.want{
				t.Errorf("Quarter(%q) = %v; want %v", tc.arg, got, tc.want)
			}
		}

}

func TestInSameQuarter(t *testing.T) {

	testCases := []struct {
		arg1  time.Time
		arg2 time.Time
		want bool
	}{
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2019, time.Month(10), 1, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(12), 31, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(3), 1, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(9), 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2019, time.Month(12), 1, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(9), 1, 0, 0, 0, 0, time.UTC), false},

	}

	for _, tc := range testCases {
		got := InSameQuarter(tc.arg1, tc.arg2)
		t.Logf("InSameQuarter(%q, %q) = %v; want %v", tc.arg1, tc.arg2, got, tc.want)

		if got != tc.want {
			t.Errorf("InSameQuarter(%q, %q) = %v; want %v", tc.arg1, tc.arg2, got, tc.want)
		}
	}

}



func TestPriorQuarter(t *testing.T) {
	testCases := []struct {
		arg1 time.Time
		arg2 time.Time
		want bool
	}{
		//q1 checking prior year q4
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.Month(10), 1, 0, 0, 0, 0, time.UTC),
			true},

		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.Month(11), 1, 0, 0, 0, 0, time.UTC),
			true},

		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.Month(12), 1, 0, 0, 0, 0, time.UTC),
			true},

			//checking same month same year
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			false},

		// checking same month prior year
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			false},

		// checking same month next year
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			false},


		// checking next month against same year
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.Month(2), 1, 0, 0, 0, 0, time.UTC),
			false},


		// checking +6 months against same year
		{time.Date(2019, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			time.Date(2019, time.Month(7), 1, 0, 0, 0, 0, time.UTC),
			false},


	}

	for _, tc := range testCases {
		got := PriorQuarter(tc.arg1, tc.arg2)

		t.Logf("PriorQuarter(%q, %q) = %v; want %v", tc.arg1, tc.arg2, got, tc.want)
		if got != tc.want {
			t.Errorf("PriorQuarter(%q, %q) = %v; want %v", tc.arg1, tc.arg2, got, tc.want)
		}
	}
}
