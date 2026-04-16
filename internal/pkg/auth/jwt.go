package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/scrypt"
)

const (
	SecretKEY              string = "JWT-Secret-Key"
	DEFAULT_EXPIRE_SECONDS int    = 3600 // 默认10分钟, 单位秒
	PasswordHashBytes             = 16
	RoleUser  string 			  = "user"		//用户
    RoleAdmin string 			  = "admin"		//管理员
)

// MyCustomClaims
// This struct is the payload
// 此结构是有效负载
type MyCustomClaims struct {
	UserID string `json:"userID"`
	Role   string   `json:"role"`   //区分用户类型
	jwt.StandardClaims
}

// JwtPayload
// This struct is the parsing of token payload
// 此结构是对token有效负载的解析
type JwtPayload struct {
	Username  string `json:"username"`
	UserID    string `json:"userID"`
	Role   	  string `json:"role"`
	IssuedAt  int64  `json:"iat"` // 发布日期
	ExpiresAt int64  `json:"exp"` // 过期时间
}

// GenerateToken
// @Title GenerateToken
// @Description "生成token"
// @Param loginInfo 		*models.User 	"登录请求"
// @Param userID 			int 					"用户ID"
// @Param expiredSeconds 	int 					"过期时间"
// @return    tokenString   string         	"编码后的token"
// @return    err   		error         	"错误信息"
func GenerateToken(userID string, role string, expiredSeconds int) (tokenString string, err error) {
	// 如果没设置过期时间，默认为 DEFAULT_EXPIRE_SECONDS 600s
	if expiredSeconds == 0 {
		expiredSeconds = DEFAULT_EXPIRE_SECONDS
	}

	// 创建声明
	mySigningKey := []byte(SecretKEY)
	// 过期时间 = 当前时间（/s）+ expiredSeconds（/s）
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	log.Info("Token 将到期于：" + time.Unix(expireAt, 0).String())

	claims := MyCustomClaims{
		userID,
		role,
		jwt.StandardClaims{
			Issuer:    "your-system",     // 发行者
			IssuedAt:  time.Now().Unix(), // 发布时间
			ExpiresAt: expireAt,          // 过期时间
		},
	}

	// 利用上面创建的声明 生成token
	// NewWithClaims(签名算法 SigningMethod, 声明 Claims) *Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 利用密钥对token签名
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		return "",
			errors.New("错误: token生成失败！")
	}
	return tokenStr, nil
}

// ValidateToken
// @Title ValidateToken
// @Description "验证token"
// @Param tokenString 		string 	"编码后的token"
// @return   	*JwtPayload     "Jwt有效负载的解析"
// @return   	error         	"错误信息"
func ValidateToken(tokenString string) (*JwtPayload, error) {
	// 获取编码前的token信息
	token, err := jwt.ParseWithClaims(tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})
	// 获取payload-声明内容
	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		log.Info(fmt.Sprintf("%v %v",
			claims.UserID,
			claims.StandardClaims.ExpiresAt, // 过期时间
		))
		log.Info("Token 将过期于：" + time.Unix(claims.StandardClaims.ExpiresAt, 0).String())
		return &JwtPayload{
			Username:  claims.StandardClaims.Issuer, // 用户名：发行者
			UserID:    claims.UserID,
			Role:	   claims.Role,
			IssuedAt:  claims.StandardClaims.IssuedAt,
			ExpiresAt: claims.StandardClaims.ExpiresAt,
		}, nil
	} else {
		log.Info(err.Error())
		return nil, errors.New("错误: token验证失败")
	}
}

// RefreshToken
// @Title RefreshToken
// @Description "更新token"
// @Param tokenString 		string 		"编码后的token"
// @return   newTokenString string    "编码后的新的token"
// @return   err   			error     "错误信息"
func RefreshToken(tokenString string) (newTokenString string, err error) {
	// 获取上一个token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})
	// 获取上一个token 的 payload-声明
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", err
	}

	// 创建新的声明
	mySigningKey := []byte(SecretKEY)
	expireAt := time.Now().Add(time.Second * time.Duration(DEFAULT_EXPIRE_SECONDS)).Unix() //new expired
	newClaims := MyCustomClaims{
		claims.UserID,
		claims.Role,
		jwt.StandardClaims{
			Issuer:    claims.StandardClaims.Issuer, //name of token issue
			IssuedAt:  time.Now().Unix(),            //time of token issue
			ExpiresAt: expireAt,
		},
	}

	// 利用新的声明，生成新的token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	// 利用签名算法对新的token进行签名
	tokenStr, err := newToken.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("错误: 新的新json web token 生成失败！")
	}

	return tokenStr, nil
}

// GenerateSalt
// @Title GenerateSalt
// @Description "生成用户的加密的钥匙|generate salt"
// @return   salt 		string    "生成用户的加密的钥匙"
// @return   err   		error     "错误信息"
func GenerateSalt() (salt string, err error) {
	buf := make([]byte, PasswordHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", errors.New("error: failed to generate user's salt")
	}

	return fmt.Sprintf("%x", buf), nil
}

// GeneratePassHash
// @Title GenerateSalt
// @Description "对密码加密|generate password hash"
// @Param password 		string 		"用户登录密码"
// @Param salt 			string 		"用户的加密的钥匙"
// @return    hash   string         "加密后的密码"
// @return    err    error         	"错误信息"
func GeneratePassHash(password string, salt string) (hash string, err error) {
	h, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, PasswordHashBytes)
	if err != nil {
		return "", errors.New("error: failed to generate password hash")
	}

	return fmt.Sprintf("%x", h), nil
}
