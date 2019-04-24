// go run main.go update.tar.gz ec_public.pem signature.sha1
// inspired by https://thanethomson.com/2018/11/30/validating-ecdsa-signatures-golang/
package main

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

// Represents the two mathematical components of an ECDSA signature once decomposed.
type ECDSASignature struct {
	R, S *big.Int
}

func main() {
	log.Printf("Usage: verify update.tar.gz ec_public.pem signature.sha1")

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	hash := h.Sum(nil)

	pubPEM, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := loadPublicKey(pubPEM)
	if err != nil {
		log.Fatal(err)
	}

	signature64, err := ioutil.ReadFile(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	// decode the signature to extract the DER-encoded byte string
	der, err := base64.StdEncoding.DecodeString(string(signature64))
	if err != nil {
		log.Fatal(err)
	}
	// unmarshal the R and S components of the ASN.1-encoded signature into our signature data structure
	sig := &ECDSASignature{}
	_, err = asn1.Unmarshal(der, sig)
	if err != nil {
		log.Fatal(err)
	}

	valid := ecdsa.Verify(publicKey, hash, sig.R, sig.S)
	if valid {
		log.Println("verification succeeds")
	} else {
		log.Println("verification failed")
	}

}

func loadPublicKey(publicKey []byte) (*ecdsa.PublicKey, error) {
	// decode the key, assuming it's in PEM format
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("Failed to decode PEM public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to parse ECDSA public key")
	}
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	}
	return nil, errors.New("Unsupported public key type")
}
