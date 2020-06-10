package stocktopus

// WatchList is the interface for interacting with a watch list
type WatchList interface {
	Add(tickers []string, key string) error
	Print(key string) (string, error)
	Remove(tickers []string, key string) error
	Clear(key string) error
}
