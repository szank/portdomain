package persistence

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/szank/portdomain/pkg/service"
)

func TestInMemoryDatabase(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		initialDataSet map[string]service.Port
		input          service.Port
		expectedError  string
		expectedUpsert bool
	}{
		// TODO: check for valid retrieval later
		{
			name:           "no upsert",
			initialDataSet: make(map[string]service.Port),
			input: service.Port{
				Name: "port1",
			},
			expectedError:  "",
			expectedUpsert: false,
		},
		{
			name: "upsert",
			initialDataSet: map[string]service.Port{
				"port1": service.Port{
					Name: "port1",
					City: "London",
				},
			},
			input: service.Port{
				Name: "port1",
			},
			expectedError:  "",
			expectedUpsert: true,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := &InMemoryDatabase{
				data: tt.initialDataSet,
			}

			upsert, err := db.Upsert(tt.input)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}

			require.Equal(t, tt.expectedUpsert, upsert)
		})
	}
}
