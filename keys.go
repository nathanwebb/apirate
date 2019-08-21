package main

// Credits: https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/url"

	"golang.org/x/crypto/ssh"
)

type key struct {
	ID                 int
	Type               string
	PublicKey          string
	PrivateKeyFilename string
}

func loadKeys(keystore string) ([]key, error) {
	k, _ := url.ParseRequestURI(keystore)
	switch k.Scheme {
	case "mongo":
		return loadKeysFromMongo(keystore)
	default:
		return loadKeysFromFile(k.Path)
	}
}

func saveKeys(keystore string, keys []key) error {
	k, _ := url.ParseRequestURI(keystore)
	switch k.Scheme {
	case "mongo":
		return saveKeysToMongo(keystore)
	default:
		return saveKeysToFile(k.Path)
	}
}

func createSSHKey(newkey key) (key, error) {
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}
	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	newkey.Type = "ssh"
	newkey.PrivateKeyFilename = "./id_rsa_test"
	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}
	newkey.PublicKey = string([]byte(publicKeyBytes))
	err = storeKey(privateKeyBytes, newkey)
	if err != nil {
		log.Fatal(err.Error())
	}
	return newkey, err
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	log.Println("Private Key generated")
	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)
	return privatePEM
}

func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	log.Println("Public key generated")
	return pubKeyBytes, nil
}

func storeKey(keyBytes []byte, newkey key) error {
	err := ioutil.WriteFile(newkey.PrivateKeyFilename, keyBytes, 0600)
	if err != nil {
		return err
	}
	log.Printf("Key saved to: %s", newkey.PrivateKeyFilename)
	return nil
}
