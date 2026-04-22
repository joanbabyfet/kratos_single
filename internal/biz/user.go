package biz

import (
	"context"
	"fmt"
	"kratos_single/internal/pkg/auth"
	"kratos_single/internal/pkg/i18n"
	"kratos_single/internal/pkg/utils"

	"github.com/go-kratos/kratos/v2/errors"
)

//业务需要的数据
type User struct {
	Id           string
	Username     string
	Password     string
	Avatar       string
	Realname     string
	Sex          int8
	Email        string
	PhoneCode    string
	Phone        string
	Address      string
	Salt      	 string
	RegIp        string
	LoginTime    int
	LoginIp      string
	LoginCountry string
	Language     string
	Status       int8
	CreateUser   string
	UpdateUser   string
}

type UserQuery struct {
	Username    string
	Realname    string
	Email    string
	Phone    string
	Page     int
	PageSize int
	Limit int
}

type UserRepo interface {
	Create(context.Context, *User) (string, error)
	Update(context.Context, *User, string) error
	GetById(context.Context, string, bool) (*User, error)
	GetByUsername(context.Context, string) (*User, error)
	List(context.Context, *UserQuery) ([]*User, int64, error)
	Delete(context.Context, string, string) error
	UpdatePassword(context.Context, string, string) error
	UpdateLogin(context.Context, string, string) error
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) Login(ctx context.Context, username string, password string, ip string) (*User, error) {

	// 查询用户
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New(400, "-1", "用户不存在")
	}

	//校验密码
	if !utils.PasswordVerify(password, user.Password) {
		return nil, errors.New(400, "-2", "密码无效")
	}

	// 状态检查
	if user.Status == 0 {
		return nil, errors.New(400, "-3", "账号已被禁用")
	}

	// 更新登录信息（异步也可以）
	err = uc.repo.UpdateLogin(ctx, user.Id, ip)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// 添加
func (uc *UserUsecase) Create(ctx context.Context, a *User) (string, error) {
	if a.Username == "" {
		return "", errors.New(400, "-1", i18n.T(auth.GetLang(ctx), "InvalidUsername"))
	}
	if a.Password == "" {
		return "", errors.New(400, "-2", i18n.T(auth.GetLang(ctx), "InvalidPassword"))
	}
	if a.Email == "" {
		return "", errors.New(400, "-3", i18n.T(auth.GetLang(ctx), "InvalidEmail"))
	}
	if a.Phone == "" {
		return "", errors.New(400, "-4", i18n.T(auth.GetLang(ctx), "InvalidPhone"))
	}

	// ===== 密码加密 =====
	hash, err := utils.PasswordHash(a.Password)
	if err != nil {
		return "", errors.New(400, "-5", "密码加密错误")
	}
	a.Id = utils.UniqueId() 
	a.Language = "cn"
	a.Password = hash

	return uc.repo.Create(ctx, a)
}

// 修改（可选：普通更新 / FOR UPDATE / SHARE MODE）
func (uc *UserUsecase) Update(ctx context.Context, a *User, lockMode string) error {
	if a.Id == "" {
		return ErrInvalidID
	}
	return uc.repo.Update(ctx, a, lockMode)
}

func (uc *UserUsecase) SetPassword(ctx context.Context, id string, oldPwd string, newPwd string) error {

	user, err := uc.repo.GetById(ctx, id, false)
	if err != nil {
		return err
	}

	if user.Id == "" {
		return errors.New(400, "-1", "用户不存在")
	}
	
	// ===== 校验旧密码 =====
	if !utils.PasswordVerify(oldPwd, user.Password) {
		return errors.New(400, "-2", "旧密码错误")
	}

	// ===== 新密码加密 =====
	hash, err := utils.PasswordHash(newPwd)
	if err != nil {
		return errors.New(400, "-3", "密码加密失败")
	}

	return uc.repo.UpdatePassword(ctx, id, hash)
}

// 获取详情(走缓存)
func (uc *UserUsecase) GetById(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, true)
}

// 获取详情(后台不走缓存)
func (uc *UserUsecase) AdminGetById(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, false)
}

// 获取列表
func (uc *UserUsecase) List(ctx context.Context, q *UserQuery) ([]*User, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return uc.repo.List(ctx, q)
}

// 软删除
func (uc *UserUsecase) Delete(ctx context.Context, id string, user string) error {
	if id == "" {
		return ErrInvalidID
	}
	return uc.repo.Delete(ctx, id, user)
}

// 同步用户 (脚本调用)
func (uc *UserUsecase) SyncUser(ctx context.Context) error {

	fmt.Println("执行数据库同步用户")

	return nil
}