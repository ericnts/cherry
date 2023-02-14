package exception

import (
	"context"
	"fmt"
	"github.com/ericnts/cherry/local"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	errMsg := map[string]map[string]string{
		local.ZH: {
			Ok.LocalKey():               "请求成功",
			Warn.LocalKey():             "请求成功，但有异常",
			None.LocalKey():             "无需动作",
			AuthInvalid.LocalKey():      "认证失败",
			LoginForbid.LocalKey():      "禁止登录",
			TokenInvalid.LocalKey():     "凭证失效",
			PermissionDenied.LocalKey(): "没有操作权限",
			ParamInvalid.LocalKey():     "无效的请求参数",
			DataInvalid.LocalKey():      "数据校验失败",
			DataNotFound.LocalKey():     "数据不存在",
			DataRepeat.LocalKey():       "数据重复",
			UploadFailed.LocalKey():     "文件上传失败",
			LicenseInvalid.LocalKey():   "服务证书验证失败",
			Timeout.LocalKey():          "服务响应超时",
			ExecFail.LocalKey():         "执行失败",
			TokenReplaced.LocalKey():    "用户在其他终端登录",
			Private.LocalKey():          "服务内部错误",
			IDBlank.LocalKey():          "ID不能为空",
			OfficeIDBlank.LocalKey():    "机构ID不能为空",
			VOInvalid.LocalKey():        "没有定义列表实体类",
		},
		local.ZhHk: {
			Ok.LocalKey():               "請求成功",
			Warn.LocalKey():             "請求成功，但有異常",
			None.LocalKey():             "无需动作",
			AuthInvalid.LocalKey():      "認證失敗",
			LoginForbid.LocalKey():      "禁止登錄",
			TokenInvalid.LocalKey():     "憑證失效",
			PermissionDenied.LocalKey(): "沒有操作許可權",
			ParamInvalid.LocalKey():     "無效的請求參數",
			DataInvalid.LocalKey():      "資料校驗失敗",
			DataNotFound.LocalKey():     "資料不存在",
			DataRepeat.LocalKey():       "資料重複",
			UploadFailed.LocalKey():     "檔上傳失敗",
			LicenseInvalid.LocalKey():   "服务证书验证失败",
			Timeout.LocalKey():          "服務回應超時",
			ExecFail.LocalKey():         "執行失敗",
			TokenReplaced.LocalKey():    "使用者在其他終端登錄",
			Private.LocalKey():          "服務內部錯誤",
			IDBlank.LocalKey():          "ID不能為空",
			OfficeIDBlank.LocalKey():    "机构ID不能為空",
			VOInvalid.LocalKey():        "沒有定義列表實體類",
		},
		local.EN: {
			Ok.LocalKey():               "Successful",
			Warn.LocalKey():             "Warning",
			None.LocalKey():             "",
			AuthInvalid.LocalKey():      "Verification failure",
			LoginForbid.LocalKey():      "Login prohibited",
			TokenInvalid.LocalKey():     "Token invalid",
			PermissionDenied.LocalKey(): "No permission for operation",
			ParamInvalid.LocalKey():     "Invalid Input",
			DataInvalid.LocalKey():      "Information calibration failure",
			DataNotFound.LocalKey():     "Information not exist",
			DataRepeat.LocalKey():       "Information duplicated",
			UploadFailed.LocalKey():     "File upload failure",
			LicenseInvalid.LocalKey():   "License invalid",
			Timeout.LocalKey():          "Service response over-time",
			ExecFail.LocalKey():         "Execution failure",
			TokenReplaced.LocalKey():    "User logined in other device",
			Private.LocalKey():          "Internal service error",
			IDBlank.LocalKey():          "ID cannot be empty",
			OfficeIDBlank.LocalKey():    "Department ID cannot be empty",
			VOInvalid.LocalKey():        "VO invalid",
		},
	}
	local.Append(errMsg)
}

type Error int

const (
	Ok               Error = 0   //请求成功
	Warn             Error = 1   //请求成功，但有异常
	None             Error = 2   //无需动作
	AuthInvalid      Error = 100 //认证失败
	LoginForbid      Error = 101 //禁止登录
	TokenInvalid     Error = 102 //Token失效
	PermissionDenied Error = 103 //没有操作权限
	ParamInvalid     Error = 104 //无效的请求参数
	DataInvalid      Error = 105 //数据校验失败
	DataNotFound     Error = 106 //数据不存在
	DataRepeat       Error = 107 //数据重复
	UploadFailed     Error = 108 //文件上传失败
	LicenseInvalid   Error = 109 //服务证书验证失败
	Timeout          Error = 110 //服务响应超时
	ExecFail         Error = 111 //执行失败
	TokenReplaced    Error = 112 //用户在其他终端登录
	Private          Error = 500 //服务内部错误
)

const (
	IDBlank       Error = iota + 150 //ID不能为空
	OfficeIDBlank                    //机构IDb不能为空

	VOInvalid Error = 501 //没有定义列表实体类
)

func (e Error) Int() int {
	return int(e)
}

func (e Error) Error() string {
	return local.Translate(context.TODO(), e.LocalKey())
}

func (e Error) Attachments() []error {
	return nil
}

func (e Error) GRPCStatus() *status.Status {
	return status.New(e.Convert(), e.Error())
}

func (e Error) Convert() codes.Code {
	code := codes.Unknown
	switch e {
	case Ok, Warn, None:
		code = codes.OK
	case AuthInvalid, LoginForbid, TokenInvalid, TokenReplaced:
		code = codes.Unauthenticated
	case PermissionDenied:
		code = codes.PermissionDenied
	case ParamInvalid:
		code = codes.InvalidArgument
	case DataInvalid:
		code = codes.Unavailable
	case DataNotFound:
		code = codes.NotFound
	case DataRepeat:
		code = codes.AlreadyExists
	case UploadFailed, ExecFail:
		code = codes.Aborted
	case Timeout:
		code = codes.DeadlineExceeded
	case Private:
		code = codes.Internal
	}
	return code
}

func ConvertCode(code codes.Code) Error {
	switch code {
	case codes.OK:
		return Ok
	case codes.Unauthenticated:
		return AuthInvalid
	case codes.PermissionDenied:
		return PermissionDenied
	case codes.InvalidArgument:
		return ParamInvalid
	case codes.Unavailable:
		return DataInvalid
	case codes.NotFound:
		return DataNotFound
	case codes.AlreadyExists:
		return DataRepeat
	case codes.Aborted:
		return ExecFail
	case codes.DeadlineExceeded:
		return Timeout
	case codes.Internal:
		return Private
	}
	return Private
}

func (e Error) LocalKey() string {
	return fmt.Sprintf("Error_%d", e.Int())
}

func IsDataNotFound(err error) bool {
	ce, ok := err.(CustomError)
	if ok && ce.ErrCode == DataNotFound.Int() {
		return true
	}
	return false
}
