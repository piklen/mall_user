package service

import (
	"context"
	"fmt"
	pb "github.com/piklen/pb/user"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/mail.v2"
	"log"
	"strconv"
	"strings"
	"time"
	"user/conf"
	"user/dao"
	"user/model"
	"user/pkg/e"
	"user/pkg/util"
	"user/serializer"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (s *Server) RegisterUser(ctx context.Context, in *pb.UserRegisterRequest) (*pb.CommonResponse, error) {
	var user model.User
	if in.Key == "" || len(in.Key) != 16 {
		return &pb.CommonResponse{
			StatusCode:   30009,
			Message:      e.GetMsg(30009),
			ResponseData: "注册失败",
		}, nil
	}
	//10000  ----->密文存储,对称加密操作
	util.Encrypt.SetKey(in.Key)
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(in.UserName)
	if err != nil {
		return &pb.CommonResponse{
			StatusCode:   30002,
			Message:      e.GetMsg(30002),
			ResponseData: "注册失败",
		}, nil
	}
	if exist {
		return &pb.CommonResponse{
			StatusCode:   30001,
			Message:      e.GetMsg(30001),
			ResponseData: "注册失败",
		}, nil
	}
	user = model.User{
		UserName: in.UserName,
		NickName: in.NickName,
		Status:   model.Active,
		Avatar:   "avatar.jpeg",
		Money:    util.Encrypt.AesEncoding("10000"), // 初始金额
	}
	// 加密密码
	//前端传入的是明文
	if err = user.SetPassword(in.Password); err != nil {
		return &pb.CommonResponse{
			StatusCode:   30002,
			Message:      e.GetMsg(30002),
			ResponseData: "注册失败",
		}, nil
	}
	// 创建用户
	err = userDao.CreateUser(&user) //传入指针,执行效率更高
	if err != nil {
		return &pb.CommonResponse{
			ResponseData: "创建用户失败！！！",
		}, nil
	}
	return &pb.CommonResponse{
		StatusCode:   200,
		Message:      e.GetMsg(200),
		ResponseData: "创建用户成功！！！",
	}, nil
}

func (service *Server) UserLogin(ctx context.Context, in *pb.UserRegisterRequest) (*pb.CommonResponse, error) {
	var user *model.User
	code := e.Success
	userDao := dao.NewUserDao(ctx)

	// 判断用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(in.UserName)
	if err != nil {
		code = e.ErrorDatabase
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "数据库查询错误",
		}, nil
	}
	if !exist {
		code = e.ErrorExistUserNotFound
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "用户不存在,请先注册!!!",
		}, nil
	}

	// 校验密码
	if !user.CheckPassword(in.Password) {
		code = e.ErrorNotCompare
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "密码错误,请重新输入密码!!!",
		}, nil
	}

	// Token签发
	token, err := util.GenerateToken(user.ID, in.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "Token签发失败!!",
		}, nil
	}

	// 将User结构体转换为map[string]interface{}
	userMap := serializer.BuildUser(user)
	// 将数据转换为google.protobuf.Struct
	dataMap := map[string]interface{}{
		"User":  userMap,
		"Token": token,
	}
	spb, err := structpb.NewStruct(dataMap)
	if err != nil {
		// 使用日志记录错误，而不是终止程序
		log.Printf("Failed to convert struct to google.protobuf.Struct: %v", err)
		code = 500
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "内部错误",
		}, nil
	}
	// 返回CommonResponse，包含protobuf.Struct类型的数据
	return &pb.CommonResponse{
		StatusCode:       int64(code),
		Message:          e.GetMsg(code),
		ResponseDataJson: spb, // 直接使用spb作为响应数据
	}, nil
}
func (s *Server) UpdateNickName(ctx context.Context, in *pb.UpdateNickNameRequest) (*pb.CommonResponse, error) {
	var user *model.User
	var err error
	code := e.Success
	// 找到用户
	userDao := dao.NewUserDao(ctx)
	userId, err := strconv.Atoi(in.UserId)
	user, err = userDao.GetUserById(uint(userId))
	if err != nil {
		code = e.Error
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "找不到该user_id!!!",
		}, nil
	}
	//修改昵称Nickname
	if in.NickName != "" {
		user.NickName = in.NickName
	}
	err = userDao.UpdateUserById(uint(userId), user)
	if err != nil {
		code = e.Error
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "修改数据库昵称错误！！！",
		}, nil
	}
	// 将User结构体转换为map[string]interface{}
	userMap := serializer.BuildUser(user)
	// 将数据转换为google.protobuf.Struct
	dataMap := map[string]interface{}{
		"User": userMap,
	}
	spb, err := structpb.NewStruct(dataMap)
	return &pb.CommonResponse{
		StatusCode:       int64(code),
		Message:          e.GetMsg(code),
		ResponseDataJson: spb, // 直接使用spb作为响应数据
		ResponseData:     "修改昵称成功！！！",
	}, nil
}

// UploadAvatar 头像更新
func (s *Server) UploadAvatar(ctx context.Context, in *pb.UploadAvatarRequest) (*pb.CommonResponse, error) {
	code := e.Success
	var user *model.User
	var err error
	userDao := dao.NewUserDao(ctx)
	userId, err := strconv.Atoi(in.UserId)
	if err != nil {
		return &pb.CommonResponse{
			StatusCode:   int64(500),
			Message:      e.GetMsg(code),
			ResponseData: "userId非法！！！！",
		}, nil
	}
	user, err = userDao.GetUserById(uint(userId))
	if err != nil {
		return &pb.CommonResponse{
			StatusCode:   int64(500),
			Message:      e.GetMsg(code),
			ResponseData: "数据库查询userId错误！！！",
		}, nil
	}
	//保存图片到本地
	path, err := UploadAvatarToLocalStatic(in.FileData, in.UserId, user.UserName)
	if err != nil {
		code = e.ErrorUploadAvatar
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "用户头像保存到本地失败！！！",
		}, nil
	}
	user.Avatar = path
	err = userDao.UpdateUserById(uint(userId), user)
	if err != nil {
		code = e.Error
		return &pb.CommonResponse{
			StatusCode:   int64(500),
			Message:      e.GetMsg(code),
			ResponseData: "更新用户数据失败！！！",
		}, nil
	}
	// 将User结构体转换为map[string]interface{}
	userMap := serializer.BuildUser(user)
	// 将数据转换为google.protobuf.Struct
	dataMap := map[string]interface{}{
		"User": userMap,
	}
	spb, err := structpb.NewStruct(dataMap)
	return &pb.CommonResponse{
		StatusCode:       int64(code),
		Message:          e.GetMsg(code),
		ResponseDataJson: spb, // 直接使用spb作为响应数据
		ResponseData:     "上传头像成功！！！！",
	}, nil
}

// SendEmail 发送邮件
func (s *Server) SendEmail(ctx context.Context, in *pb.SendEmailRequest) (*pb.CommonResponse, error) {
	code := e.Success
	var address string
	var notice *model.Notice //绑定邮箱,修改密码都有模板通知
	userId, err := strconv.Atoi(in.UserId)
	operationType, err := strconv.Atoi(in.OperationType)
	token, err := util.GenerateEmailToken(uint(userId), uint(operationType), in.Email, in.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "签发邮箱验证Token！！！",
		}, nil
	}
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(uint(operationType))
	if err != nil {
		code = e.Error
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "获取notice类型失败！！！",
		}, nil
	}
	address = conf.ValidEmail + token //发送方
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "Email", address, -1) //字符串替换
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", in.Email)
	m.SetHeader("Subject", "xiaobao")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		fmt.Println("err", err)
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "发送邮件失败！！！！",
		}, nil
	}
	return &pb.CommonResponse{
		StatusCode:   int64(code),
		Message:      e.GetMsg(code),
		ResponseData: "发送邮件成功！！！",
	}, nil
}

// ValidEmail 验证邮箱
func (s *Server) ValidEmail(ctx context.Context, in *pb.ValidEmailRequest) (*pb.CommonResponse, error) {
	var userID uint
	var email string
	var password string
	var operationType uint
	code := e.Success

	// 验证token
	if in.Token == "" {
		code = e.InvalidParams
	} else {
		claims, err := util.ParseEmailToken(in.Token)
		if err != nil {
			//如果解析token错误就返回错误
			code = e.ErrorAuthToken
		} else if time.Now().Unix() > claims.ExpiresAt {
			//如果超时就返回验证时间超时
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			//不然就是成功了，就直接构建用户结构体
			userID = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
		}
	}
	if code != e.Success {
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "发生错误！！！",
		}, nil
	}

	// 获取该用户信息
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(userID)
	if err != nil {
		code = e.Error
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "验证邮件时，在获取用户信息时失败！！！",
		}, nil
	}
	if operationType == 1 {
		// 1:绑定邮箱
		user.Email = email
	} else if operationType == 2 {
		// 2：解绑邮箱
		user.Email = ""
	} else if operationType == 3 {
		// 3：修改密码
		err = user.SetPassword(password)
		if err != nil {
			code = e.Error
			return &pb.CommonResponse{
				StatusCode:   int64(code),
				Message:      e.GetMsg(code),
				ResponseData: "更改用户密码错误！！！",
			}, nil
		}
	}
	err = userDao.UpdateUserById(userID, user)
	if err != nil {
		code = e.Error
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "更新用户邮件内容时发生错误！！！",
		}, nil
	}
	// 成功则返回用户的信息
	// 将User结构体转换为map[string]interface{}
	userMap := serializer.BuildUser(user)
	// 将数据转换为google.protobuf.Struct
	dataMap := map[string]interface{}{
		"User": userMap,
	}
	spb, err := structpb.NewStruct(dataMap)
	return &pb.CommonResponse{
		StatusCode:       int64(code),
		Message:          e.GetMsg(code),
		ResponseDataJson: spb, // 直接使用spb作为响应数据
		ResponseData:     "验证邮箱成功！！！",
	}, nil
}
