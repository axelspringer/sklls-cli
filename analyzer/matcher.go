package analyzer

type ParserResponse int

const (
	NoDependencyFound ParserResponse = iota
	DependencyFound
	StartOfMultilineDependencyFound
	EndOfMultilineDependencyFound
)

type Matcher interface {
	// The parser response is just an FYI - the matchers are expected to keep track of
	// multiline import states themselves!
	ParseLine(line string) (string, ParserResponse)
}
