package utils

// 建立 200 500 的 Interface
type IH200 struct {
	Status  bool        `default:"true" json:"status"`
	Msg     string      `default:"成功" json:"msg"`
	Predict interface{} `default:"" json:"predict"`
}

type IH500 struct {
	Status bool   `default:"false" json:"status"`
	Msg    string `default:"desc..." json:"msg"`
}

// Response 的格式
type Response struct {
	Status  bool        `default:"true" json:"status"`
	Msg     string      `default:"成功" json:"msg"`
	Predict interface{} `default:"" json:"predict"`
}

// 返回200
func H200(data interface{}, msg string) Response {
	return Response{
		Status:  true,
		Msg:     msg,
		Predict: data,
	}
}

// 返回500
func H500(msg string) Response {
	return Response{
		Status: false,
		Msg:    msg,
	}
}
