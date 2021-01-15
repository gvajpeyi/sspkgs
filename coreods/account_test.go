package coreods

import "testing"

func TestOdsDB_LoadDevicesAndAccountNums(t *testing.T) {
	dbs, err := setupODSTests()
	if err != nil {
		t.Fatal(err)

	}
	
	
	testCases := []struct {
		Want int
	}{
		{Want: 980000},
	}

	for _, tc := range testCases {
		got, err := dbs.LoadAllDevicesAndAccountNums()
		if err != nil{
			t.Errorf("Got Error: %v Want: >= %v",  err.Error(), tc.Want)
			continue
		}
		if  len(*got) < tc.Want{
				t.Errorf("Got: %v; Want: >=  %v",  len(*got), tc.Want)
				continue
			
		}
		
		
		t.Logf("Got: %v; Want: >= %v", len(*got), tc.Want)

		// if got != tc.Want {
		// 	t.Errorf("err:%s:  Got: %v; Want:  %v", tc.T.convertFromCurrency, got, tc.Want)
		// }
	}
	
}