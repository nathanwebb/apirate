package main

// Credits: https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

type key struct {
	ID                 int
	Type               string
	Name               string
	Description        string
	PublicKey          string
	PrivateKeyFilename string `json:"-"`
}

func createSSHKey(sshkey key) (key, error) {
	sshkey.Type = "ssh"
	sshkey.PrivateKeyFilename = "./id_rsa_test"
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}
	privateKeyBytes := encodePrivateKeyToPEM(privateKey)
	err = writeKeyToFile(privateKeyBytes, sshkey.PrivateKeyFilename)
	if err != nil {
		log.Fatal(err.Error())
	}
	sshkey.PublicKey = string([]byte(publicKeyBytes))
	return sshkey, err
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

func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
