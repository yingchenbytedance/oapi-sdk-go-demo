package main

import (
	"context"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"net/http"
	oapi_sdk_go_demo "oapi-sdk-go-demo"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	// 创建长连接
	wg.Add(1)
	go func() {
		CreateLongLink()
		wg.Done()
	}()
	// 创建告警群并拉人入群
	chatId, err := CreateAlertChat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("chatId: " + chatId)

	// 发送告警通知
	err = SendAlertMessage(chatId)
	if err != nil {
		fmt.Println(err)
		return
	}

	//  注册事件回调
	eventHandler := dispatcher.NewEventDispatcher(oapi_sdk_go_demo.VerificationToken, oapi_sdk_go_demo.EncryptKey)
	eventHandler.OnP2MessageReceiveV1(DoP2ImMessageReceiveV1)

	// 注册卡片回调
	cardHandler := larkcard.NewCardActionHandler(oapi_sdk_go_demo.VerificationToken, oapi_sdk_go_demo.EncryptKey, DoInteractiveCard)

	http.HandleFunc("/event", httpserverext.NewEventHandlerFunc(eventHandler,
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))
	http.HandleFunc("/card", httpserverext.NewCardActionHandlerFunc(cardHandler,
		larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

func CreateLongLink() {
	// 注册事件回调，OnP2MessageReceiveV1 为接收消息 v2.0；OnCustomizedEvent 内的 message 为接收消息 v1.0。
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			fmt.Printf("[ OnP2MessageReceiveV1 access ], data: %s\n", larkcore.Prettify(event))
			return nil
		}).
		OnCustomizedEvent("message", func(ctx context.Context, event *larkevent.EventReq) error {
			fmt.Printf("[ OnCustomizedEvent access ], type: message, data: %s\n", string(event.Body))
			return nil
		})
	// 创建Client
	cli := larkws.NewClient(oapi_sdk_go_demo.AppId, oapi_sdk_go_demo.AppSecret, // 统一
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)
	// 启动客户端
	err := cli.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
