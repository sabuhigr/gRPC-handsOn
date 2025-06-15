package types

type ErrDetails struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
