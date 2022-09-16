package rpc

type Ctx struct {
}

type (
	BaseRsp struct {
		Error
		Hint string `json:"hint,omitempty"`
		Data any    `json:"data,omitempty"`
	}
)
