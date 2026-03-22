package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"server/model"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < num; i++ {
		digit := r.Intn(10)
		code += strconv.Itoa(digit)
	}
	return code
}

func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func IsBcryptHash(hash string) bool {
	_, err := bcrypt.Cost([]byte(hash))
	return err == nil
}

func GenerateUUID() string {
	return uuid.New().String()
}

func ConvertToModelMessage(sessionID string, userName string, msg *schema.Message) *model.Message {
	isUser := false
	if msg.Role == schema.User {
		isUser = true
	}
	return &model.Message{
		IdempotencyKey: GenerateUUID(),
		SessionID:      sessionID,
		UserName:       userName,
		Content:        msg.Content,
		IsUser:         isUser,
	}
}

func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsgs := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		role := schema.Assistant
		if m.IsUser {
			role = schema.User
		}
		schemaMsgs = append(schemaMsgs, &schema.Message{
			Role:    role,
			Content: m.Content,
		})
	}
	return schemaMsgs
}
