package app

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key := strings.ReplaceAll(uuid.NewString(), "-", "")
	t.Log(key, len(key))

	data := []byte("hello world")
	crypted, err := encrypt(key, data)
	require.NoError(t, err)
	require.NotEmpty(t, crypted)

	decrypted, err := decrypt(key, crypted)
	require.NoError(t, err)
	require.Equal(t, data, decrypted)
}
