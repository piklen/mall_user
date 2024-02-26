package util

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成验证码函数
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}
