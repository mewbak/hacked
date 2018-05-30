package world

// Language defines the human language of a resource.
type Language byte

const (
	// LangAny identifies language agnostic from.
	LangAny Language = 0xFF
	// LangDefault identifies the default language, typically English - unless modded.
	LangDefault Language = 0
	// LangFrench identifies the French language.
	LangFrench Language = 1
	// LangGerman identifies the German language.
	LangGerman Language = 2
)

// Languages returns a slice of all human languages. Does not include "Any" selector.
func Languages() []Language {
	return []Language{LangDefault, LangFrench, LangGerman}
}