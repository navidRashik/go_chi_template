package request

type StatusType int

const AuthHeader = "Authorization"
const ContentTypeHeader = "Content-Type"
const ContentTypeValue = "application/json"
const ApiKey = "x-api-key"

const (
	BuildRequestError StatusType = iota
	NetworkError
	InformationalError
	RequestSuccess
	RedirectError
	BadRequestError
	ServerError
	UnknownError
)

func GetStatusType(status int) StatusType {
	switch {
	case status == -1 || status == -2:
		return BuildRequestError
	case status == -3:
		return NetworkError
	case 100 <= status && status <= 199:
		return InformationalError
	case 200 <= status && status <= 299:
		return RequestSuccess
	case 300 <= status && status <= 399:
		return RedirectError
	case 400 <= status && status <= 499:
		return BadRequestError
	case 500 <= status && status <= 599:
		return ServerError
	default:
		return UnknownError
	}
}
