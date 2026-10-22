package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	payloads := [][]byte{
		nil,
		{},
		[]byte("hello world"),
		{0, 255, 128, 1},
		bytes.Repeat([]byte{0, 127, 255}, 1024),
	}
	for _, keySize := range []int{16, 24, 32} {
		for index, data := range payloads {
			t.Run(fmt.Sprintf("AES%d/payload%d", keySize*8, index), func(t *testing.T) {
				key := strings.Repeat("k", keySize)
				original := bytes.Clone(data)
				crypted, err := encrypt(key, data)
				require.NoError(t, err)
				require.Len(t, crypted, len(data)+12+16)
				require.Equal(t, original, data, "encryption modified input")

				encryptedOriginal := bytes.Clone(crypted)
				decrypted, err := decrypt(key, crypted)
				require.NoError(t, err)
				require.True(t, bytes.Equal(data, decrypted))
				require.Equal(t, encryptedOriginal, crypted, "decryption modified input")

				// Previous releases used a separate 12-byte nonce with NewGCM.
				block, err := aes.NewCipher([]byte(key))
				require.NoError(t, err)
				legacy, err := cipher.NewGCM(block)
				require.NoError(t, err)
				decrypted, err = legacy.Open(nil, crypted[:12], crypted[12:], nil)
				require.NoError(t, err)
				require.True(t, bytes.Equal(data, decrypted), "legacy decryption failed")

				nonce := make([]byte, legacy.NonceSize())
				_, err = rand.Read(nonce)
				require.NoError(t, err)
				legacyCiphertext := legacy.Seal(nonce, nonce, data, nil)
				decrypted, err = decrypt(key, legacyCiphertext)
				require.NoError(t, err)
				require.True(t, bytes.Equal(data, decrypted), "legacy ciphertext is incompatible")
			})
		}
	}
}

// Fixed test-only keys and nonces capture the format used by previous releases.
var legacyGCMFixtures = []struct {
	key        string
	ciphertext string
}{
	{
		key:        "kkkkkkkkkkkkkkkk",
		ciphertext: "000102030405060708090a0bd2990dd056da16df7a2ba1f63ead70101b936faa43f0ac1e6619ca0c551fe9990fb46ff7e184de61c0f96a5e",
	},
	{
		key:        "kkkkkkkkkkkkkkkkkkkkkkkk",
		ciphertext: "000102030405060708090a0b9ee3f939dfb134f201b4345fb25706fe8022c47bf7ceeded43f9f10b73915c981c73c870dbc4546e3870a563",
	},
	{
		key:        "00112233445566778899aabbccddeeff",
		ciphertext: "000102030405060708090a0b9564ec325028d7a626c4aa550c0d00efacdf9177f6bb1cc73441bf488517781363605ff10ba8fcfaf4175c26",
	},
}

func TestDecryptLegacyFixtures(t *testing.T) {
	for _, fixture := range legacyGCMFixtures {
		t.Run(fmt.Sprintf("AES%d", len(fixture.key)*8), func(t *testing.T) {
			ciphertext, err := hex.DecodeString(fixture.ciphertext)
			require.NoError(t, err)
			plaintext, err := decrypt(fixture.key, ciphertext)
			require.NoError(t, err)
			require.Equal(t, `{"Title":"saved","Value":42}`, string(plaintext))
		})
	}
}

func TestDecryptLegacyPersistedState(t *testing.T) {
	ctx := Context{localStorage: newMemoryStorage()}
	require.NoError(t, ctx.LocalStorage().Set("/go-app/deviceID", "00112233-4455-6677-8899-aabbccddeeff"))
	ciphertext, err := hex.DecodeString(legacyGCMFixtures[2].ciphertext)
	require.NoError(t, err)
	const stateName = "/test/legacy-encrypted-state"
	require.NoError(t, ctx.LocalStorage().Set(stateName, storableState{EncryptedValue: ciphertext}))

	var manager stateManager
	var got struct {
		Title string
		Value int
	}
	manager.Get(ctx, stateName, &got)
	require.Equal(t, "saved", got.Title)
	require.Equal(t, 42, got.Value)
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	fixture := legacyGCMFixtures[2]
	ciphertext, err := hex.DecodeString(fixture.ciphertext)
	require.NoError(t, err)

	t.Run("truncated", func(t *testing.T) {
		for size := 0; size < len(ciphertext); size++ {
			input := bytes.Clone(ciphertext[:size])
			original := bytes.Clone(input)
			plaintext, err := decrypt(fixture.key, input)
			require.Error(t, err, "length %d", size)
			require.Nil(t, plaintext)
			if size < 12 {
				require.IsType(t, errors.Error{}, err)
				require.Equal(t, "ciphertext too short", err.(errors.Error).Message)
			} else {
				require.EqualError(t, err, "cipher: message authentication failed")
			}
			require.Equal(t, original, input)
		}
	})

	t.Run("tampered", func(t *testing.T) {
		for index := range ciphertext {
			input := bytes.Clone(ciphertext)
			input[index] ^= 1
			original := bytes.Clone(input)
			plaintext, err := decrypt(fixture.key, input)
			require.Error(t, err, "byte %d", index)
			require.Nil(t, plaintext)
			require.Equal(t, original, input)
		}
	})

	t.Run("wrong key", func(t *testing.T) {
		original := bytes.Clone(ciphertext)
		plaintext, err := decrypt(strings.Repeat("x", 32), ciphertext)
		require.Error(t, err)
		require.Nil(t, plaintext)
		require.Equal(t, original, ciphertext)
	})
}

func TestEncryptDecryptInvalidKeys(t *testing.T) {
	for _, size := range []int{0, 1, 15, 17, 23, 25, 31, 33} {
		t.Run(fmt.Sprintf("length%d", size), func(t *testing.T) {
			key := strings.Repeat("k", size)
			ciphertext, err := encrypt(key, []byte("hello"))
			require.Equal(t, aes.KeySizeError(size), err)
			require.Nil(t, ciphertext)

			// Invalid-key errors still take precedence over ciphertext errors.
			for _, input := range [][]byte{nil, make([]byte, 28)} {
				plaintext, err := decrypt(key, input)
				require.Equal(t, aes.KeySizeError(size), err)
				require.Nil(t, plaintext)
			}
		})
	}
}

type failingNonceReader struct {
	calls int
}

func (r *failingNonceReader) Read([]byte) (int, error) {
	r.calls++
	return 0, io.ErrUnexpectedEOF
}

func TestEncryptUsesInternalRandomness(t *testing.T) {
	// Do not run in parallel: this test temporarily replaces a global reader.
	original := rand.Reader
	t.Cleanup(func() { rand.Reader = original })
	failing := &failingNonceReader{}
	rand.Reader = failing

	key := strings.Repeat("k", 32)
	ciphertext, err := encrypt(key, []byte("hello"))
	require.NoError(t, err)
	require.Zero(t, failing.calls)
	plaintext, err := decrypt(key, ciphertext)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), plaintext)
}
