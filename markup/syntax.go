package markup

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

const tokenSpace = "\u00A0\u200D\u200B"

var prismLanguages = map[string]bool{
	"extend": true,
	"plain":  true, "plaintext": true, "text": true, "txt": true, "markup": true, "html": true, "mathml": true,
	"svg": true, "xml": true, "ssml": true, "atom": true, "rss": true, "css": true, "clike": true,
	"javascript": true, "js": true, "typescript": true, "ts": true, "jsx": true, "tsx": true, "python": true,
	"py": true, "java": true, "c": true, "cpp": true, "csharp": true, "cs": true, "dotnet": true, "php": true,
	"ruby": true, "rb": true, "go": true, "rust": true, "sql": true, "bash": true, "sh": true, "shell": true,
	"json": true, "webmanifest": true, "yaml": true, "yml": true, "markdown": true, "md": true,
}

var languageAliases = map[string]string{"jsx": "javascript", "tsx": "typescript"}

var chromaLexerNames = map[string]string{
	"markup": "html", "mathml": "xml", "svg": "xml", "ssml": "xml", "atom": "xml", "rss": "xml",
	"dotnet": "csharp", "webmanifest": "json", "shell": "bash",
	"plain": "", "plaintext": "", "text": "", "txt": "", "clike": "",
}

var chromaTokenNames = map[chroma.TokenType]string{
	chroma.Comment: "comment", chroma.Keyword: "keyword", chroma.KeywordConstant: "boolean",
	chroma.KeywordType: "class-name", chroma.LiteralString: "string", chroma.LiteralStringRegex: "regex",
	chroma.LiteralStringSymbol: "symbol", chroma.LiteralNumber: "number", chroma.NameFunction: "function",
	chroma.NameClass: "class-name", chroma.NameConstant: "constant", chroma.NameBuiltin: "builtin",
	chroma.NameVariable: "variable", chroma.NameTag: "tag", chroma.NameAttribute: "attr-name",
	chroma.NameDecorator: "function", chroma.NameNamespace: "namespace", chroma.NameEntity: "entity",
	chroma.NameLabel: "symbol", chroma.NameProperty: "property", chroma.Operator: "operator",
	chroma.OperatorWord: "keyword", chroma.Punctuation: "punctuation", chroma.GenericDeleted: "deleted",
	chroma.GenericInserted: "inserted", chroma.GenericEmph: "italic", chroma.GenericStrong: "bold",
	chroma.GenericHeading: "important", chroma.GenericSubheading: "important",
}

func processCodeBlock(attributes map[string]string, code string, budget *renderBudget) string {
	requested := strings.ToLower(attributes["language"])
	if requested == "" {
		requested = "javascript"
	}
	language := requested
	if alias, ok := languageAliases[requested]; ok {
		language = alias
	}
	themeName := attributes["theme"]
	if themeName == "" {
		themeName = "oneDark"
	}
	theme := resolveCodeTheme(themeName)
	customStyle := orderedStyle(extractAllStyleAttributes(attributes))

	if code == "" {
		return "<pre><code></code></pre>"
	}

	length := len(utf16.Encode([]rune(code)))
	if !prismLanguages[language] || length > maxCodeblockLength || budget.highlighted+length > maxTotalCodeblockLength {
		return fmt.Sprintf(
			`<pre style="%s"><code>%s</code></pre>`,
			preStyle(theme.base, []styleEntry{{"width", "100%"}}, customStyle),
			escapeCode(code),
		)
	}
	budget.highlighted += length

	inherited := map[string]string{}
	if fontFamily := attributes["font-family"]; fontFamily != "" {
		inherited["font-family"] = fontFamily
	}
	lines := highlightLines(code, language, theme, inherited)
	numbered := attributes["line-numbers"] == "true"
	var html strings.Builder
	for index, line := range lines {
		number := ""
		if numbered {
			number = fmt.Sprintf(`<span style="width:2em;display:inline-block">%d</span>`, index+1)
		}
		html.WriteString(fmt.Sprintf(`<p style="margin:0;min-height:1em">%s%s</p>`, number, line))
	}

	return fmt.Sprintf(
		`<pre style="%s"><code>%s</code></pre>`,
		preStyle(theme.base, []styleEntry{{"width", "100%"}, {"box-sizing", "border-box"}}, customStyle),
		html.String(),
	)
}

func orderedStyle(style map[string]string) []styleEntry {
	keys := make([]string, 0, len(style))
	for key := range style {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return stylePriority[keys[i]] < stylePriority[keys[j]]
	})
	entries := make([]styleEntry, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, styleEntry{key, style[key]})
	}
	return entries
}

func preStyle(layers ...[]styleEntry) string {
	var merged []styleEntry
	for _, layer := range layers {
		for _, entry := range layer {
			replaced := false
			for index := range merged {
				if merged[index].key == entry.key {
					merged[index].value = entry.value
					replaced = true
					break
				}
			}
			if !replaced {
				merged = append(merged, entry)
			}
		}
	}

	parts := make([]string, 0, len(merged))
	for _, entry := range merged {
		parts = append(parts, kebabProperty(entry.key)+":"+entry.value)
	}
	return strings.Join(parts, ";")
}

func kebabProperty(key string) string {
	var property strings.Builder
	for _, char := range key {
		if char >= 'A' && char <= 'Z' {
			property.WriteByte('-')
			property.WriteRune(char + ('a' - 'A'))
		} else {
			property.WriteRune(char)
		}
	}
	if strings.HasPrefix(property.String(), "ms-") {
		return "-" + property.String()
	}
	return property.String()
}

func highlightLines(code, language string, theme codeTheme, inherited map[string]string) []string {
	type styledText struct {
		style string
		text  string
	}

	inheritedStyle := styleToString(inherited)
	tokens := []styledText{{inheritedStyle, code}}
	if lexer := codeLexer(language); lexer != nil {
		iterator, err := lexer.Tokenise(&chroma.TokeniseOptions{State: "root", EnsureLF: false}, code)
		if err == nil {
			chromaTokens := iterator.Tokens()
			if last := len(chromaTokens) - 1; last >= 0 && !strings.HasSuffix(code, "\n") {
				chromaTokens[last].Value = strings.TrimSuffix(chromaTokens[last].Value, "\n")
			}
			tokens = make([]styledText, 0, len(chromaTokens))
			for _, token := range chromaTokens {
				tokens = append(tokens, styledText{tokenStyle(theme, token.Type, inherited), token.Value})
			}
		}
	}

	lines := []string{""}
	for _, token := range tokens {
		for index, part := range strings.Split(token.text, "\n") {
			if index > 0 {
				lines = append(lines, "")
			}
			if part == "" {
				continue
			}
			escaped := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;", " ", tokenSpace).Replace(part)
			if token.style != "" {
				escaped = fmt.Sprintf(`<span style="%s">%s</span>`, token.style, escaped)
			}
			lines[len(lines)-1] += escaped
		}
	}
	return lines
}

func codeLexer(language string) chroma.Lexer {
	name, mapped := chromaLexerNames[language]
	if !mapped {
		name = language
	}
	if name == "" {
		return nil
	}
	lexer := lexers.Get(name)
	if lexer == nil {
		return nil
	}
	return chroma.Coalesce(lexer)
}

func tokenStyle(theme codeTheme, tokenType chroma.TokenType, inherited map[string]string) string {
	for current := tokenType; ; current = current.Parent() {
		if name, ok := chromaTokenNames[current]; ok {
			merged := map[string]string{}
			for key, value := range inherited {
				merged[key] = value
			}
			for key, value := range theme.tokens[name] {
				merged[key] = value
			}
			return styleToString(merged)
		}
		if current.Parent() == current || current == 0 {
			break
		}
	}
	return styleToString(inherited)
}

func escapeCode(code string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(code, "<", "&lt;"),
		">",
		"&gt;",
	)
}
