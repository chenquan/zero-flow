package md

type Carrier interface {
	// Get returns the value associated with the passed key.
	Get(key string) []string

	// Set stores the key-value pair.
	Set(key string, values ...string)

	Append(k string, values ...string)

	// Keys lists the keys stored in this carrier.
	Keys() []string
}
