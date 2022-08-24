package ckret

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestGetCkret(t *testing.T) {
	Init(&aws.Config{Region: aws.String("ap-south-1")})
	os.Setenv("ENVIRONMENT", "sandbox")
	fmt.Println(GetCkret())
}
