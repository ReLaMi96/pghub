package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)

	headers = NewHeaders()
	data = []byte("Prime:   Agen\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Equal(t, "Agen", headers["prime"])
	data = []byte("Prime: loves-zig...agen\r\n\r\n")
	n, done, err = headers.Parse(data)
	data = []byte("Prime: gopheragen\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Equal(t, "Agen, loves-zig...agen, gopheragen", headers["prime"])

}
