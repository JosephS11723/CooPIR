package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/JosephS11723/CooPIR/src/api/config"
)

var PublicKey *rsa.PublicKey
var PrivateKey *rsa.PrivateKey

// VerifyKey checks if the RSA key files exist and if the private key matches the public key. If they do not exist, it creates them.
func VerifyKeys() {
	// check if key file exists
	if _, err := os.Stat(config.RSAKeyDirectory + string(os.PathSeparator) + config.RSAKeyFile); os.IsNotExist(err) {
		log.Println("No private key file file found. Generating new key...")
		// generate key
		if err := GenerateKey(); err != nil {
			log.Fatal(err)
		}
	}

	// check if public key file exists
	if _, err := os.Stat(config.RSAKeyDirectory + string(os.PathSeparator) + config.RSAPublicKeyFile); os.IsNotExist(err) {
		log.Println("No public key file file found. Generating new key...")
		// generate key
		if err := GenerateKey(); err != nil {
			log.Fatal(err)
		}
	}

	// check if public key matches private key
	// read private key from file
	privateKey, err := GetKey()
	if err != nil {
		log.Fatal(err)
	}

	// read public key from file
	publicKey, err := GetPublicKey()
	if err != nil {
		log.Fatal(err)
	}

	// check if public key matches private key
	if publicKey.N.Cmp(privateKey.PublicKey.N) != 0 {
		log.Println("Public key does not match private key. Regenerating key.")
		// generate key
		if err := GenerateKey(); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Keys verified.")

	LoadKeys()
}

// LoadKeys loads the public and private keys from the data folder
func LoadKeys() {
	// read private key from file
	privateKey, err := GetKey()
	if err != nil {
		log.Fatal(err)
	}

	// read public key from file
	publicKey, err := GetPublicKey()
	if err != nil {
		log.Fatal(err)
	}

	// set public key
	PublicKey = publicKey

	// set private key
	PrivateKey = privateKey
}

// GenerateKey generates a new RSA key and certificate.
func GenerateKey() error {
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, config.RSAKeySize)
	if err != nil {
		return err
	}

	// write private key to file
	privateKeyFile, err := os.OpenFile(config.RSAKeyDirectory+string(os.PathSeparator)+config.RSAKeyFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	// write private key to file
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// write private key to file
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// write public key to file
	publicKeyFile, err := os.OpenFile(config.RSAKeyDirectory+string(os.PathSeparator)+config.RSAPublicKeyFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	// write public key to file
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	// write public key to file
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return err
	}

	log.Println("New RSA keys generated.")

	return nil
}

// GetKey reads the private key from the private.pem file located in the data folder and returns the private key
func GetKey() (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(config.RSAKeyDirectory + string(os.PathSeparator) + config.RSAKeyFile)
	if err != nil {
		return nil, err
	}

	// read the contents of the file into memory
	keyData := make([]byte, config.RSAKeySize)
	if _, err := privateKeyFile.Read(keyData); err != nil {
		return nil, err
	}

	// read private key from file
	privateKeyPEM, _ := pem.Decode(keyData)

	// parse private key from PEM block
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// GetPublicKey reads the public key from the public.pem file located in the data folder and returns the public key
func GetPublicKey() (*rsa.PublicKey, error) {
	publicKeyFile, err := os.Open(config.RSAKeyDirectory + string(os.PathSeparator) + config.RSAPublicKeyFile)
	if err != nil {
		return nil, err
	}

	// read the contents of the file into memory
	keyData := make([]byte, config.RSAKeySize)
	if _, err := publicKeyFile.Read(keyData); err != nil {
		return nil, err
	}

	// read public key from file
	publicKeyPEM, _ := pem.Decode(keyData)

	// parse public key from PEM block
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyPEM.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
