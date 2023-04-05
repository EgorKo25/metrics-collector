package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

type Encryptor struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewEncryptor(path, keyType string) (*Encryptor, error) {

	var encrypter Encryptor

	_ = encrypter.GetKey(path, keyType)

	return &encrypter, nil
}

func (e *Encryptor) GetKey(pathKey string, keyType string) (err error) {

	var pemString []byte

	pemString, err = os.ReadFile(pathKey)
	if err != nil {
		log.Println("no key for encrypt")
		return nil
	}

	block, _ := pem.Decode(pemString)

	switch keyType {
	case "public":
		e.publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return err
		}

		return nil
	case "private":
		e.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (e *Encryptor) Encrypt(message []byte) (encryptString []byte, err error) {

	if e.publicKey == nil {
		return message, nil
	}

	hash := sha256.New()
	label := []byte("")

	encryptString, err = rsa.EncryptOAEP(hash, rand.Reader, e.publicKey, message, label)
	if err != nil {
		return []byte(""), nil
	}

	return encryptString, nil
}

func (e *Encryptor) Decrypt(message []byte) (encryptString []byte, err error) {

	if e.privateKey == nil {
		return message, nil
	}

	hash := sha256.New()
	label := []byte("")

	encryptString, err = rsa.DecryptOAEP(hash, rand.Reader, e.privateKey, message, label)
	if err != nil {
		return []byte(""), nil
	}

	return encryptString, nil
}
