package types

const (
	Static_token = "dwwf=Adowwlfpfwpifwpwpw27O!"
)

type ErrDetails struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
