package models

type ResponsePacket struct {
	Error   bool   `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
