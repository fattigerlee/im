package connect

import (
	"fmt"
	"im/pkg/logger"
	"im/pkg/util"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}

	conn := &Conn{
		CoonType: ConnTypeWS,
		WS:       wsConn,
	}
	DoConn(conn)
}

// DoConn 处理连接
func DoConn(conn *Conn) {
	defer util.RecoverPanic()

	for {
		err := conn.WS.SetReadDeadline(time.Now().Add(12 * time.Minute))
		_, data, err := conn.WS.ReadMessage()
		if err != nil {
			fmt.Println(err)
			HandleReadErr(conn, err)
			return
		}

		conn.HandleMessage(data)
	}
}

// HandleReadErr 读取conn错误
func HandleReadErr(conn *Conn, err error) {
	logger.Logger.Debug("read tcp error：", zap.Int64("user_id", conn.UserId),
		zap.Int64("device_id", conn.DeviceId), zap.Error(err))
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	conn.Close()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

func StartWSServer(address string) {
	logger.Info("websocket server start, address: %s", address)

	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
