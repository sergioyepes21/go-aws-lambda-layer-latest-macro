package macromodels

type Event struct {
	AccountId string           `json:"accountId"`
	Region    string           `json:"region"`
	RequestId string           `json:"requestId"`
	Params    MacroEventParams `json:"params"`
}
