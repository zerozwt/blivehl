package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/zerozwt/blivehl/server/db"
)

type LoginService struct{}

var gLoginService LoginService = LoginService{}
var ErrUserBanned error = fmt.Errorf("user banned")
var ErrPassword error = fmt.Errorf("password incorrect")
var ErrParamTooLong error = fmt.Errorf("param too long")

const (
	MAX_LOGIN_PARAM_LEN = 64
)

func GetLoginService() *LoginService {
	return &gLoginService
}

func (s *LoginService) CheckCookie(cookie string) (string, bool, error) {
	userID, err := db.GetCookie(cookie)
	if err != nil {
		return "", false, err
	}

	return s.GetUserInfo(userID)
}

func (s *LoginService) LoginAndGenerateCookie(user, pass string) (string, error) {
	if len(user) > MAX_LOGIN_PARAM_LEN || len(pass) > MAX_LOGIN_PARAM_LEN {
		return "", ErrParamTooLong
	}

	userInfo, err := db.GetUser(user)
	if err != nil {
		return "", err
	}

	if userInfo.Banned {
		return "", ErrUserBanned
	}

	encPass := s.encryptByPassword(pass, pass)
	if encPass != userInfo.Password {
		return "", ErrPassword
	}

	cookie := s.encryptByPassword(pass, user+fmt.Sprint(time.Now().Unix()))
	err = db.PutCookie(cookie, user, time.Now().Unix()+24*3600)
	if err != nil {
		return "", err
	}

	return cookie, nil
}

func (s *LoginService) encryptByPassword(pass, value string) string {
	hash := md5.Sum([]byte(pass))
	key, _ := aes.NewCipher(hash[:])
	aead, _ := cipher.NewGCM(key)
	nonce := []byte("123456781234")
	encData := aead.Seal(nil, nonce, []byte(value), nil)
	return base64.RawStdEncoding.EncodeToString(encData)
}

func (s *LoginService) DeleteCookie(cookie string) error {
	return db.DeleteCookie(cookie)
}

func (s *LoginService) GetUserInfo(userID string) (string, bool, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		return "", false, err
	}

	if user.Banned {
		return "", false, ErrUserBanned
	}

	return user.Name, user.IsAdmin, nil
}

func (s *LoginService) ChangePassword(user, oldPass, newPass string) error {
	if len(oldPass) > MAX_LOGIN_PARAM_LEN || len(newPass) > MAX_LOGIN_PARAM_LEN {
		return ErrParamTooLong
	}

	userInfo, err := db.GetUser(user)
	if err != nil {
		return err
	}

	if userInfo.Banned {
		return ErrUserBanned
	}

	encPass := s.encryptByPassword(oldPass, oldPass)
	if encPass != userInfo.Password {
		return ErrPassword
	}

	userInfo.Password = s.encryptByPassword(newPass, newPass)

	return db.PutUser(userInfo)
}

func (s *LoginService) PutUser(user db.User) error {
	user.Password = s.encryptByPassword(user.Password, user.Password)
	return db.PutUser(user)
}
