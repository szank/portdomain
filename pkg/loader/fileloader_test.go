package loader

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/edsrzf/mmap-go"
	"github.com/stretchr/testify/require"
)

func TestFileLoaderNoError(t *testing.T) {
	t.Parallel()

	input := `
	{
		"AEAJM": {
		  "name": "Ajman",
		  "city": "Ajman",
		  "country": "United Arab Emirates",
		  "alias": [],
		  "regions": [],
		  "coordinates": [
			55.5136433,
			25.4052165
		  ],
		  "province": "Ajman",
		  "timezone": "Asia/Dubai",
		  "unlocs": [
			"AEAJM"
		  ],
		  "code": "52000"
		},
		"AEAUH": {
		  "name": "Abu Dhabi",
		  "coordinates": [
			54.37,
			24.47
		  ],
		  "city": "Abu Dhabi",
		  "province": "Abu Z¸aby [Abu Dhabi]",
		  "country": "United Arab Emirates",
		  "alias": [],
		  "regions": [],
		  "timezone": "Asia/Dubai",
		  "unlocs": [
			"AEAUH"
		  ],
		  "code": "52001"
		}
	}
	`

	loader := File{
		decoder: json.NewDecoder(bytes.NewReader([]byte(input))),
		mmap:    mmap.MMap(input),
	}

	port, err := loader.Next()
	require.NoError(t, err)
	require.Equal(t, "Ajman", port.Name)

	port, err = loader.Next()
	require.NoError(t, err)
	require.Equal(t, "Abu Dhabi", port.Name)

	port, err = loader.Next()
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, "", port.Name)
}

func TestFileLoaderInvalidData(t *testing.T) {
	t.Parallel()

	input := `
	{
		"AEAJM": {
		  "name": "Ajman",
		  "city": "Ajman",
		  "country": "United Arab Emirates",
		  "alias": [],
		  "regions": [],
		  "coordinates": [
			55.5136433,
			25.4052165
		  ],
		  "province": "Ajman",
		  "timezone": "Asia/Dubai",
		  "unlocs": [
	}
	`

	loader := File{
		decoder: json.NewDecoder(bytes.NewReader([]byte(input))),
		mmap:    mmap.MMap(input),
	}

	_, err := loader.Next()
	require.Error(t, err)
	require.EqualError(t, err, "invalid character '}' looking for beginning of value")
	// errors.Is does not give match. Sigh.
	require.IsType(t, &json.SyntaxError{}, err)
}
