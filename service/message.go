package service

import (
	"errors"
	"fmt"
	"net"
	"project/dao/mongo"
	"project/middleware/kafka"
	"project/models"
	"sync"
	"sync/atomic"
	"time"
)

var chatConnMap = sync.Map{}
var messageIdSequence = int64(1)

func RunMessageServer() {
	//listen, err := net.Listen("tcp", "127.0.0.1:9090")
	//if err != nil {
	//	fmt.Printf("Run message sever failed: %v\n", err)
	//	return
	//}
	//
	//for {
	//	conn, err := listen.Accept()
	//	if err != nil {
	//		fmt.Printf("Accept conn failed: %v\n", err)
	//		continue
	//	}
	//
	//	go process(conn)
	//}

}

func process(conn net.Conn) {
	//defer conn.Close()
	//
	//var buf [256]byte
	//for {
	//	n, err := conn.Read(buf[:])
	//	if n == 0 {
	//		if err == io.EOF {
	//			break
	//		}
	//		fmt.Printf("Read message failed: %v\n", err)
	//		continue
	//	}
	//
	//	var event = models.MessageSendEvent{}
	//	_ = json.Unmarshal(buf[:n], &event)
	//	fmt.Printf("Receive Message：%+v\n", event)
	//
	//	fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
	//	if len(event.MsgContent) == 0 {
	//		chatConnMap.Store(fromChatKey, conn)
	//		continue
	//	}
	//
	//	toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
	//	writeConn, exist := chatConnMap.Load(toChatKey)
	//	if !exist {
	//		fmt.Printf("User %d offline\n", event.ToUserId)
	//		continue
	//	}
	//
	//	pushEvent := models.MessagePushEvent{
	//		FromUserId: event.UserId,
	//		MsgContent: event.MsgContent,
	//	}
	//	pushData, _ := json.Marshal(pushEvent)
	//	_, err = writeConn.(net.Conn).Write(pushData)
	//	if err != nil {
	//		fmt.Printf("Push message failed: %v\n", err)
	//	}
	//}
}

func SendMessage(fromUserId, toUserId uint, content string) (err error) {

	// TODO 1、非好友关系
	if !IsInMyFollowList(fromUserId,toUserId) || !IsInMyFollowList(toUserId,fromUserId) {
		return errors.New("对方不是您的好友，无法发送消息。")
	}

	// TODO 2、非法action_type
	if false {
		return errors.New("无效的请求：错误的action_type")
	}

	// TODO 3、对数据进行非对称加密
	messageData := &models.Message{
		ID:         uint(atomic.AddInt64(&messageIdSequence, 1)),
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Content:    content,
		CreateTime: time.Now().Unix(),
	}

	fmt.Println("Message Data=====")
	fmt.Println(messageData)
	// 往mongo发送聊天记录
	//err = mongo.SendMessage(messageData)
	//if err != nil {
	//	log.Println("mongo.SendMessage err:", err)
	//	return err
	//}

	// 聊天记录发向kafka
	go kafka.Produce(messageData)

	return
}

func GetMessageList(fromUserId, toUserId uint, preMsgTime int64) ([]models.Message, error) {
	msgList, err := mongo.GetMessageList(fromUserId, toUserId, preMsgTime)
	if err != nil {
		fmt.Println("mongo.GetMessageList err:", err)
		return nil, err
	}
	return msgList, nil
}

//func genChatKey(userIdA uint, userIdB uint) string {
//	if userIdA > userIdB {
//		return fmt.Sprintf("%d_%d", userIdB, userIdA)
//	}
//	return fmt.Sprintf("%d_%d", userIdA, userIdB)
//}
