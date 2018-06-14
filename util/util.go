package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"os"
)

const (
	RsaDirectory = "/etc/rancher/cube"
	RsaBitSize   = 4096
)

type pkcs1PublicKey struct {
	N *big.Int
	E int
}

func GenerateRSA256() error {
	// make sure the rsa directory is exist
	if _, err := os.Stat(RsaDirectory); err != nil {
		err = os.MkdirAll(RsaDirectory, os.ModeDir|0700)
		if err != nil {
			return err
		}
	}

	// if no rsa private/public key file re-generate it
	if !CheckRSAKeyFileExist() {
		privateKey, err := GeneratePrivateKey()
		if err != nil {
			logrus.Errorf("generate private key error: %v", err)
			return err
		}

		publicKeyBytes, err := GeneratePublicKey(privateKey)
		if err != nil {
			logrus.Errorf("generate public key bytes error: %v", err)
			return err
		}

		privateKeyBytes := PrivateKeyToPEM(privateKey)

		err = WriteKeyToFile(privateKeyBytes, RsaDirectory+"/id_rsa")
		if err != nil {
			logrus.Errorf("write private key file error: %v", err)
			return err
		}

		err = WriteKeyToFile([]byte(publicKeyBytes), RsaDirectory+"/id_rsa.pub")
		if err != nil {
			logrus.Errorf("write public key file error: %v", err)
			return err
		}
	}

	return nil
}

func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RsaBitSize)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func PrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateDER,
	}

	privatePEM := pem.EncodeToMemory(&privateBlock)

	return privatePEM
}

func GeneratePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicKey := privateKey.PublicKey
	publicDER := MarshalPKCS1PublicKey(&publicKey)

	publicBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDER,
	}

	publicPEM := pem.EncodeToMemory(&publicBlock)

	return publicPEM, nil
}

func CheckRSAKeyFileExist() bool {
	if _, err := os.Stat(RsaDirectory + "/id_rsa"); err == nil {
		if _, err = os.Stat(RsaDirectory + "/id_rsa.pub"); err == nil {
			return true
		}

		err = os.Remove(RsaDirectory + "/id_rsa.pub")
		if err != nil {
			logrus.Errorf("remove id_rsa.pub file error: %v", err)
		}

		err = os.Remove(RsaDirectory + "/id_rsa")
		if err != nil {
			logrus.Errorf("remove id_rsa file error: %v", err)
		}
	}

	return false
}

func WriteKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

// MarshalPKCS1PublicKey converts an RSA public key to PKCS#1, ASN.1 DER form.
func MarshalPKCS1PublicKey(key *rsa.PublicKey) []byte {
	derBytes, _ := asn1.Marshal(pkcs1PublicKey{
		N: key.N,
		E: key.E,
	})
	return derBytes
}
