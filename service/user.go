package service

import (
	"context"
	pb "github.com/piklen/pb/user"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"user/dao"
	"user/model"
	"user/pkg/e"
	"user/pkg/util"
	"user/serializer"
	//"user/serializer"
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

// // Login 用户登陆函数
//
//	func (service *Server) UserLogin(ctx context.Context, in *pb.UserRegisterRequest) (*pb.CommonResponse, error) {
//		var user *model.User
//		code := e.Success
//		userDao := dao.NewUserDao(ctx)
//
//		// 判断用户是否存在
//		user, exist, err := userDao.ExistOrNotByUserName(in.UserName)
//		if err != nil {
//			code = e.ErrorDatabase
//			return &pb.CommonResponse{
//				StatusCode:   int64(code),
//				Message:      e.GetMsg(code),
//				ResponseData: "数据库查询错误",
//			}, nil
//		}
//		if !exist {
//			code = e.ErrorExistUserNotFound
//			return &pb.CommonResponse{
//				StatusCode:   int64(code),
//				Message:      e.GetMsg(code),
//				ResponseData: "用户不存在,请先注册!!!",
//			}, nil
//		}
//
//		// 校验密码
//		if !user.CheckPassword(in.Password) {
//			code = e.ErrorNotCompare
//			return &pb.CommonResponse{
//				StatusCode:   int64(code),
//				Message:      e.GetMsg(code),
//				ResponseData: "密码错误,请重新输入密码!!!",
//			}, nil
//		}
//
//		// Token签发
//		token, err := util.GenerateToken(user.ID, in.UserName, 0)
//		if err != nil {
//			code = e.ErrorAuthToken
//			return &pb.CommonResponse{
//				StatusCode:   int64(code),
//				Message:      e.GetMsg(code),
//				ResponseData: "Token签发失败!!",
//			}, nil
//		}
//		// 将User结构体转换为map[string]interface{}
//		builtUser := serializer.BuildUser(user)
//		userMap := map[string]interface{}{
//			"ID":       builtUser.ID,
//			"UserName": builtUser.UserName,
//			"NickName": builtUser.NickName,
//			"Email":    builtUser.Email,
//			"Status":   builtUser.Status,
//			"Avatar":   builtUser.Avatar,
//			"CreateAt": builtUser.CreateAt,
//		}
//
//		// 将数据转换为google.protobuf.Struct
//		dataMap := map[string]interface{}{
//			"User":  userMap,
//			"Token": token,
//		}
//		spb, err := structpb.NewStruct(dataMap)
//		if err != nil {
//			log.Fatal("Failed to convert struct to google.protobuf.Struct: ", err)
//		}
//
//		return &pb.CommonResponse{
//			StatusCode:       int64(code),
//			Message:          e.GetMsg(code),
//			ResponseDataJson: spb.MarshalJSON(),
//			//ResponseData:
//		}, nil
//	}
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
	builtUser := serializer.BuildUser(user)
	userMap := map[string]interface{}{
		"ID":       builtUser.ID,
		"UserName": builtUser.UserName,
		"NickName": builtUser.NickName,
		"Email":    builtUser.Email,
		"Status":   builtUser.Status,
		"Avatar":   builtUser.Avatar,
		"CreateAt": builtUser.CreateAt,
	}

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

//func (service *Server) UserLogin(ctx context.Context, in *pb.UserRegisterRequest) (*pb.CommonResponse, error) {
//	var user *model.User
//	code := e.Success
//	userDao := dao.NewUserDao(ctx)
//	//判断用户是否存在
//	user, exist, err := userDao.ExistOrNotByUserName(in.UserName)
//	if !exist || err != nil { // 如果查询不到，返回相应的错误
//		code = e.ErrorExistUserNotFound
//		return &pb.CommonResponse{
//			StatusCode:   int64(code),
//			Message:      e.GetMsg(code),
//			ResponseData: "用户不存在,请先注册!!!",
//		}, nil
//	}
//	//校验密码
//	if user.CheckPassword(in.Password) == false {
//		code = e.ErrorNotCompare
//		return &pb.CommonResponse{
//			StatusCode:   int64(code),
//			Message:      e.GetMsg(code),
//			ResponseData: "密码错误,请重新输入密码!!!",
//		}, nil
//	}
//	//http 无状态(服务器需要token来认证)
//	//Token签发
//	token, err := util.GenerateToken(user.ID, in.UserName, 0)
//	if err != nil {
//		code = e.ErrorAuthToken
//		return &pb.CommonResponse{
//			StatusCode:   int64(code),
//			Message:      e.GetMsg(code),
//			ResponseData: "Token签发失败!!",
//		}, nil
//	}
//	p := serializer.TokenData{User: serializer.BuildUser(user), Token: token}
//	// 将结构体转换为JSON字符串
//	jsonString, err := json.Marshal(p)
//	if err != nil {
//		log.Fatal("JSON marshaling failed: ", err)
//	}
//	dataMap := make(map[string]interface{})
//	value := reflect.ValueOf(p)
//	for i := 0; i < value.NumField(); i++ {
//		fieldName := value.Type().Field(i).Name
//		fieldValue := value.Field(i).Interface()
//		dataMap[fieldName] = fieldValue
//	}
//	dataMap = map[string]interface{}{
//		"User":  serializer.BuildUser(user),
//		"Token": token,
//	}
//	spb, err := structpb.NewStruct(dataMap)
//	if err != nil {
//		log.Fatal("Failed to convert struct to google.protobuf.Struct: ", err)
//	}
//	//p := serializer.TokenData{User: serializer.BuildUser(user), Token: token}
//	//// 将结构体转换为map[string]interface{}
//	//dataMap := make(map[string]interface{})
//	//value := reflect.ValueOf(p)
//	//for i := 0; i < value.NumField(); i++ {
//	//	fieldName := value.Type().Field(i).Name
//	//	fieldValue := value.Field(i).Interface()
//	//	dataMap[fieldName] = fieldValue
//	//}
//	//// 将map[string]interface{}转换为google.protobuf.Struct
//	//spb, err := structpb.NewStruct(dataMap)
//	//if err != nil {
//	//	log.Fatal("Failed to convert struct to google.protobuf.Struct: ", err)
//	//}
//	//}
//	return &pb.CommonResponse{
//		StatusCode:       int64(code),
//		Message:          e.GetMsg(code),
//		ResponseDataJson: spb,
//		ResponseData:     string(jsonString),
//	}, nil
//}
