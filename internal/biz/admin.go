package biz

import (
	"context"
	"kratos_single/internal/pkg/auth"
	"kratos_single/internal/pkg/i18n"
	"kratos_single/internal/pkg/utils"

	"github.com/go-kratos/kratos/v2/errors"
)

//业务需要的数据
type Admin struct {
	Id           string
	Username     string
	Password     string
	Realname     string
	Sex          int8
	Email        string
	Salt         string
	RoleId       int
	RegIp        string
	LoginTime    int
	LoginIp      string
	LoginCountry string
	Status       int8
	CreateUser   string
	UpdateUser   string
}

type AdminQuery struct {
	Username    string
	Realname    string
	Email    string
	Phone    string
	Page     int
	PageSize int
	Limit int
}

type AdminRepo interface {
	Create(context.Context, *Admin) (string, error)
	Update(context.Context, *Admin) error
	GetById(context.Context, string, bool) (*Admin, error)
	GetByUsername(context.Context, string) (*Admin, error)
	List(context.Context, *AdminQuery) ([]*Admin, int64, error)
	Delete(context.Context, string, string) error
	UpdatePassword(context.Context, string, string) error
	UpdateLogin(context.Context, string, string) error
}

type AdminUsecase struct {
	repo AdminRepo
}

func NewAdminUsecase(repo AdminRepo) *AdminUsecase {
	return &AdminUsecase{repo: repo}
}

func (uc *AdminUsecase) Login(ctx context.Context, username string, password string, ip string) (*Admin, error) {

	// 查询管理员
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New(400, "-1", "管理员不存在")
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
func (uc *AdminUsecase) Create(ctx context.Context, a *Admin) (string, error) {
	if a.Username == "" {
		return "", errors.New(400, "-1", i18n.T(auth.GetLang(ctx), "InvalidUsername"))
	}
	if a.Password == "" {
		return "", errors.New(400, "-2", i18n.T(auth.GetLang(ctx), "InvalidPassword"))
	}
	if a.Email == "" {
		return "", errors.New(400, "-3", i18n.T(auth.GetLang(ctx), "InvalidEmail"))
	}

	// ===== 密码加密 =====
	hash, err := utils.PasswordHash(a.Password)
	if err != nil {
		return "", errors.New(400, "-5", "密码加密错误")
	}
	a.Id = utils.UniqueId() 
	a.Password = hash

	return uc.repo.Create(ctx, a)
}

// 修改
func (uc *AdminUsecase) Update(ctx context.Context, a *Admin) error {
	if a.Id == "" {
		return ErrInvalidID
	}
	return uc.repo.Update(ctx, a)
}

func (uc *AdminUsecase) SetPassword(ctx context.Context, id string, oldPwd string, newPwd string) error {

	user, err := uc.repo.GetById(ctx, id, false)
	if err != nil {
		return err
	}

	if user.Id == "" {
		return errors.New(400, "-1", "管理员不存在")
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
func (uc *AdminUsecase) GetById(ctx context.Context, id string) (*Admin, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, true)
}

// 获取详情(后台不走缓存)
func (uc *AdminUsecase) AdminGetById(ctx context.Context, id string) (*Admin, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	return uc.repo.GetById(ctx, id, false)
}

// 获取列表
func (uc *AdminUsecase) List(ctx context.Context, q *AdminQuery) ([]*Admin, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return uc.repo.List(ctx, q)
}

// 软删除
func (uc *AdminUsecase) Delete(ctx context.Context, id string, user string) error {
	if id == "" {
		return ErrInvalidID
	}
	return uc.repo.Delete(ctx, id, user)
}