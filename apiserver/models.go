package apiserver

type genericResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

// LinkModel is a model of saved links
type LinkModel struct {
	// LinkID        string `json:"LinkID"`
	LinkPath      string `json:"LinkPath"` // LinkPath is now the unique identifier
	CanonicalName string `json:"CanonicalName"`
	TargetURL     string `json:"TargetURL"`
	Enabled       bool   `json:"Enabled"`
	// Audit info
	CreatedTime    int64  `json:"CreatedTime"`
	LastModified   int64  `json:"LastModified"`
	LastModifiedBy string `json:"LastModifiedBy"`
}
