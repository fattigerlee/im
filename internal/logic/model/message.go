package model

import (
	"im/pkg/logger"
	"im/pkg/pb"
	"im/pkg/util"
	"strconv"
	"strings"
	"time"
)

// Message 消息
type Message struct {
	Id           int64     // 自增主键
	UserId       int64     // 所属类型id
	RequestId    int64     // 请求id
	SenderType   int32     // 发送者类型
	SenderId     int64     // 发送者账户id
	ReceiverType int32     // 接收者账户id
	ReceiverId   int64     // 接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id
	ToUserIds    string    // 需要@的用户id列表，多个用户用，隔开
	Type         int       // 消息类型
	Content      []byte    // 消息内容
	Seq          int64     // 消息同步序列
	SendTime     time.Time // 消息发送时间
	Status       int32     // 创建时间
}

func FormatUserIds(userId []int64) string {
	build := strings.Builder{}
	for i, v := range userId {
		build.WriteString(strconv.FormatInt(v, 10))
		if i != len(userId)-1 {
			build.WriteString(",")
		}
	}
	return build.String()
}

func UnformatUserIds(userIdStr string) []int64 {
	if userIdStr == "" {
		return []int64{}
	}
	toUserIdStrs := strings.Split(userIdStr, ",")
	toUserIds := make([]int64, 0, len(toUserIdStrs))
	for i := range toUserIdStrs {
		userId, err := strconv.ParseInt(toUserIdStrs[i], 10, 64)
		if err != nil {
			logger.Sugar.Error(err)
			continue
		}
		toUserIds = append(toUserIds, userId)
	}
	return toUserIds
}

func MessageToPB(message *Message) *pb.Message {
	return &pb.Message{
		Sender: &pb.Sender{
			SenderType: pb.SenderType(message.SenderType),
			SenderId:   message.SenderId,
		},
		ReceiverType:   pb.ReceiverType(message.ReceiverType),
		ReceiverId:     message.ReceiverId,
		ToUserIds:      UnformatUserIds(message.ToUserIds),
		MessageType:    pb.MessageType(message.Type),
		MessageContent: message.Content,
		Seq:            message.Seq,
		SendTime:       util.UnixMilliTime(message.SendTime),
		Status:         pb.MessageStatus(message.Status),
	}
}

func MessagesToPB(messages []Message) []*pb.Message {
	pbMessages := make([]*pb.Message, 0, len(messages))
	for i := range messages {
		pbMessage := MessageToPB(&messages[i])
		if pbMessages != nil {
			pbMessages = append(pbMessages, pbMessage)
		}
	}
	return pbMessages
}
