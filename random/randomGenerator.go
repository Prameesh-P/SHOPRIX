package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphaBet = "abcdefghijklmnopqrstuvwxyz"
const Gmail = "@gmail.com"

func init() {
	rand.Seed(time.Now().UnixNano())
}
func RandomGmailGenerator(length int) string {
	var sb strings.Builder
	k := len(alphaBet)
	for i := 0; i < length; i++ {
		c := alphaBet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	fmt.Println(sb.String() + Gmail)
	return sb.String() + Gmail
}
func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphaBet)
	for i := 0; i < length; i++ {
		c := alphaBet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomInteger(i, j int64) int64 {
	return i + rand.Int63n(j-i+1)
}
