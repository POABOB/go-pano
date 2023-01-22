package utils

// 建立 200 500 的 Interface
type IH200 struct {
	Status bool        `default:"true" json:"status"`
	Msg    string      `default:"成功" json:"msg"`
	Data   interface{} `default:"" json:"data"`
}

type IH500 struct {
	Status bool   `default:"false" json:"status"`
	Msg    string `default:"desc..." json:"msg"`
}

// Response 的格式
type Response struct {
	Status bool        `default:"true" json:"status"`
	Msg    string      `default:"成功" json:"msg"`
	Data   interface{} `default:"" json:"data"`
}

// 返回200
func H200(data interface{}, msg string) Response {
	return Response{
		Status: true,
		Msg:    msg,
		Data:   data,
	}
}

// 返回500
func H500(msg string) Response {
	return Response{
		Status: false,
		Msg:    msg,
	}
}
