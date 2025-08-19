package utils

type CarryOverTodos struct {
	Enabled                   bool
	OnlyIncomplete            bool
	HeadingRegex              string
	HeadingRegexCaseSensitive bool
}

type BaselineConfig struct {
	RootPath string `validate:"required"`
}

type Config struct {
	RootPath       string `validate:"required"`
	KeepDays       int
	KeepMonths     int
	CarryOverTodos CarryOverTodos
	JournalsPath   string
	PagesPath      string
}
