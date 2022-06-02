package ckret

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var ckretCache *secretcache.Cache

func Init(config *aws.Config) error {
	c, err := secretcache.New(func(cache *secretcache.Cache) {
		sess, _ := session.NewSession(config)
		cache.Client = secretsmanager.New(sess)
	})
	if err != nil {
		return err
	}
	ckretCache = c
	return nil
}

func GetInstance() *secretcache.Cache {
	return ckretCache
}

// This will look for secret in aws secret manager based on environment variable
// warning: function is very specific to perticular use case.
// aws secret must be valid JSON
func GetCkret() (data map[string]any) {
	var secretName string = ""
	switch strings.ToLower(os.Getenv("ENVIRONMENT")) {
	case "prod", "prduction":
		secretName = "ckret/prod"
	case "stage", "staging":
		secretName = "ckret/stage"
	case "dev", "development":
		secretName = "ckret/dev"
	default:
		secretName = "ckret/local"
	}
	s, err := ckretCache.GetSecretString(secretName)
	if err != nil {
		panic("can not read ckret from secret manager")
	}
	json.Unmarshal([]byte(s), &data)
	return
}
