package ckret

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

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
		cache.CacheItemTTL = int64(time.Minute) * 10 // 10 minute
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
var once sync.Once

func GetCkret() (data map[string]any) {
	var secretName string = ""
	switch strings.ToLower(os.Getenv("ENVIRONMENT")) {
	case "prod", "production":
		secretName = "ckret/prod"
	case "stage", "staging":
		secretName = "ckret/stage"
	case "dev", "development":
		secretName = "ckret/dev"
	case "sandbox":
		secretName = "ckret/sandbox"
	default:
		secretName = "ckret/local"
	}
	s, err := ckretCache.GetSecretString(secretName)
	if err != nil {
		panic("can not read ckret from secret manager")
	}
	once.Do(func() {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`{"selected_ckret":"%s"}`, secretName))
	})
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(`{"ckret_json_error":"%s"}`, err.Error()))
	}
	return
}
