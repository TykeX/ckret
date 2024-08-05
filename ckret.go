package ckret

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type instance struct {
	ckretCache *secretcache.Cache
	name       *string
}

// create new ckret instance using name (secretId) and aws config.
// panic on error
func New(config *aws.Config, name string) *instance {
	c, err := secretcache.New(func(cache *secretcache.Cache) {
		sess, _ := session.NewSession(config)
		cache.Client = secretsmanager.New(sess)
		cache.CacheItemTTL = int64(time.Minute) * 10 // 10 minute
	})
	if err != nil {
		panic("failed to initialize aws cache")
	}
	return &instance{
		ckretCache: c,
		name:       &name,
	}
}

// get secretId (name) of this ckret instance
func (i instance) Name() *string { return i.name }

// get value of this ckret instance
func (i instance) GetCkret() map[string]any {
	s, err := i.ckretCache.GetSecretString(*i.name)
	if err != nil {
		panic("failed to GetSecretString from cache")
	}
	var data map[string]any
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		panic("failed to unmarshal secret")
	}
	return data
}

var (
	// default ckret instance
	defaultInstance *instance = nil
	// initialize default ckret instance
	Init = func(config *aws.Config, name string) { defaultInstance = New(config, name) }
	// get value of default ckret instance
	GetCkret = func() map[string]any { return defaultInstance.GetCkret() }
	// get secretId (name) for default ckret instance
	Name = func() *string { return defaultInstance.Name() }
)
