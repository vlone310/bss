package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func RequireBodyMatch[T any](t *testing.T, body *bytes.Buffer, reqBody T) {
	t.Helper()
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotBody T
	err = json.Unmarshal(data, &gotBody)
	require.NoError(t, err)

	require.Equal(t, reqBody, gotBody)
}
