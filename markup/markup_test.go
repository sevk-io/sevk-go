package markup

import (
	"strings"
	"testing"
)

// ============================================
// isTruthy TESTS
// ============================================

func TestIsTruthy_NilIsFalsy(t *testing.T) {
	if isTruthy(nil) {
		t.Error("expected nil to be falsy")
	}
}

func TestIsTruthy_FalseBoolIsFalsy(t *testing.T) {
	if isTruthy(false) {
		t.Error("expected false to be falsy")
	}
}

func TestIsTruthy_TrueBoolIsTruthy(t *testing.T) {
	if !isTruthy(true) {
		t.Error("expected true to be truthy")
	}
}

func TestIsTruthy_EmptyStringIsFalsy(t *testing.T) {
	if isTruthy("") {
		t.Error("expected empty string to be falsy")
	}
}

func TestIsTruthy_NonEmptyStringIsTruthy(t *testing.T) {
	if !isTruthy("hello") {
		t.Error("expected non-empty string to be truthy")
	}
}

func TestIsTruthy_ZeroFloat64IsFalsy(t *testing.T) {
	if isTruthy(float64(0)) {
		t.Error("expected 0.0 to be falsy")
	}
}

func TestIsTruthy_NonZeroFloat64IsTruthy(t *testing.T) {
	if !isTruthy(float64(42)) {
		t.Error("expected 42.0 to be truthy")
	}
}

func TestIsTruthy_ZeroIntIsFalsy(t *testing.T) {
	if isTruthy(0) {
		t.Error("expected 0 to be falsy")
	}
}

func TestIsTruthy_NonZeroIntIsTruthy(t *testing.T) {
	if !isTruthy(7) {
		t.Error("expected 7 to be truthy")
	}
}

func TestIsTruthy_EmptySliceIsFalsy(t *testing.T) {
	if isTruthy([]interface{}{}) {
		t.Error("expected empty slice to be falsy")
	}
}

func TestIsTruthy_NonEmptySliceIsTruthy(t *testing.T) {
	if !isTruthy([]interface{}{"a"}) {
		t.Error("expected non-empty slice to be truthy")
	}
}

func TestIsTruthy_MapIsTruthy(t *testing.T) {
	// map falls into default case, always truthy
	if !isTruthy(map[string]string{"a": "b"}) {
		t.Error("expected map to be truthy")
	}
}

// ============================================
// renderTemplate TESTS
// ============================================

func TestRenderTemplate_SimpleVariable(t *testing.T) {
	result := renderTemplate("Hello {%name%}!", map[string]interface{}{
		"name": "World",
	})
	if result != "Hello World!" {
		t.Errorf("expected 'Hello World!', got '%s'", result)
	}
}

func TestRenderTemplate_MissingVariable(t *testing.T) {
	result := renderTemplate("Hello {%name%}!", map[string]interface{}{})
	if result != "Hello !" {
		t.Errorf("expected 'Hello !', got '%s'", result)
	}
}

func TestRenderTemplate_VariableWithFallback(t *testing.T) {
	result := renderTemplate("Hello {%name ?? Guest%}!", map[string]interface{}{})
	if result != "Hello Guest!" {
		t.Errorf("expected 'Hello Guest!', got '%s'", result)
	}
}

func TestRenderTemplate_VariableWithFallbackPresentValue(t *testing.T) {
	result := renderTemplate("Hello {%name ?? Guest%}!", map[string]interface{}{
		"name": "Alice",
	})
	if result != "Hello Alice!" {
		t.Errorf("expected 'Hello Alice!', got '%s'", result)
	}
}

func TestRenderTemplate_EachLoop(t *testing.T) {
	result := renderTemplate("{%#each items as item%}<li>{%item.name%}</li>{%/each%}", map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "Apple"},
			map[string]interface{}{"name": "Banana"},
		},
	})
	if result != "<li>Apple</li><li>Banana</li>" {
		t.Errorf("unexpected result: '%s'", result)
	}
}

func TestRenderTemplate_EachLoopDefaultAlias(t *testing.T) {
	result := renderTemplate("{%#each items%}<li>{%item.name%}</li>{%/each%}", map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "X"},
		},
	})
	if result != "<li>X</li>" {
		t.Errorf("unexpected result: '%s'", result)
	}
}

func TestRenderTemplate_EachLoopEmptyArray(t *testing.T) {
	result := renderTemplate("{%#each items as item%}<li>{%item.name%}</li>{%/each%}", map[string]interface{}{
		"items": []interface{}{},
	})
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestRenderTemplate_EachLoopMissingArray(t *testing.T) {
	result := renderTemplate("{%#each items as item%}<li>{%item.name%}</li>{%/each%}", map[string]interface{}{})
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestRenderTemplate_EachLoopWithFallback(t *testing.T) {
	result := renderTemplate("{%#each items as item%}{%item.url ?? #%}{%/each%}", map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"url": "https://example.com"},
			map[string]interface{}{},
		},
	})
	if result != "https://example.com#" {
		t.Errorf("unexpected result: '%s'", result)
	}
}

func TestRenderTemplate_IfTrue(t *testing.T) {
	result := renderTemplate("{%#if show%}visible{%/if%}", map[string]interface{}{
		"show": true,
	})
	if result != "visible" {
		t.Errorf("expected 'visible', got '%s'", result)
	}
}

func TestRenderTemplate_IfFalse(t *testing.T) {
	result := renderTemplate("{%#if show%}visible{%/if%}", map[string]interface{}{
		"show": false,
	})
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestRenderTemplate_IfMissing(t *testing.T) {
	result := renderTemplate("{%#if show%}visible{%/if%}", map[string]interface{}{})
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestRenderTemplate_IfElseTrue(t *testing.T) {
	result := renderTemplate("{%#if premium%}Pro{%else%}Free{%/if%}", map[string]interface{}{
		"premium": true,
	})
	if result != "Pro" {
		t.Errorf("expected 'Pro', got '%s'", result)
	}
}

func TestRenderTemplate_IfElseFalse(t *testing.T) {
	result := renderTemplate("{%#if premium%}Pro{%else%}Free{%/if%}", map[string]interface{}{
		"premium": false,
	})
	if result != "Free" {
		t.Errorf("expected 'Free', got '%s'", result)
	}
}

func TestRenderTemplate_NestedIfs(t *testing.T) {
	tmpl := "{%#if a%}{%#if b%}both{%/if%}{%/if%}"
	result := renderTemplate(tmpl, map[string]interface{}{
		"a": true,
		"b": true,
	})
	if result != "both" {
		t.Errorf("expected 'both', got '%s'", result)
	}
}

func TestRenderTemplate_NestedIfsOuterFalse(t *testing.T) {
	tmpl := "{%#if a%}{%#if b%}both{%/if%}{%/if%}"
	result := renderTemplate(tmpl, map[string]interface{}{
		"a": false,
		"b": true,
	})
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestRenderTemplate_DoublebracePreserved(t *testing.T) {
	// {{variable}} syntax should NOT be processed by renderTemplate
	result := renderTemplate("Hello {{name}}", map[string]interface{}{
		"name": "World",
	})
	if result != "Hello {{name}}" {
		t.Errorf("expected '{{name}}' to be preserved, got '%s'", result)
	}
}

func TestRenderTemplate_CombinedVariablesAndIf(t *testing.T) {
	tmpl := "Hi {%name%}, {%#if premium%}thanks for subscribing{%else%}upgrade now{%/if%}!"
	result := renderTemplate(tmpl, map[string]interface{}{
		"name":    "Alice",
		"premium": true,
	})
	if result != "Hi Alice, thanks for subscribing!" {
		t.Errorf("unexpected result: '%s'", result)
	}
}

func TestRenderTemplate_TruthyStringInIf(t *testing.T) {
	result := renderTemplate("{%#if name%}yes{%/if%}", map[string]interface{}{
		"name": "hello",
	})
	if result != "yes" {
		t.Errorf("expected 'yes', got '%s'", result)
	}
}

func TestRenderTemplate_FalsyEmptyStringInIf(t *testing.T) {
	result := renderTemplate("{%#if name%}yes{%else%}no{%/if%}", map[string]interface{}{
		"name": "",
	})
	if result != "no" {
		t.Errorf("expected 'no', got '%s'", result)
	}
}

func TestRenderTemplate_MultipleVariables(t *testing.T) {
	result := renderTemplate("{%first%} {%last%}", map[string]interface{}{
		"first": "John",
		"last":  "Doe",
	})
	if result != "John Doe" {
		t.Errorf("expected 'John Doe', got '%s'", result)
	}
}

// ============================================
// processBlockTag TESTS
// ============================================

func TestProcessBlockTag_BasicConfig(t *testing.T) {
	attrs := map[string]string{
		"config": `{"name":"World"}`,
	}
	result := processBlockTag(attrs, "Hello {%name%}!")
	if result != "Hello World!" {
		t.Errorf("expected 'Hello World!', got '%s'", result)
	}
}

func TestProcessBlockTag_HTMLEntitiesInConfig(t *testing.T) {
	attrs := map[string]string{
		"config": `{&quot;name&quot;:&quot;World&quot;}`,
	}
	result := processBlockTag(attrs, "Hello {%name%}!")
	if result != "Hello World!" {
		t.Errorf("expected 'Hello World!', got '%s'", result)
	}
}

func TestProcessBlockTag_SingleQuotesInConfig(t *testing.T) {
	attrs := map[string]string{
		"config": "{'name':'World'}",
	}
	result := processBlockTag(attrs, "Hello {%name%}!")
	if result != "Hello World!" {
		t.Errorf("expected 'Hello World!', got '%s'", result)
	}
}

func TestProcessBlockTag_AmpersandEntity(t *testing.T) {
	attrs := map[string]string{
		"config": `{"title":"Tom &amp; Jerry"}`,
	}
	result := processBlockTag(attrs, "{%title%}")
	if result != "Tom & Jerry" {
		t.Errorf("expected 'Tom & Jerry', got '%s'", result)
	}
}

func TestProcessBlockTag_EmptyConfig(t *testing.T) {
	attrs := map[string]string{}
	result := processBlockTag(attrs, "Hello {%name ?? Guest%}!")
	if result != "Hello Guest!" {
		t.Errorf("expected 'Hello Guest!', got '%s'", result)
	}
}

func TestProcessBlockTag_InvalidJSON(t *testing.T) {
	attrs := map[string]string{
		"config": "not json",
	}
	result := processBlockTag(attrs, "Hello {%name ?? Fallback%}!")
	if result != "Hello Fallback!" {
		t.Errorf("expected 'Hello Fallback!', got '%s'", result)
	}
}

func TestProcessBlockTag_EmptyTemplate(t *testing.T) {
	attrs := map[string]string{
		"config": `{"name":"World"}`,
	}
	result := processBlockTag(attrs, "")
	// Empty inner + no template attr => empty result
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestProcessBlockTag_TemplateFromAttr(t *testing.T) {
	attrs := map[string]string{
		"config":   `{"name":"World"}`,
		"template": "Hello {%name%}!",
	}
	// Empty inner falls back to template attr
	result := processBlockTag(attrs, "")
	if result != "Hello World!" {
		t.Errorf("expected 'Hello World!', got '%s'", result)
	}
}

// ============================================
// Full pipeline tests using GenerateEmailFromMarkup
// ============================================

func TestFullPipeline_BlockWithParagraph(t *testing.T) {
	markup := `<block config='{"text":"Hello World"}'><paragraph>{%text%}</paragraph></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "<p") {
		t.Error("expected <p> tag in output")
	}
	if !strings.Contains(result, "Hello World") {
		t.Error("expected 'Hello World' in output")
	}
}

func TestFullPipeline_BlockWithHeading(t *testing.T) {
	markup := `<block config='{"title":"Welcome"}'><heading level="2">{%title%}</heading></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "<h2") {
		t.Error("expected <h2> tag in output")
	}
	if !strings.Contains(result, "Welcome") {
		t.Error("expected 'Welcome' in output")
	}
}

func TestFullPipeline_BlockWithButton(t *testing.T) {
	markup := `<block config='{"label":"Click Me","url":"https://example.com"}'><button href="{%url%}">{%label%}</button></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Click Me") {
		t.Error("expected 'Click Me' in output")
	}
	if !strings.Contains(result, "https://example.com") {
		t.Error("expected URL in output")
	}
	if !strings.Contains(result, "<a") {
		t.Error("expected <a> tag in output")
	}
}

func TestFullPipeline_BlockWithImage(t *testing.T) {
	markup := `<block config='{"src":"https://img.example.com/photo.png","alt":"Photo"}'><image src="{%src%}" alt="{%alt%}"></image></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "<img") {
		t.Error("expected <img> tag in output")
	}
	if !strings.Contains(result, "https://img.example.com/photo.png") {
		t.Error("expected image src in output")
	}
	if !strings.Contains(result, "Photo") {
		t.Error("expected alt text in output")
	}
}

func TestFullPipeline_BlockWithSection(t *testing.T) {
	markup := `<block config='{"text":"Inside section"}'><section><paragraph>{%text%}</paragraph></section></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Inside section") {
		t.Error("expected 'Inside section' in output")
	}
	if !strings.Contains(result, "<table") {
		t.Error("expected <table> tag from section in output")
	}
}

func TestFullPipeline_BlockWithLink(t *testing.T) {
	markup := `<block config='{"url":"https://example.com","label":"Go"}'><link href="{%url%}">{%label%}</link></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "<a") {
		t.Error("expected <a> tag in output")
	}
	if !strings.Contains(result, "https://example.com") {
		t.Error("expected URL in output")
	}
	if !strings.Contains(result, "Go") {
		t.Error("expected 'Go' in output")
	}
}

func TestFullPipeline_BlockWithEachLoop(t *testing.T) {
	markup := `<block config='{"items":[{"name":"A"},{"name":"B"},{"name":"C"}]}'>{%#each items as item%}<paragraph>{%item.name%}</paragraph>{%/each%}</block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "A") {
		t.Error("expected 'A' in output")
	}
	if !strings.Contains(result, "B") {
		t.Error("expected 'B' in output")
	}
	if !strings.Contains(result, "C") {
		t.Error("expected 'C' in output")
	}
	// Should have 3 paragraph tags
	if strings.Count(result, "<p") != 3 {
		t.Errorf("expected 3 <p> tags, got %d", strings.Count(result, "<p"))
	}
}

func TestFullPipeline_BlockWithIfElse(t *testing.T) {
	markupTrue := `<block config='{"premium":true}'>{%#if premium%}<paragraph>Pro user</paragraph>{%else%}<paragraph>Free user</paragraph>{%/if%}</block>`
	resultTrue := GenerateEmailFromMarkup(markupTrue, nil)
	if !strings.Contains(resultTrue, "Pro user") {
		t.Error("expected 'Pro user' when premium is true")
	}
	if strings.Contains(resultTrue, "Free user") {
		t.Error("should not contain 'Free user' when premium is true")
	}

	markupFalse := `<block config='{"premium":false}'>{%#if premium%}<paragraph>Pro user</paragraph>{%else%}<paragraph>Free user</paragraph>{%/if%}</block>`
	resultFalse := GenerateEmailFromMarkup(markupFalse, nil)
	if !strings.Contains(resultFalse, "Free user") {
		t.Error("expected 'Free user' when premium is false")
	}
	if strings.Contains(resultFalse, "Pro user") {
		t.Error("should not contain 'Pro user' when premium is false")
	}
}

func TestFullPipeline_DoubleBracePreserved(t *testing.T) {
	markup := `<block config='{"name":"World"}'><paragraph>{{contactName}} says {%name%}</paragraph></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "{{contactName}}") {
		t.Error("expected {{contactName}} to be preserved in output")
	}
	if !strings.Contains(result, "World") {
		t.Error("expected 'World' in output")
	}
}

func TestFullPipeline_MultipleBlocks(t *testing.T) {
	markup := `<block config='{"title":"Header"}'><heading level="1">{%title%}</heading></block><block config='{"body":"Content here"}'><paragraph>{%body%}</paragraph></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Header") {
		t.Error("expected 'Header' in output")
	}
	if !strings.Contains(result, "Content here") {
		t.Error("expected 'Content here' in output")
	}
	if !strings.Contains(result, "<h1") {
		t.Error("expected <h1> tag in output")
	}
	if !strings.Contains(result, "<p") {
		t.Error("expected <p> tag in output")
	}
}

func TestFullPipeline_FallbackValues(t *testing.T) {
	markup := `<block config='{}'><paragraph>{%greeting ?? Hello%} {%name ?? Guest%}</paragraph></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Hello") {
		t.Error("expected fallback 'Hello' in output")
	}
	if !strings.Contains(result, "Guest") {
		t.Error("expected fallback 'Guest' in output")
	}
}

func TestFullPipeline_FallbackOverriddenByValue(t *testing.T) {
	markup := `<block config='{"name":"Alice"}'><paragraph>{%name ?? Guest%}</paragraph></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Alice") {
		t.Error("expected 'Alice' in output")
	}
	if strings.Contains(result, "Guest") {
		t.Error("should not contain fallback 'Guest' when value is present")
	}
}

func TestFullPipeline_BlockNoConfig(t *testing.T) {
	// Block without config attribute still renders template (with empty config)
	markup := `<block><paragraph>Static content</paragraph></block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Static content") {
		t.Error("expected 'Static content' in output")
	}
}

// ============================================
// Lang and Dir tests
// ============================================

func TestLangAttribute_CustomLang(t *testing.T) {
	markup := `<mail lang="tr"><body><paragraph>Merhaba</paragraph></body></mail>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, `lang="tr"`) {
		t.Error("expected lang=\"tr\" in output")
	}
}

func TestDirAttribute_RTL(t *testing.T) {
	markup := `<mail dir="rtl"><body><paragraph>مرحبا</paragraph></body></mail>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, `dir="rtl"`) {
		t.Error("expected dir=\"rtl\" in output")
	}
}

func TestLangAndDir_BothCustom(t *testing.T) {
	markup := `<mail lang="ar" dir="rtl"><body><paragraph>مرحبا</paragraph></body></mail>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, `lang="ar"`) {
		t.Error("expected lang=\"ar\" in output")
	}
	if !strings.Contains(result, `dir="rtl"`) {
		t.Error("expected dir=\"rtl\" in output")
	}
}

func TestLangAndDir_Defaults(t *testing.T) {
	markup := `<mail><body><paragraph>Hello</paragraph></body></mail>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, `lang="en"`) {
		t.Error("expected default lang=\"en\" in output")
	}
	if !strings.Contains(result, `dir="ltr"`) {
		t.Error("expected default dir=\"ltr\" in output")
	}
}

func TestLangAndDir_EmailTag(t *testing.T) {
	markup := `<email lang="fr" dir="ltr"><body><paragraph>Bonjour</paragraph></body></email>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, `lang="fr"`) {
		t.Error("expected lang=\"fr\" in output")
	}
}

func TestLangAndDir_HeadSettingsOverride(t *testing.T) {
	markup := `<mail lang="tr"><body><paragraph>Test</paragraph></body></mail>`
	settings := &EmailHeadSettings{Lang: "de", Dir: "rtl"}
	result := GenerateEmailFromMarkup(markup, settings)
	if !strings.Contains(result, `lang="de"`) {
		t.Error("expected lang=\"de\" from headSettings override")
	}
	if !strings.Contains(result, `dir="rtl"`) {
		t.Error("expected dir=\"rtl\" from headSettings override")
	}
}

func TestFullPipeline_ComplexBlock(t *testing.T) {
	markup := `<block config='{"title":"Newsletter","items":[{"name":"Feature A","desc":"Fast"},{"name":"Feature B","desc":"Secure"}],"showFooter":true,"footer":"Thanks!"}'>
<heading level="1">{%title%}</heading>
{%#each items as item%}<paragraph><b>{%item.name%}</b>: {%item.desc%}</paragraph>{%/each%}
{%#if showFooter%}<paragraph>{%footer%}</paragraph>{%/if%}
</block>`
	result := GenerateEmailFromMarkup(markup, nil)
	if !strings.Contains(result, "Newsletter") {
		t.Error("expected 'Newsletter' in output")
	}
	if !strings.Contains(result, "Feature A") {
		t.Error("expected 'Feature A' in output")
	}
	if !strings.Contains(result, "Feature B") {
		t.Error("expected 'Feature B' in output")
	}
	if !strings.Contains(result, "Fast") {
		t.Error("expected 'Fast' in output")
	}
	if !strings.Contains(result, "Secure") {
		t.Error("expected 'Secure' in output")
	}
	if !strings.Contains(result, "Thanks!") {
		t.Error("expected 'Thanks!' in output")
	}
}
