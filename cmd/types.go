package cmd

type CarryOverTodos struct {
	Enabled                   bool
	OnlyIncomplete            bool
	HeadingRegex              string
	HeadingRegexCaseSensitive bool
}

type Config struct {
	RootPath       string `validate:"required"`
	KeepDays       int
	CarryOverTodos CarryOverTodos
	JournalsPath   string
	PagesPath      string
}
