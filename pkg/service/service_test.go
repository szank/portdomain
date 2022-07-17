package service

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortDomainService(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		recordProvider *mockRecordProvider
		store          *mockStore
		expectedError  string
	}{
		{
			name: "OK",
			recordProvider: &mockRecordProvider{
				returnValues: []recordProviderReturnValues{
					{
						port: Port{
							Name: "port1",
						},
						err: nil,
					},
					{
						port: Port{
							Name: "port2",
						},
						err: nil,
					},
					{
						port: Port{},
						err:  io.EOF,
					},
				},
			},
			store: &mockStore{
				returnValues: []storeReturnValues{
					{
						upsert: false,
						err:    nil,
					},
					{
						upsert: false,
						err:    nil,
					},
				},
			},
			expectedError: "",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := &PortDomainService{
				db:       tt.store,
				provider: tt.recordProvider,
			}

			err := service.Start()
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}

			// and so on, there's a lot of tests and a lot of assertions to write
			require.Equal(t, 3, tt.recordProvider.callCount)
			require.Len(t, tt.store.gotPorts, 2)
		})
	}
}
