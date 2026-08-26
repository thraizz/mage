package bot

import (
	"encoding/json"
	"regexp"
	"strings"
)

// json5_test.go is a Go port of reference/json5_utils.py (upstream
// src/magebench/common/json5_utils.py), used to write golden files.
//
// Why not plain JSON: a golden file whose entire value is one 4KB string with
// every newline escaped as \n is unreviewable. Any change to the board renderer
// shows up as a single modified line in the diff and a human cannot tell a
// typo from a redesign. JSON5 line continuations put each rendered board line
// on its own file line, so `git diff` on a golden shows exactly which board
// lines moved.

var trailingCommaRe = regexp.MustCompile(`([^\s,\[\{])\n(\s*[\]\}])`)

// dumpsJSON5 serialises obj with 2-space indent, sorted keys (encoding/json
// sorts map keys for us, matching python's sort_keys=True), trailing commas,
// and \n inside strings expanded into line continuations.
// Port of json5_utils.dumps_json5.
func dumpsJSON5(obj any) (string, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	// Python's ensure_ascii=False: keep non-ASCII literal. Go's default HTML
	// escaping would mangle <, > and & inside oracle text, so turn it off.
	enc.SetEscapeHTML(false)
	if err := enc.Encode(obj); err != nil {
		return "", err
	}
	text := strings.TrimRight(buf.String(), "\n")
	text = addTrailingCommas(text)
	return expandMultilineStrings(text), nil
}

// addTrailingCommas ports json5_utils._add_trailing_commas.
func addTrailingCommas(text string) string {
	return trailingCommaRe.ReplaceAllString(text, "$1,\n$2")
}

// expandMultilineStrings ports json5_utils._expand_multiline_strings: walk the
// text tracking string context so that only \n inside a string (and not the
// \\n that means a literal backslash followed by n) becomes a continuation.
func expandMultilineStrings(text string) string {
	var out strings.Builder
	inString := false
	for i := 0; i < len(text); {
		ch := text[i]
		if inString {
			if ch == '\\' {
				if i+1 < len(text) {
					next := text[i+1]
					if next == 'n' {
						out.WriteString("\\n\\\n")
						i += 2
						continue
					}
					out.WriteByte(ch)
					out.WriteByte(next)
					i += 2
					continue
				}
			} else if ch == '"' {
				inString = false
			}
		} else if ch == '"' {
			inString = true
		}
		out.WriteByte(ch)
		i++
	}
	return out.String()
}
