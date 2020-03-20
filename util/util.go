package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"hash"
	"io"
	"os"
	"path/filepath"
	"fmt"
	"time"
	db "go_dev/video_service/db"
)

const (
	PWD_SALT = "sfdk4r97EHko"
	TOKEN_SALT = "jidsajfi2SXfhUd6i"
	TOKEN_INVALID_TIME = 86400
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
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

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

//生成用户登录的token
func GenToken(username string) string {
	//40位的一个token md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := MD5([]byte(ts+username+TOKEN_SALT))
    return tokenPrefix + ts[:8]
}

//验证token是否有效
func IsTokenValid(token string, username string) bool {
	//1 验证token是否过期

	if len(token) != 40 {
		return false
	}

	
	ts := fmt.Sprintf("%x",time.Now().Unix())[:8]
	usertime,_ := strconv.ParseUint(token[32:], 16, 32)
	now ,_ := strconv.ParseUint(ts[:8], 16, 32)

	if usertime + TOKEN_INVALID_TIME <  now {
		fmt.Println("token超时")
		return false
	}
	
	//2 数据库查询username对应的token信息
   istokenused :=	db.CheckUserToken(username,token)
   if !istokenused {
	   fmt.Println("token错误")
	   return false
   }

	return true

}