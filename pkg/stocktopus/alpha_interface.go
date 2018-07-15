//+build ALPHA

package stocktopus

import (
	"os"

	"github.com/thorfour/stocktopus/pkg/stock"
)

var stockInterface stock.Lookup

func init() {
	alphaAPI := os.Getenv("ALPHA_API")
	stockInterface = &stock.AlphaWrapper{
		APIKey: alphaAPI,
	}
}
