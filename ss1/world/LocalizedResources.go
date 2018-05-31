package world

import (
	"strings"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// LocalizedResources associates a language with a resource provider under a specific identifier.
type LocalizedResources struct {
	// ID is the identifier of the provider. This could be a filename for instance.
	ID string
	// Language specifies for which language the provider has resources.
	Language Language
	// Provider is the actual container of the resources.
	Provider resource.Provider
}

type languageSpecificFilenames struct {
	cybstrng string
	mfdart   string
	citalog  string
	citbark  string
}

func (spec languageSpecificFilenames) hasFilename(filename string) bool {
	return spec.cybstrng == filename
}

var localizedFilenames = map[Language]languageSpecificFilenames{
	LangDefault: {
		cybstrng: "cybstrng.res",
		mfdart:   "mfdart.res",
		citalog:  "citalog.res",
		citbark:  "citbark.res",
	},
	LangFrench: {
		cybstrng: "frnstrng.res",
		mfdart:   "mfdfrn.res",
		citalog:  "frnalog.res",
		citbark:  "frnbark.res",
	},
	LangGerman: {
		cybstrng: "gerstrng.res",
		mfdart:   "mfdger.res",
		citalog:  "geralog.res",
		citbark:  "gerbark.res",
	},
}

// LocalizeResources creates an instance of LocalizedResources based on a filename and a provider.
// The given filename is taken as an ID, as well as a hint to identify the language.
func LocalizeResourcesByFilename(provider resource.Provider, filename string) (res LocalizedResources) {
	res.ID = filename
	res.Provider = provider
	res.Language = LangAny
	for lang, loc := range localizedFilenames {
		if loc.hasFilename(strings.ToLower(filename)) {
			res.Language = lang
		}
	}

	return
}
