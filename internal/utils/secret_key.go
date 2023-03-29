package utils

import (
	"github.com/catfishlty/webhooks-hub/internal/common"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetSecretKey(inputSecretKey string) string {
	if inputSecretKey != "" {
		log.Debugf("using input secret key: %s", inputSecretKey)
		return inputSecretKey
	}
	bytes, err := os.ReadFile(common.SecretKeyFile)
	if err == nil {
		secretKey := string(bytes)
		log.Debugf("using from secret key file: %s", secretKey)
		return secretKey
	}
	log.Debugf("failed to read secret key file: %s", err.Error())
	secretKey := NewRandom().String(64)
	err = os.WriteFile(common.SecretKeyFile, []byte(secretKey), 0644)
	if err != nil {
		log.Fatal("failed to write secret key", err)
	}
	log.Debugf("using random secret key: %s, and save to %s", secretKey, common.SecretKeyFile)
	return secretKey
}
