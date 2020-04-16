package utils

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"time"
	"unicode"

	"github.com/Unknwon/goconfig"

	"github.com/gin-gonic/gin"
)

func Consolelog(contents ...interface{}) {
	cfg := GetCfg()
	enviroment, _ := cfg.GetValue("project", "enviroment")
	if enviroment == "debug" {
		for _, v := range contents {
			fmt.Fprintln(gin.DefaultWriter, v)
		}
	}
}

func Int64to32(i64 int64) int {
	i32, _ := strconv.Atoi(strconv.FormatInt(i64, 10))
	return i32
}

func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}

func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(str))])
	}
	return string(result)
}

func SaveFile(file multipart.File, ext string, category string) (string, error) {
	cfg := GetCfg()
	uploadDir, _ := cfg.GetValue("project", "uploaddir")
	dateString := time.Now().Format("20060102150405")
	filename := dateString + GetRandomString(20) + ext
	filepath := uploadDir + "/" + category + "/" + filename
	storepath := "upload/" + category + "/" + filename
	os.MkdirAll(path.Dir(filepath), os.ModePerm)
	out, err := os.Create(filepath)
	defer out.Close()
	_, err = io.Copy(out, file)
	return storepath, err
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetCfg() *goconfig.ConfigFile {
	configFilePath := "config/config.ini"
	cfg, _ := goconfig.LoadConfigFile(configFilePath)
	return cfg
}

func GetTestCfg(path string) *goconfig.ConfigFile {
	cfg, _ := goconfig.LoadConfigFile(path)
	return cfg
}

func ConvertTextImgUrl(origin string) string {
	server, _ := GetCfg().GetValue("project", "server")
	replaceTo := `<img src="` + server + "/"
	reg := regexp.MustCompile(`<img src="`)
	res := reg.ReplaceAllString(origin, replaceTo)
	return res
}

func EngNumToInt(engNum string) (int, error) {
	num := 0
	for _, v := range engNum {
		if !unicode.IsLetter(v) {
			return num, errors.New("格式错误")
		}
		if unicode.IsUpper(v) {
			num = num*26 + (int(v) - 64)
		} else {
			num = num*26 + (int(v) - 96)
		}
	}
	return num, nil
}

func CheckError(c *gin.Context, err error, msg string) error {
	if err != nil {
		Consolelog(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": msg,
			"show":    true,
		})
		return err
	} else {
		return nil
	}
}
