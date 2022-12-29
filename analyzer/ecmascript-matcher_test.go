package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcmascriptSingleLineImportMatcher(t *testing.T) {
	const importDoubleQuotes = `import React from "react"`
	const importSingleQuotes = `import React from 'react'`
	const importMultiline = `} from 'react';`
	const notAnImport = "console.log('hello!');"

	matcher := NewEcmascriptSingleLineImportMatcher()
	dep, res := matcher.ParseLine(importDoubleQuotes)
	assert.Equal(t, "react", dep, "Extracts correct dependency for double quotes")
	assert.Equal(t, DependencyFound, res, "Returns correct ParserResponse for double quotes")

	dep, res = matcher.ParseLine(importSingleQuotes)
	assert.Equal(t, "react", dep, "Extracts correct dependency for single quotes")
	assert.Equal(t, DependencyFound, res, "Returns correct ParserResponse for single quotes")

	dep, res = matcher.ParseLine(importMultiline)
	assert.Equal(t, "react", dep, "Extracts correct dependency for multiline imports")
	assert.Equal(t, DependencyFound, res, "Returns correct ParserResponse for multiline imports")

	dep, res = matcher.ParseLine(notAnImport)
	assert.Equal(t, "", dep, "Returns empty string if no import was found")
	assert.Equal(t, NoDependencyFound, res, "Returns correct parser response for no import found")
}

func TestEcmascriptRequireMatcher(t *testing.T) {
	const requireWithDoubleQuotes = `const expressLib = require("express")`
	const requireWithSingleQuotes = `const expressLib = require('express')`
	const requireStandalone = `require('./some/stylesheet.css')`
	const notAnImport = "console.log('hello!');"

	matcher := NewEcmascriptRequireMatcher()
	dep, res := matcher.ParseLine(requireWithDoubleQuotes)
	assert.Equal(t, "express", dep, "Extracts correct dependency for double quotes")
	assert.Equal(t, DependencyFound, res, "Returns correct ParserResponse for double quotes")

	dep, res = matcher.ParseLine(requireWithSingleQuotes)
	assert.Equal(t, "express", dep, "Extracts correct dependency for single quotes")
	assert.Equal(t, DependencyFound, res, "Returns correct ParserResponse for single quotes")

	dep, res = matcher.ParseLine(requireStandalone)
	assert.Equal(t, "./some/stylesheet.css", dep, "Extracts correct dependency for standalone require statement")
	assert.Equal(t, DependencyFound, res, "Returns correct ParserResponse for single quotes")

	dep, res = matcher.ParseLine(notAnImport)
	assert.Equal(t, "", dep, "Returns empty string if no import was found")
	assert.Equal(t, NoDependencyFound, res, "Returns correct parser response for no import found")
}

func BenchmarkEcmascriptSingleLineImportMatcher(b *testing.B) {
	const importDoubleQuotes = `import React from "react"`
	matcher := NewEcmascriptSingleLineImportMatcher()

	for i := 0; i < b.N; i++ {
		matcher.ParseLine(importDoubleQuotes)
	}
}
