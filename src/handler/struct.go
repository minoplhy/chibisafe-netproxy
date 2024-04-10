package handler

type UploadPostMeta struct {
	Name        string `json:"name"`
	FileSize    int64  `json:"size"`
	ContentType string `json:"contentType"`
}

type UploadResponseMeta struct {
	URL        string `json:"url"`
	Identifier string `json:"identifier"`
	PublicURL  string `json:"publicUrl"`
}

type UploadProcessMeta struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	ContentType string `json:"type"`
}

type UploadProcessResponseMeta struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}

type URLRequest struct {
	URL           string
	ContentType   string
	Method        string
	ContentLength *int
	Header        map[string]string
}
