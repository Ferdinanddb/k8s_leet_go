package api


type CodeRequest struct {
	Language     string  `json:"language"`
    Content  string  `json:"content"`
}

type ExtentedCodeRequest struct {
	CodeReq CodeRequest
	ReqTS int64
	UniqueID string

}