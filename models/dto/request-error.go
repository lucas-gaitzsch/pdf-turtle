package dto

type RequestError struct {
	Msg       string `json:"msg"`
	Err       string `json:"err"`
	RequestId string `json:"requestId"`
} // @name RequestError
