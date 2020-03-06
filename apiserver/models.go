package apiserver

type genericResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

type linkModel struct {
	LinkID        string `json:"linkID"`
	CanonicalName string `json:"canonicalName"`
	URI           string `json:"URI"`
	TargetURL     string `json:"targetURL"`
	// Audit info
	Created        int64  `json:"created"`
	LastModified   int64  `json:"lastModified"`
	LastModifiedBy string `json:"lastModifiedBy"`
}
