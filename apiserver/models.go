package apiserver

type genericResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

// LinkModel is a model of saved links
type LinkModel struct {
	LinkID        string `json:"linkID"`
	CanonicalName string `json:"canonicalName"`
	LinkPath      string `json:"URI"`
	TargetURL     string `json:"targetURL"`
	// Audit info
	CreatedTime    int64  `json:"createdTime"`
	LastModified   int64  `json:"lastModified"`
	LastModifiedBy string `json:"lastModifiedBy"`
}
