package char

import (
	"bytes"
	"math/rand"
)

const (
	lowerSnake byte = iota
	upperSnake
)

func IsUpper(character byte) bool {
	return character >= 'A' && character <= 'Z'
}

func IsLower(character byte) bool {
	return character >= 'a' && character <= 'z'
}

func camel2snake(str string, style byte) string {
	strBytes := []byte(str)
	strLen := len(strBytes)
	result := bytes.Buffer{}
	for pos, ascii := range strBytes {
		changedASCII := ascii
		if style == lowerSnake && IsUpper(ascii) {
			changedASCII += 32
		} else if style == upperSnake && IsLower(ascii) {
			changedASCII -= 32
		}
		result.WriteByte(changedASCII)
		if pos >= strLen-1 {
			continue
		}
		nextChar := strBytes[pos+1]
		if IsLower(ascii) && IsUpper(nextChar) {
			result.WriteByte('_')
		}
	}

	return result.String()
}

// CamelToUpperSnake 大驼峰字符串转换为大写蛇形
func CamelToUpperSnake(str string) string {
	return camel2snake(str, upperSnake)
}

// CamelToLowerSnake 大驼峰字符串转换为小写蛇形
func CamelToLowerSnake(str string) string {
	return camel2snake(str, lowerSnake)
}

// RandomStrings 生成随机字符串
func RandomStrings(length int) string {
	base := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	baseLen := len(base)
	shuffle := make([]byte, length)

	rand.Shuffle(baseLen, func(a int, b int) {
		base[a], base[b] = base[b], base[a]
	})

	for range shuffle {
		pos := rand.Intn(baseLen)
		shuffle = append(shuffle, base[pos])
	}

	return string(shuffle)
}
