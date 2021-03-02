package global

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func StrLen(str string) int {
	return strings.Count(str, "") - 1
}

func PwdPlaintext2CipherText(pwd string, salt string) string {
	pwd = salt + "{_}" + pwd + "{_}" + salt
	has := md5.Sum([]byte(pwd))
	return fmt.Sprintf("%x", has)
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateRangeNum(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(max - min) + min
	return randNum
}

func Hour2Unix(hour string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormatStr, time.Now().Format(DateFormatStr) + " " + hour, time.Local)
}

func Md5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func Json2Map(j string) map[string]interface{} {
	r := make(map[string]interface{})
	_ = json.Unmarshal([]byte(j), &r)
	return r
}

func FileExists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}


func RandFloats(min, max float64, n int) float64 {
	rand.Seed(time.Now().UnixNano())
	res := min + rand.Float64() * (max - min)
	res, _ =  strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(n)+"f", res), 64)
	return res
}