package models

type (
	SuccessfulResponse struct {
		APIVersion string                 `json:"api_version"`
		Context    string                 `json:"context"`
		Method     string                 `json:"method"`
		Params     map[string]interface{} `json:"params,omitempty"`
		Data       interface{}            `json:"data,omitempty"`
	}

	Error struct {
		Code    uint16                   `json:"code"`
		Message string                   `json:"message"`
		Errors  []map[string]interface{} `json:"errors"`
	}

	BadResponse struct {
		APIVersion string                 `json:"api_version"`
		Context    string                 `json:"context"`
		Method     string                 `json:"method"`
		Params     map[string]interface{} `json:"params,omitempty"`
		Error      `json:"error"`
	}
)
