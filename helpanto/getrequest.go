package helpanto

import (
	"encoding/csv"
	"fmt"
	"net/http"
)


func  BuildQueryString(qp map[string]string) string{

	var qs string
	c:=0
	for _,k:=range qp{

		if c == 0{
			qs = "?"
			c++
		}
		qs = fmt.Sprintf("%s%s", qs, fmt.Sprintf("%s=%s", k, qp[k]))

	}
	return qs

}

func ReadCsvFromUrl(url string, client *http.Client) ([][]string, error) {

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil

}
