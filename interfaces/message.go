package interfaces

type ErrorResponse struct{
	Error MessageResponse
}

type SuccessResponse struct {
	Response MessageResponse
}

type MessageResponse struct{
	Message string
}