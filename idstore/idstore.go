package idstore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	ErrStorageFailed   = errors.New("failed to store data")
	ErrDecodeFailed    = errors.New("failed to decode data")
	ErrDecryptFailed   = errors.New("failed to decrypt data")
	ErrReadFile        = errors.New("failed to read encoded file")
	ErrNotFound        = errors.New("data not found")
	ErrInitChatStorage = errors.New("cannot init chat storage")
)

const (
	privateKeyFileName = "privatekey.txt"
	publicKeyFileName  = "publickey.txt"

	walletPath = "wallets"
	chatPath   = "chats"
)

type IDStore struct {
	iv   string
	root string
	m    map[int64]struct{}

	cipherBlock cipher.Block
}

func NewIDStore(key, iv, root string) (*IDStore, error) {
	cipherBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, ErrStorageFailed
	}
	err = os.MkdirAll(root, os.FileMode(0o750))
	if err != nil {
		return nil, ErrInitChatStorage
	}
	return &IDStore{
		cipherBlock: cipherBlock,
		iv:          iv,
		root:        root,
		m:           map[int64]struct{}{},
	}, nil
}

func (s *IDStore) StorePrivateKey(chatID int64, privateKey string) error {
	err := os.MkdirAll(filepath.Join(s.root, walletPath, strconv.FormatInt(chatID, 10)), os.FileMode(0o750))
	if err != nil {
		return ErrInitChatStorage
	}
	encrypted := s.encrypt(privateKey)
	err = os.WriteFile(filepath.Join(s.root, walletPath, strconv.FormatInt(chatID, 10), privateKeyFileName), []byte(encrypted), os.FileMode(0o640))
	if err != nil {
		fmt.Println(err)
		return ErrStorageFailed
	}
	return nil
}

func (s *IDStore) StorePublicKey(chatID int64, publicKey string) error {
	err := os.MkdirAll(filepath.Join(s.root, walletPath, strconv.FormatInt(chatID, 10)), os.FileMode(0o750))
	if err != nil {
		return ErrInitChatStorage
	}
	encrypted := s.encrypt(publicKey)
	err = os.WriteFile(filepath.Join(s.root, walletPath, strconv.FormatInt(chatID, 10), publicKeyFileName), []byte(encrypted), os.FileMode(0o640))
	if err != nil {
		return ErrStorageFailed
	}
	return nil
}

func (s *IDStore) ReadPrivateKey(chatID int64) (string, error) {
	return s.ReadFile(filepath.Join(walletPath, strconv.FormatInt(chatID, 10)), privateKeyFileName)
}

func (s *IDStore) ReadPublicKey(chatID int64) (string, error) {
	return s.ReadFile(filepath.Join(walletPath, strconv.FormatInt(chatID, 10)), publicKeyFileName)
}

func (s *IDStore) StoreChatID(publicKey string, chatID int64) error {
	err := os.RemoveAll(filepath.Join(s.root, chatPath, publicKey))
	if err != nil {
		return ErrInitChatStorage
	}
	err = os.MkdirAll(filepath.Join(s.root, chatPath, publicKey), os.FileMode(0o750))
	if err != nil {
		return ErrInitChatStorage
	}
	encrypted := s.encrypt(strconv.FormatInt(chatID, 10))
	err = os.WriteFile(filepath.Join(s.root, chatPath, publicKey, encrypted), []byte(``), os.FileMode(0o640))
	if err != nil {
		return ErrStorageFailed
	}
	return nil
}

func (s *IDStore) ReadChatID(publicKey string) (int64, error) {
	entries, err := os.ReadDir(filepath.Join(s.root, chatPath, publicKey))
	if err != nil {
		return 0, ErrNotFound
	}
	if len(entries) == 0 {
		return 0, ErrNotFound
	}
	chatIDString, err := s.decrypt(entries[0].Name())
	if err != nil {
		return 0, ErrDecryptFailed
	}
	chatID, err := strconv.ParseInt(chatIDString, 10, 64)
	if err != nil {
		return 0, ErrDecryptFailed
	}
	return chatID, nil
}

func (s *IDStore) ReadFile(path string, fileName string) (string, error) {
	fileContent, err := os.ReadFile(filepath.Join(s.root, path, fileName))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return "", ErrNotFound
		}
		return "", ErrReadFile
	}
	decrypted, err := s.decrypt(string(fileContent))
	if err != nil {
		return "", ErrDecryptFailed
	}
	return decrypted, nil
}

func (s *IDStore) encrypt(data string) string {
	plainTextBlock := PKCS5Padding([]byte(data), aes.BlockSize)
	cipherText := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(s.cipherBlock, []byte(s.iv))
	mode.CryptBlocks(cipherText, plainTextBlock)
	str := base64.StdEncoding.EncodeToString(cipherText)
	return str
}

func (s *IDStore) decrypt(data string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", ErrDecodeFailed
	}
	mode := cipher.NewCBCDecrypter(s.cipherBlock, []byte(s.iv))
	mode.CryptBlocks(cipherText, cipherText)
	return string(PKCS5Unpading(cipherText, aes.BlockSize)), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpading(data []byte, blockSize int) []byte {
	n := int(data[len(data)-1])
	return data[:len(data)-n]
}
