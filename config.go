package oapi_sdk_go_demo

import "github.com/larksuite/oapi-sdk-go/v3"

const ( // 正式
	AppId             = "cli_a61dccc07dfd900b"
	AppSecret         = "KJPCMvEz4gQDlJikPx8cUfAnxyi22NaC"
	EncryptKey        = "W9t8WbHuey1RFpGTv9n4xINsWRwaJQIw"
	VerificationToken = "gBLa23S6EvEO60fZVk3IGdcrwpTFCBrG"
)

//const ( // 测试
//	AppId             = "cli_a61d18391cb7d00c"
//	AppSecret         = "QfE1wMDqoVnU2dllNOMEcbvQcekLHSiK"
//	EncryptKey        = "EP05QILeWqmdmPrD3J4yxhN04JuNZqWH"
//	VerificationToken = "Dz5AEqgEDqAvCJXGY17nwnuOJoijUOPG"
//)

var Client = lark.NewClient(AppId, AppSecret)
