package api


type CodeRequest struct {
	Language     string  `json:"language"`
    Content  string  `json:"content"`
}

type ExtentedCodeRequest struct {
	UserID uint
	Language string
	Content string
	ReqTS int64
	UniqueID string

}