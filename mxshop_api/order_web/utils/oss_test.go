package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"testing"
)

func TestOss(t *testing.T) {
	fmt.Println("OSS Go SDK Version:", oss.Version)
}
