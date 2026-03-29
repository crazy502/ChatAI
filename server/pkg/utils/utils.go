package utils

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < num; i++ {
		code += strconv.Itoa(r.Intn(10))
	}
	return code
}

func GenerateUUID() string {
	return uuid.New().String()
}
