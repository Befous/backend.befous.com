package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Befous/backend.befous.com/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func ReadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file %s: %w", filename, err)
	}
	privateBlock, _ := pem.Decode(privateKeyBytes)
	if privateBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func ReadPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file %s: %w", filename, err)
	}
	publicBlock, _ := pem.Decode(publicKeyBytes)
	if publicBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not an RSA public key")
	}
	return publicKey, nil
}

func ReadPrivateKeyFromEnv(private string) (*rsa.PrivateKey, error) {
	privateKeyPEM := os.Getenv(private)
	if privateKeyPEM == "" {
		log.Fatalf("PRIVATE_KEY environment variable not set")
	}

	// Replace escaped newlines with actual newlines
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, `\n`, "\n")

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKCS1PrivateKey: %v", err)
	}

	return privateKey, nil
}

func ReadPublicKeyFromEnv(public string) (*rsa.PublicKey, error) {
	publicKeyPEM := os.Getenv(public)
	if publicKeyPEM == "" {
		log.Fatalf("PUBLIC_KEY environment variable not set")
	}
	publicKeyPEM = strings.ReplaceAll(publicKeyPEM, `\n`, "\n")
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.ParsePKIXPublicKey: %v", err)
	}
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not ok: %v", err)
	}
	return publicKey, nil
}

func GenerateRSAPem(privateFilename string, publicFilename string, bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %v", err)
	}
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	privFile, err := os.Create(privateFilename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", privateFilename, err)
	}
	defer privFile.Close()
	err = pem.Encode(privFile, privPem)
	if err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}
	publicKey := &privateKey.PublicKey
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}
	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	pubFile, err := os.Create(publicFilename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", publicFilename, err)
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, pubPem)
	if err != nil {
		return fmt.Errorf("failed to write public key to file: %v", err)
	}
	return nil
}

func GenerateRSAEnv(privateKeyPath string) (string, string, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read private key file: %v", err)
	}
	cleanPrivateKey := CleanPEMString(string(privateKeyPEM))
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", "", fmt.Errorf("failed to decode private key PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse private key: %v", err)
	}
	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
	cleanPublicKey := CleanPEMString(string(publicKeyPEM))

	return cleanPrivateKey, cleanPublicKey, nil
}

func CleanPEMString(pem string) string {
	pem = strings.ReplaceAll(pem, "\n", `\n`)
	return pem
}

func SignedJWT(mongoenv *mongo.Database, user models.Users, userAgent string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 2).Unix()
	issuedTime := time.Now().Unix()
	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": expirationTime,
		"iat": issuedTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := ReadPrivateKeyFromEnv("private_key")
	if err != nil {
		return "", fmt.Errorf("error loading private key: %v", err)
	}

	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing string: %v", err)
	}

	session := models.Session{
		Username:   user.Username,
		Token:      t,
		User_Agent: userAgent,
		Expire_At:  time.Unix(expirationTime, 0),
		Issued_At:  time.Unix(issuedTime, 0),
	}
	InsertSession(SetConnection(), session)
	return t, nil
}

func DecodeJWT(r *http.Request) (datauser models.Users) {
	tokenString := r.Header.Get("Authorization")
	parts := strings.Split(tokenString, " ")
	tokenString = parts[1]
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return ReadPublicKeyFromEnv("public_key")
	})
	claims := token.Claims.(jwt.MapClaims)
	username := claims["sub"].(string)
	datauser = FindUser(SetConnection(), models.Users{Username: username})
	return datauser
}
