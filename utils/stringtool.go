package utils

import "fmt"
import "crypto/md5"
import "math/rand"
import "time"
import "strconv"
import "strings"

//將字符串加密成 md5
func String2md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //將[]byte轉成16進制
}

//RandomString 在數字、大寫字母、小寫字母範圍內生成num位的隨機字符串
func RandomString(length int) string {
	// 48 ~ 57 數字
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
	// 一共62個字符，在0~61進行隨機，小於10時，在數字範圍隨機，
	// 小於36在大寫範圍內隨機，其他在小寫範圍隨機
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	for i := 0; i < length; i++ {
		t := rand.Intn(62)
		if t < 10 {
			result = append(result, strconv.Itoa(rand.Intn(10)))
		} else if t < 36 {
			result = append(result, string(rand.Intn(26)+65))
		} else {
			result = append(result, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(result, "")
}
