package service

import (
	"context"
	pb "github.com/piklen/pb/user"
	logging "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
	"strconv"
	"user/dao"
	"user/model"
	"user/pkg/e"
	"user/serializer"
)

func (s *Server) UserCreateAddress(ctx context.Context, in *pb.UserCreateAddressRequest) (*pb.CommonResponse, error) {
	code := e.Success
	addressDao := dao.NewAddressDao(ctx)
	userId, err := strconv.Atoi(in.UserId)
	address := &model.Address{
		UserID:  uint(userId),
		Name:    in.Name,
		Phone:   in.Phone,
		Address: in.Address,
	}
	err = addressDao.CreateAddress(address)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "创建用户地址错误！！！！",
		}, nil
	}
	addressDao = dao.NewAddressDaoByDB(addressDao.DB)
	var addresses []*model.Address
	addresses, err = addressDao.ListAddressByUid(uint(userId))
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "获取用户全部地址失败！！！！",
		}, nil
	}
	addressesList := serializer.BuildAddresses(addresses)
	addressesMap := make(map[string]interface{})
	addressesData := make(map[string]interface{})
	for i, addr := range addressesList {
		addressesMap[strconv.Itoa(i)] = map[string]interface{}{
			"id":        addr.ID,
			"user_id":   addr.UserID,
			"name":      addr.Name,
			"phone":     addr.Phone,
			"address":   addr.Address,
			"seen":      addr.Seen,
			"create_at": addr.CreateAt,
		}
	}
	addressesData["addresses"] = addressesMap
	spb, err := structpb.NewStruct(addressesData)
	if err != nil {
		code = 400
		return &pb.CommonResponse{
			StatusCode:   int64(code),
			Message:      e.GetMsg(code),
			ResponseData: "序列化解析json失败！！！！",
		}, nil
	}
	return &pb.CommonResponse{
		StatusCode:       int64(code),
		Message:          e.GetMsg(code),
		ResponseDataJson: spb, // 直接使用spb作为响应数据
		ResponseData:     "用户全部地址返回成功！！！！",
	}, nil
}
