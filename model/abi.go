package model

type Abi struct {
	Anonymous bool       `json:"anonymous"`
	Inputs    []AbiInput `json:"inputs"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
}

type AbiInput struct {
	Indexed      bool   `json:"indexed"`
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}
