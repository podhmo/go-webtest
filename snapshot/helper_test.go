package snapshot

// test helpers (these identifiers are not included as actual dependencies, on used by library)

// Dummy : dummy data
type Dummy struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// GetData :
func GetData() Dummy {
	return Dummy{
		Name: "foo",
		Age:  20,
	}
}
