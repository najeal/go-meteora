package idstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryption(t *testing.T) {
	testKey := "my32digitkey12345678901234567890"
	testIV := "my16digitIvKey12"

	s, err := NewIDStore(testKey, testIV, "/tmp")

	tests := []struct {
		data string
	}{
		{
			data: "mydatatest1234",
		},
		{
			data: "mydatatest123456",
		},
		{
			data: "mydatatest1234567",
		},
		{
			data: "mydatatest123456b",
		},
		{
			data: "mydatatest12345b",
		},
	}

	for _, test := range tests {
		require.NoError(t, err)
		encrypted := s.encrypt(test.data)
		encrypted2 := s.encrypt(test.data)
		require.Equal(t, encrypted, encrypted2)
		require.NotEqual(t, test.data, encrypted)
		decrypted, err := s.decrypt(encrypted)
		require.NoError(t, err)
		require.Equal(t, test.data, decrypted)
	}
}

func TestWriteRead(t *testing.T) {
	testKey := "my32digitkey12345678901234567890"
	testIV := "my16digitIvKey12"

	require.NoError(t, os.RemoveAll("/tmp/testbot"))
	require.NoError(t, os.MkdirAll("/tmp/testbot", os.FileMode(0o750)))
	s, err := NewIDStore(testKey, testIV, "/tmp/testbot")
	require.NoError(t, err)

	t.Run("PrivateKey", func(t *testing.T) {
		_, err = s.ReadPrivateKey(-1)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)
		require.NoError(t, s.StorePrivateKey(-1, "my-private-key"))
		privateKey, err := s.ReadPrivateKey(-1)
		require.NoError(t, err)
		require.Equal(t, "my-private-key", privateKey)
	})
	t.Run("PublicKey", func(t *testing.T) {
		_, err = s.ReadPublicKey(-1)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)
		require.NoError(t, s.StorePublicKey(-1, "my-public-key"))
		privateKey, err := s.ReadPublicKey(-1)
		require.NoError(t, err)
		require.Equal(t, "my-public-key", privateKey)
	})
	t.Run("ChatID", func(t *testing.T) {
		err = s.StoreChatID("my-public-key", -1)
		require.NoError(t, err)
		res, err := s.ReadChatID("my-public-key")
		require.NoError(t, err)
		require.Equal(t, int64(-1), res)
	})
}
