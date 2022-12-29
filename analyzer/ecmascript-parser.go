package analyzer

// TODO: Add infos about amount of testing somebody wrote (and what test libraries they used!)
type EcmaScriptParser struct {
	DisplayName        string
	Extensions         map[string]bool
	EcmaScriptPackages *EcmaScriptPackages

	SingleLineMatcher EcmascriptSingleLineImportMatcher
	RequireMatcher    EcmascriptRequireMatcher
}

func NewEcmaScriptParser(ecmaScriptPackages *EcmaScriptPackages) *EcmaScriptParser {
	displayName := "EcmaScript"
	extensions := map[string]bool{
		".js":  true,
		".ts":  true,
		".jsx": true,
		".tsx": true,
		".vue": true,
	}

	singleLineMatcher := NewEcmascriptSingleLineImportMatcher()
	requireMatcher := NewEcmascriptRequireMatcher()

	parser := EcmaScriptParser{
		DisplayName:        displayName,
		Extensions:         extensions,
		EcmaScriptPackages: ecmaScriptPackages,
		SingleLineMatcher:  *singleLineMatcher,
		RequireMatcher:     *requireMatcher,
	}

	return &parser
}

func (parser *EcmaScriptParser) GetDisplayName() string {
	return parser.DisplayName
}

func (parser *EcmaScriptParser) GetSupportedExt() []string {
	supportedExt := []string{}
	for ext, _ := range parser.Extensions {
		supportedExt = append(supportedExt, ext)
	}

	return supportedExt
}

func (parser *EcmaScriptParser) ExtSupported(fileExtension string) bool {
	ok, _ := parser.Extensions[fileExtension]
	return ok
}

// TODO - Optional: Test out how much faster this would be with checking for the "from" statement first:
// https://pkg.go.dev/github.com/boyter/go-string#IndexAll
func (parser *EcmaScriptParser) ParseLine(line string) (string, string) {
	dep, _ := parser.SingleLineMatcher.ParseLine(line)
	if dep != "" {
		if found, dependency, version := parser.EcmaScriptPackages.GetDepFromPackages(Dependency(dep)); found {
			return string(dependency), string(version)
		}
	}

	dep, _ = parser.RequireMatcher.ParseLine(line)
	if dep != "" {
		if found, dependency, version := parser.EcmaScriptPackages.GetDepFromPackages(Dependency(dep)); found {
			return string(dependency), string(version)
		}
	}

	return "", ""
}
