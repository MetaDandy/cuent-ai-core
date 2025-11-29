//go:build containers

package integration

// MockFindAllOptions es un mock para FindAllOptions compartido entre paquetes
type MockFindAllOptions struct {
	Limit  int
	Offset int
}

func (m *MockFindAllOptions) GetLimit() int        { return m.Limit }
func (m *MockFindAllOptions) GetOffset() int       { return m.Offset }
func (m *MockFindAllOptions) GetSort() string      { return "" }
func (m *MockFindAllOptions) GetSortOrder() string { return "ASC" }
