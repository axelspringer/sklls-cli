package analyzer

import "regexp"

// Import matcher
type EcmascriptSingleLineImportMatcher struct {
	regex *regexp.Regexp
}

func NewEcmascriptSingleLineImportMatcher() *EcmascriptSingleLineImportMatcher {
	regex := regexp.MustCompile(`from\s*("|')(?P<import>.*)('|")(;?)\s*$`)
	matcher := &EcmascriptSingleLineImportMatcher{regex}
	return matcher
}

func (matcher *EcmascriptSingleLineImportMatcher) ParseLine(line string) (string, ParserResponse) {
	dependencies := matcher.regex.FindStringSubmatch(line)
	if len(dependencies) > 0 {
		dependency := dependencies[2]
		return dependency, DependencyFound
	}

	return "", NoDependencyFound
}

// Require matcher
type EcmascriptRequireMatcher struct {
	regex *regexp.Regexp
}

func NewEcmascriptRequireMatcher() *EcmascriptRequireMatcher {
	regex := regexp.MustCompile(`require\(('|")(.*)('|")\)`)
	matcher := &EcmascriptRequireMatcher{regex}
	return matcher
}

func (matcher *EcmascriptRequireMatcher) ParseLine(line string) (string, ParserResponse) {
	dependencies := matcher.regex.FindStringSubmatch(line)
	if len(dependencies) > 0 {
		dependency := dependencies[2]
		return dependency, DependencyFound
	}

	return "", NoDependencyFound
}
