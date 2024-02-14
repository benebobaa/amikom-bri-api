package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestNoWhiteSpace(t *testing.T) {
	word := "GoLanguage"

	whitespace := regexp.MustCompile(`\s`).MatchString(word)
	assert.False(t, whitespace)
}
