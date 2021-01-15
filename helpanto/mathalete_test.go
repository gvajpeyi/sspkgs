package helpanto

import "fmt"

func ExamplePercentageDiff() {

	janSales := 750.00
	febSales := 590.00

	febDiff := PercentageDiff(janSales, febSales)

	// Output: 21.333333333333336
}
