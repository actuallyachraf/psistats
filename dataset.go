package psistats

// Dataset is a structured key,value interface.
type Dataset interface {
	Put(key []byte, val []byte)
	Get(key []byte)
}
