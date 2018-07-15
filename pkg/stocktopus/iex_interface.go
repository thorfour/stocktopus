//+build !ALPHA

package stocktopus

import "github.com/thorfour/stocktopus/pkg/stock"

var stockInterface stock.Lookup

func init() {
	stockInterface = new(stock.IexWrapper)
}
