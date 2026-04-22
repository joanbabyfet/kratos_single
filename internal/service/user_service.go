package service

import (
	"context"
	v1 "kratos_single/api/client/v1"
	"kratos_single/internal/biz"
	"kratos_single/internal/pkg/auth"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/jinzhu/copier"
)

//一套 proto = 一个 service 实现
type UserService struct {
	v1.UnimplementedUserServer
	uc *biz.UserUsecase
	log *log.Helper //不要全局用 log.xxx，应该依赖注入 logger
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc: uc,
		log: log.NewHelper(log.With(logger, "module", "user-service")),
	}
}

// 登录
func (s *UserService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {

	ip := auth.GetClientIp(ctx)

	user, err := s.uc.Login(ctx, req.Username, req.Password, ip)
	if err != nil {
		s.log.Errorw(
			"msg", "用户登录失败",
			"username", req.Username,
			"err", err,
		)
		return nil, err
	}

	var res v1.LoginReply
	if err := copier.Copy(&res, user); err != nil {
		return nil, err
	}

	//获取jwt token
	token, err := auth.GenerateToken(user.Id, "user", 0)
	if err != nil {
		return nil, err
	}
	res.Token = token

	return &res, nil
}

// 登出
func (s *UserService) Logout(ctx context.Context, req *v1.LogoutReq) (*v1.LogoutReply, error) {

	var token string

	// HTTP Header / gRPC metadata 取 Authorization
	if md, ok := metadata.FromServerContext(ctx); ok {
		token = md.Get("authorization")
	}

	if token == "" {
		return nil, errors.Unauthorized(
			"UNAUTHORIZED",
			"未登录",
		)
	}

	// 去掉 Bearer
	token = strings.TrimSpace(token)
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, errors.Unauthorized(
			"UNAUTHORIZED",
			"未登录",
		)
	}

	// 写入 Redis 黑名单
	// err := s.rdb.Set(ctx, "jwt:blacklist:"+token, 1, 24*time.Hour).Err()

	// if err != nil {
	// 	return nil, errors.InternalServer(
	// 		"REDIS_ERROR",
	// 		"退出失败",
	// 	)
	// }

	return &v1.LogoutReply{}, nil
}

// 注册
func (s *UserService) Register(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {

	var u biz.User

	// 自动拷贝
	if err := copier.Copy(&u, req); err != nil {
		return nil, err
	}

	u.RegIp = auth.GetClientIp(ctx)

	id, err := s.uc.Create(ctx, &u)
	if err != nil {
		s.log.Errorw(
			"msg", "用户注册失败",
			"username", req.Username,
			"err", err,
		)
		return nil, err
	}

	return &v1.RegisterReply{
		Id: id,
	}, nil
}

// 修改用户信息
func (s *UserService) UpdateProfile(ctx context.Context, req *v1.UpdateProfileReq) (*v1.UpdateProfileReply, error) {

	var u biz.User

	if err := copier.Copy(&u, req); err != nil {
		return nil, err
	}
	
	//从jwt获取用户id
	uid := auth.GetUser(ctx)
	u.Id = uid

	err := s.uc.Update(ctx, &u, "update")
	if err != nil {
		s.log.Errorw(
			"msg", "修改用户资料失败",
			"id", uid,
			"err", err,
		)
		return nil, err
	}

	return &v1.UpdateProfileReply{}, nil
}

// 获取用户信息
func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.GetUserReply, error) {
	
	//从jwt token获取用户id
	uid := auth.GetUser(ctx)
	//使用缓存
	u, err := s.uc.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}

	var res v1.GetUserReply

	if err := copier.Copy(&res, u); err != nil {
		return nil, err
	}

	return &res, nil
}

// 修改密码
func (s *UserService) SetPassword(ctx context.Context, req *v1.SetPasswordReq) (*v1.SetPasswordReply, error) {

	//从jwt获取用户id
	uid := auth.GetUser(ctx)

	err := s.uc.SetPassword(ctx, uid, req.OldPassword, req.NewPassword)
	if err != nil {
		s.log.Errorw(
			"msg", "修改密码失败",
			"id", uid,
			"err", err,
		)
		return nil, err
	}

	return &v1.SetPasswordReply{}, nil
}