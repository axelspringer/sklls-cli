package analyzer

/*
	Two types of dependency parsers:
	(1) Line by line
	    --> e.g. Node.js / Typescript / JS: require('<dependency name>') and .... from "<dependency name>" are both recognizable from a single line
	(2) Lookahead
		--> e.g. go: After an "import" statement, all the dependencies are either in the same line (e.g. import "<dependency name>") or in multiple lines - in the latter case an "import" statement needs to be recognized, followed by recording the following lines until a closing ) is detected.
*/

type DependencyParser interface {
	GetDisplayName() string
	GetSupportedExt() []string
	ExtSupported(fileExtension string) bool
	ParseLine(line string) (string, string)
}
