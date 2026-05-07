package schemas

type (
	SearchIndexReq struct {
		Numbers []int `json:"number"`
		Target  int   `json:"target"`
	}

	SearchIndexResp struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}
)
