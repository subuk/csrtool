package csrtool

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
)

// KeyType represents the type of key to generate
type KeyType string

const (
	RSA2048 KeyType = "RSA2048"
	RSA4096 KeyType = "RSA4096"
	EC256   KeyType = "EC256"
	EC384   KeyType = "EC384"
)

// GeneratePrivateKey generates a new private key based on the specified type
func GeneratePrivateKey(keyType KeyType) (interface{}, error) {
	switch keyType {
	case RSA2048:
		return rsa.GenerateKey(rand.Reader, 2048)
	case RSA4096:
		return rsa.GenerateKey(rand.Reader, 4096)
	case EC256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case EC384:
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	default:
		return nil, errors.New("unsupported key type")
	}
}

type pkcs10CertificationRequest struct {
	Raw                      asn1.RawContent
	CertificationRequestInfo certificationRequestInfo
	SignatureAlgorithm       pkix.AlgorithmIdentifier
	SignatureValue           asn1.BitString
}

type certificationRequestInfo struct {
	Raw           asn1.RawContent
	Version       int
	Subject       asn1.RawValue
	PublicKey     publicKeyInfo
	RawAttributes []asn1.RawValue `asn1:"tag:0"`
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

// GenerateCSR generates a Certificate Signing Request using ASN.1 directly
func GenerateCSR(privateKey interface{}, subject pkix.Name, dnsNames []string) ([]byte, error) {
	// Create the ASN.1 structure for the subject
	subjectRDN, err := asn1.Marshal(subject.ToRDNSequence())
	if err != nil {
		return nil, err
	}

	// Get the public key info
	publicKeyInfo, err := getPublicKeyInfo(privateKey)
	if err != nil {
		return nil, err
	}

	// Create the certification request info
	certReqInfo := certificationRequestInfo{
		Version: 0,
		Subject: asn1.RawValue{
			FullBytes: subjectRDN,
		},
		PublicKey:     publicKeyInfo,
		RawAttributes: []asn1.RawValue{},
	}

	// Marshal the certification request info to be signed
	certReqInfoBytes, err := asn1.Marshal(certReqInfo)
	if err != nil {
		return nil, err
	}

	// Sign the certification request info
	signature, err := signCSR(privateKey, certReqInfoBytes)
	if err != nil {
		return nil, err
	}

	// Create the final CSR structure
	csr := pkcs10CertificationRequest{
		CertificationRequestInfo: certReqInfo,
		SignatureAlgorithm:       getASN1SignatureAlgorithm(privateKey),
		SignatureValue: asn1.BitString{
			Bytes:     signature,
			BitLength: len(signature) * 8,
		},
	}

	// Marshal the CSR
	csrBytes, err := asn1.Marshal(csr)
	if err != nil {
		return nil, err
	}

	// Create PEM block
	block := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	return pem.EncodeToMemory(block), nil
}

// getPublicKeyInfo returns the ASN.1 encoded public key info
func getPublicKeyInfo(privateKey interface{}) (publicKeyInfo, error) {
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			return publicKeyInfo{}, err
		}
		var pki struct {
			Algorithm pkix.AlgorithmIdentifier
			PublicKey asn1.BitString
		}
		if _, err := asn1.Unmarshal(pubKeyBytes, &pki); err != nil {
			return publicKeyInfo{}, err
		}
		return publicKeyInfo{
			Raw:       pubKeyBytes,
			Algorithm: pki.Algorithm,
			PublicKey: pki.PublicKey,
		}, nil
	case *ecdsa.PrivateKey:
		pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			return publicKeyInfo{}, err
		}
		var pki struct {
			Algorithm pkix.AlgorithmIdentifier
			PublicKey asn1.BitString
		}
		if _, err := asn1.Unmarshal(pubKeyBytes, &pki); err != nil {
			return publicKeyInfo{}, err
		}
		return publicKeyInfo{
			Raw:       pubKeyBytes,
			Algorithm: pki.Algorithm,
			PublicKey: pki.PublicKey,
		}, nil
	default:
		return publicKeyInfo{}, errors.New("unsupported private key type")
	}
}

// getSignatureAlgorithm returns the appropriate x509.SignatureAlgorithm for the private key
func getSignatureAlgorithm(privateKey interface{}) x509.SignatureAlgorithm {
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return x509.SHA256WithRSA
	case *ecdsa.PrivateKey:
		switch key.Curve {
		case elliptic.P256():
			return x509.ECDSAWithSHA256
		case elliptic.P384():
			return x509.ECDSAWithSHA384
		default:
			return x509.ECDSAWithSHA256
		}
	default:
		return x509.UnknownSignatureAlgorithm
	}
}

// getASN1SignatureAlgorithm returns the ASN.1 algorithm identifier for the private key
func getASN1SignatureAlgorithm(privateKey interface{}) pkix.AlgorithmIdentifier {
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return pkix.AlgorithmIdentifier{
			Algorithm:  asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11}, // SHA256WithRSA
			Parameters: asn1.RawValue{Tag: 5},                              // NULL
		}
	case *ecdsa.PrivateKey:
		switch key.Curve {
		case elliptic.P256():
			return pkix.AlgorithmIdentifier{
				Algorithm: asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}, // ECDSAWithSHA256
			}
		case elliptic.P384():
			return pkix.AlgorithmIdentifier{
				Algorithm: asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 3}, // ECDSAWithSHA384
			}
		default:
			return pkix.AlgorithmIdentifier{
				Algorithm: asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}, // ECDSAWithSHA256
			}
		}
	default:
		return pkix.AlgorithmIdentifier{}
	}
}

// signCSR signs the CSR bytes with the private key
func signCSR(privateKey interface{}, csrBytes []byte) ([]byte, error) {
	// Create a hash of the CSR bytes
	hash := sha256.New()
	hash.Write(csrBytes)
	hashed := hash.Sum(nil)

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed)
	case *ecdsa.PrivateKey:
		return ecdsa.SignASN1(rand.Reader, key, hashed)
	default:
		return nil, errors.New("unsupported private key type")
	}
}
