package service

type storeReturnValues struct {
	upsert bool
	err    error
}
type mockStore struct {
	gotPorts     []Port
	returnValues []storeReturnValues
}

// Please be careful with this quick and dirty mock. Can panic in badly written tests
func (m *mockStore) Upsert(port Port) (bool, error) {
	m.gotPorts = append(m.gotPorts, port)
	returnValue := m.returnValues[0]
	m.returnValues = m.returnValues[1:]

	return returnValue.upsert, returnValue.err
}

type recordProviderReturnValues struct {
	port Port
	err  error
}
type mockRecordProvider struct {
	callCount    int
	returnValues []recordProviderReturnValues
}

// Please be careful with this quick and dirty mock. Can panic in badly written tests
func (m *mockRecordProvider) Next() (Port, error) {
	m.callCount++

	returnValue := m.returnValues[0]
	m.returnValues = m.returnValues[1:]

	return returnValue.port, returnValue.err
}
