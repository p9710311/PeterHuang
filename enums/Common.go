package enums

type JsonResultCode int

const (
	JRCodeSucc JsonResultCode = iota
	JRCodeFailed
	JRCode302 = 302 //跳轉至地址
	JRCode401 = 401 //未授權訪問
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)
