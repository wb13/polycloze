package api

type Review struct {
	Word    string `json:"word"`
	Correct bool   `json:"correct"`
}

type Reviews struct {
	Reviews []Review `json:"reviews"`
}
