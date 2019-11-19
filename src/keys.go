// Credits: https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
	"golang.org/x/crypto/ssh"
)

type key struct {
	ID                 string
	Type               string `json:"type"`
	PublicKey          string
	PrivateKeyFilename string
}

func loadKeys() ([]key, error) {
	keystore := os.Getenv("KEYSTORE")
	if keystore == "" {
		return []key{}, errors.New("KEYSTORE environment variable must be defined on the server")
	}
	k, _ := url.ParseRequestURI(keystore)
	switch k.Scheme {
	case "mongo":
		return loadKeysFromMongo(keystore)
	default:
		return loadKeysFromFile(k.Path)
	}
}

func saveKeys(keys []key) error {
	keystore := os.Getenv("KEYSTORE")
	if keystore == "" {
		return errors.New("KEYSTORE environment variable must be defined on the server")
	}
	k, _ := url.ParseRequestURI(keystore)
	switch k.Scheme {
	case "mongo":
		return saveKeysToMongo(keystore, keys)
	default:
		return saveKeysToFile(k.Path, keys)
	}
}

func deleteAllKeys() error {
	keystore := os.Getenv("KEYSTORE")
	if keystore == "" {
		return errors.New("KEYSTORE environment variable must be defined on the server")
	}
	k, _ := url.ParseRequestURI(keystore)
	switch k.Scheme {
	case "mongo":
		return deleteAllKeysFromMongo(keystore)
	default:
		return deleteAllKeysFromFile(k.Path)
	}
}

func deleteKey(id string) error {
	keystore := os.Getenv("KEYSTORE")
	if keystore == "" {
		return errors.New("KEYSTORE environment variable must be defined on the server")
	}
	k, _ := url.ParseRequestURI(keystore)
	switch k.Scheme {
	case "mongo":
		return deleteKeyFromMongo(keystore, id)
	default:
		return deleteKeyFromFile(k.Path, id)
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
	newkey.ID = xid.New().String()
	newkey.PrivateKeyFilename = generateKeyFilename(newkey)

	err = addPublicKey(privateKey.PublicKey, &newkey)
	if err != nil {
		return newkey, err
	}
	err = storeKey(privateKeyBytes, newkey)
	if err != nil {
		return newkey, err
	}
	return newkey, err
}

func generateKeyFilename(newkey key) string {
	keystore := os.Getenv("KEYSTORE")
	k, err := url.ParseRequestURI(keystore)
	log.Println(k.Path)
	if err != nil || k.Scheme == "mongo" {
		k.Scheme = "file"
		k.Path = "/var/local/apirate/"
	} else if k.Path == "" {
		k.Path = "."
	} else {
		k.Path = filepath.Dir(k.Path)
	}
	k.Path = filepath.Join(k.Path, "keys")
	os.MkdirAll(k.Path, 0700)
	return filepath.Join(k.Path, newkey.ID+"_id_rsa")
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

func addPublicKey(publicKey rsa.PublicKey, newkey *key) error {
	publicKeyBytes, err := generatePublicKey(publicKey)
	if err != nil {
		return err
	}
	publicKeyString := string([]byte(publicKeyBytes))
	publicKeyString = strings.TrimSuffix(publicKeyString, "\n")
	newkey.PublicKey = publicKeyString + fmt.Sprintf(" %s@apirate", newkey.ID)
	return nil
}

func generatePublicKey(publicKey rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		log.Println(err.Error())
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

func getAllPrivateKeyFilenames() ([]string, error) {
	keyfilenames := []string{}
	keys, err := loadKeys()
	if err != nil {
		return keyfilenames, nil
	}
	for _, k := range keys {
		keyfilenames = append(keyfilenames, k.PrivateKeyFilename)
	}
	return keyfilenames, nil
}
