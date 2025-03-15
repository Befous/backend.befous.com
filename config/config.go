package config

import (
	"os"
)

var PrivateKey string = os.Getenv("privatekey")
var PublicKey string = os.Getenv("publickey")
