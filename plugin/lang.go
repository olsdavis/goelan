package plugin

type Language string

func (lang Language) IsSupported() bool {
	return lang != UnsupportedLanguage
}

const (
	PythonLanguage      Language = "Python"
	JavaScriptLanguage  Language = "JavaScript"
	UnsupportedLanguage Language = "UnsupportedLanguage"
)

// IsSupportedLanguage returns true if the language associated to the file extension
// is supported by goelan.
func IsSupportedLanguage(ext string) bool {
	return AssociatedLanguage(ext).IsSupported()
}

// AssociatedLanguage returns the language associated to the file extension.
// If it is not supported by goelan, returns UnsupportedLanguage.
func AssociatedLanguage(ext string) Language {
	switch ext {
	case "py":
		return PythonLanguage
	case "js":
		return JavaScriptLanguage
	default:
		return UnsupportedLanguage
	}
}
