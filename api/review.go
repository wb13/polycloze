// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

type Review struct {
	Word    string `json:"word"`
	Correct bool   `json:"correct"`
}

type Reviews struct {
	Reviews []Review `json:"reviews"`
}
