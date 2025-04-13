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

// PKCS#9 attribute OIDs
var (
	oidChallengePassword = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 7}
)

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

// GenerateCSR generates a Certificate Signing Request using ASN.1 directly.
// The privateKey parameter should be either *rsa.PrivateKey or *ecdsa.PrivateKey.
// Users should generate their own private keys using the standard crypto package.
func GenerateCSR(privateKey interface{}, subject pkix.Name, dnsNames []string, challengePassword string) ([]byte, error) {
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

	// Prepare attributes
	var attributes []asn1.RawValue

	// Add challengePassword if provided
	if challengePassword != "" {
		challengePasswordRaw, err := asn1.Marshal(challengePassword)
		if err != nil {
			return nil, err
		}

		type attribute struct {
			Type   asn1.ObjectIdentifier
			Values []asn1.RawValue `asn1:"set"`
		}

		attr := attribute{
			Type: oidChallengePassword,
			Values: []asn1.RawValue{
				{
					Class:     asn1.ClassUniversal,
					Tag:       asn1.TagPrintableString,
					Bytes:     challengePasswordRaw,
					FullBytes: challengePasswordRaw,
				},
			},
		}

		attrBytes, err := asn1.Marshal(attr)
		if err != nil {
			return nil, err
		}

		attributes = append(attributes, asn1.RawValue{FullBytes: attrBytes})
	}

	// Create the certification request info
	certReqInfo := certificationRequestInfo{
		Version: 0,
		Subject: asn1.RawValue{
			FullBytes: subjectRDN,
		},
		PublicKey:     publicKeyInfo,
		RawAttributes: attributes,
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
