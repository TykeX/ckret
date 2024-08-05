# ckret

## ckret can be used in two ways.

### Using default instance
```go
ckret.Init(&aws.Config({Region: aws.String("ap-south-1")}), "ckret/local") // initialize default ckret instance
ckret.GetCkret() // get secret value from default ckret instance
ckret.Name() // --> ckret/local i.e. secretId or name of the secret
```

### Using new instance
```go
ck := ckret.New(&aws.Config({Region: aws.String("ap-south-1")}), "ckret/mango") // create and initialize new ckret instance
ck.GetCkret() // get secret value from the instance
ckret.Name() // --> ckret/mango i.e. secretId or name of the secret
```

### Warning & Suggestions
1. Avoid using default instance.
2. While using default instance make sure it has been initialized.
3. Don't called Init() more than once.


