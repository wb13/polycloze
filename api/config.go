package api

type Config struct {
	ReviewDb      string
	Lang1Db       string
	Lang2Db       string
	TranslationDb string

	AllowCORS bool
}
