package markup

import (
	"fmt"
	"regexp"
	"strings"
)

func Render(markupContent string, headSettings ...*EmailHeadSettings) (html string, err error) {
	// the processors stay plain string functions and report a tripped limit through
	// failRender; this recover is the only place that turns it into the returned error
	defer func() {
		if recovered := recover(); recovered != nil {
			renderErr, ok := recovered.(*RenderError)
			if !ok {
				panic(recovered)
			}
			html = ""
			err = renderErr
		}
	}()

	if len(headSettings) > 0 {
		return renderDocument(markupContent, headSettings[0]), nil
	}
	return renderDocument(markupContent, nil), nil
}

func renderDocument(htmlContent string, headSettings *EmailHeadSettings) string {
	parsed := ParseEmailHTML(htmlContent)
	settings := parsed.HeadSettings
	if headSettings != nil {
		settings = *headSettings
	}

	processed, gapStyles := processMarkup(normalizeMarkup(parsed.Body))

	lang := settings.Lang
	if lang == "" {
		lang = "en"
	}
	dir := settings.Dir
	if dir == "" {
		dir = "ltr"
	}

	titleTag := ""
	if settings.Title != "" {
		titleTag = fmt.Sprintf("<title>%s</title>", settings.Title)
	}
	customStyles := ""
	if settings.Styles != "" {
		customStyles = fmt.Sprintf(`<style type="text/css">%s</style>`, settings.Styles)
	}
	gapStyleTag := ""
	if gapStyles != "" {
		gapStyleTag = fmt.Sprintf(`<style type="text/css">%s</style>`, gapStyles)
	}
	previewText := ""
	if settings.PreviewText != "" {
		previewText = fmt.Sprintf(`<div style="display:none;font-size:1px;color:#ffffff;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;">%s</div>`, settings.PreviewText)
	}

	return fmt.Sprintf(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="%s" dir="%s" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
<meta content="text/html; charset=UTF-8" http-equiv="Content-Type"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<meta name="x-apple-disable-message-reformatting"/>
<meta content="IE=edge" http-equiv="X-UA-Compatible"/>
<meta name="format-detection" content="telephone=no,address=no,email=no,date=no,url=no"/>
<!--[if mso]>
<noscript>
<xml>
<o:OfficeDocumentSettings>
<o:AllowPNG/>
<o:PixelsPerInch>96</o:PixelsPerInch>
</o:OfficeDocumentSettings>
</xml>
</noscript>
<![endif]-->
<style type="text/css">
#outlook a { padding: 0; }
body { margin: 0; padding: 0; -webkit-text-size-adjust: 100%%; -ms-text-size-adjust: 100%%; }
table, td { border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
.sevk-row-table { border-collapse: separate !important; }
img { border: 0; height: auto; line-height: 100%%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; }
@media only screen and (max-width: 479px) {
  .sevk-row-table { width: 100%% !important; }
  .sevk-column { display: block !important; width: 100%% !important; max-width: 100%% !important; box-sizing: border-box !important; }
}
</style>
%s
%s
%s
%s
</head>
<body style="margin:0;padding:0;word-spacing:normal;-webkit-text-size-adjust:100%%;-ms-text-size-adjust:100%%;font-family:ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
<div aria-roledescription="email" role="article">
%s
%s
</div>
</body>
</html>`, lang, dir, gapStyleTag, titleTag, generateFontLinks(settings.Fonts), customStyles, previewText, processed)
}

func ParseEmailHTML(content string) ParsedEmailContent {
	if regexp.MustCompile(`(?i)<(?:email|mail)[\s>]`).MatchString(content) {
		return parseSevkMarkup(content)
	}
	result := ParsedEmailContent{Body: content}
	result.HeadSettings.Title = extractTagContent(content, "title")
	styles := regexp.MustCompile(`(?i)<style[^>]*>([\s\S]*?)</style>`).FindAllStringSubmatch(content, -1)
	styleContents := make([]string, 0, len(styles))
	for _, style := range styles {
		styleContents = append(styleContents, style[1])
	}
	result.HeadSettings.Styles = strings.Join(styleContents, "\n")
	fonts := regexp.MustCompile(`(?i)<link\b[^>]*href=["']([^"']*fonts[^"']*)["'][^>]*>`).FindAllStringSubmatch(content, -1)
	for index, font := range fonts {
		name := fmt.Sprintf("Font %d", index+1)
		if family := regexp.MustCompile(`family=([^:&]+)`).FindStringSubmatch(font[1]); len(family) > 1 {
			name = strings.ReplaceAll(family[1], "+", " ")
		}
		result.HeadSettings.Fonts = append(result.HeadSettings.Fonts, FontConfig{
			ID: fmt.Sprintf("font-%d", index), Name: name, URL: font[1],
		})
	}

	if body := regexp.MustCompile(`(?i)<body[^>]*>([\s\S]*?)</body>`).FindStringSubmatch(content); len(body) > 1 {
		result.Body = body[1]
	}
	if inner, full, found := findHiddenPreviewDiv(content); found {
		result.HeadSettings.PreviewText = strings.TrimSpace(inner)
		result.Body = strings.Replace(result.Body, full, "", 1)
	}
	result.Body = strings.TrimSpace(result.Body)
	return result
}

func parseSevkMarkup(content string) ParsedEmailContent {
	var settings EmailHeadSettings

	match := regexp.MustCompile(`(?i)<email([^>]*)>`).FindStringSubmatch(content)
	if len(match) < 2 {
		match = regexp.MustCompile(`(?i)<mail([^>]*)>`).FindStringSubmatch(content)
	}
	if len(match) > 1 {
		if lang := regexp.MustCompile(`(?i)lang=["']([^"']+)["']`).FindStringSubmatch(match[1]); len(lang) > 1 {
			settings.Lang = lang[1]
		}
		if dir := regexp.MustCompile(`(?i)dir=["']([^"']+)["']`).FindStringSubmatch(match[1]); len(dir) > 1 {
			settings.Dir = dir[1]
		}
	}

	settings.Title = extractTagContent(content, "title")
	settings.PreviewText = extractTagContent(content, "preview")
	settings.Styles = extractTagContent(content, "style")

	fontPattern := regexp.MustCompile(`(?i)<font([^>]*)>`)
	seenFontURLs := map[string]bool{}
	for _, match := range fontPattern.FindAllStringSubmatch(content, -1) {
		name := regexp.MustCompile(`(?i)name=["']([^"']*)["']`).FindStringSubmatch(match[1])
		url := regexp.MustCompile(`(?i)url=["']([^"']*)["']`).FindStringSubmatch(match[1])
		if len(name) > 1 && len(url) > 1 && !seenFontURLs[url[1]] {
			seenFontURLs[url[1]] = true
			settings.Fonts = append(settings.Fonts, FontConfig{
				ID:   fmt.Sprintf("font-%d", len(settings.Fonts)),
				Name: name[1],
				URL:  url[1],
			})
		}
	}

	bodyPattern := regexp.MustCompile(`(?i)<body[^>]*>([\s\S]*?)</body>`)
	if match := bodyPattern.FindStringSubmatch(content); len(match) > 1 && strings.TrimSpace(match[1]) != "" {
		return ParsedEmailContent{
			Body:         strings.TrimSpace(match[1]),
			HeadSettings: settings,
		}
	}

	body := content
	patterns := []string{
		`<(?:email|mail)[^>]*>`, `</(?:email|mail)>`,
		`<head[^>]*>[\s\S]*?</head>`,
		`<title[^>]*>[\s\S]*?</title>`,
		`<preview[^>]*>[\s\S]*?</preview>`,
		`<style[^>]*>[\s\S]*?</style>`,
		`<font[^>]*>`, `</font>`,
	}
	for _, pattern := range patterns {
		body = regexp.MustCompile("(?i)"+pattern).ReplaceAllString(body, "")
	}

	return ParsedEmailContent{
		Body:         strings.TrimSpace(body),
		HeadSettings: settings,
	}
}

func extractTagContent(content, tag string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?i)<%s[^>]*>([\s\S]*?)</%s>`, tag, tag))
	if match := pattern.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func normalizeMarkup(content string) string {
	result := content
	result = regexp.MustCompile(`(?i)<link(\s+(?:[\w-]+=(?:"[^"]*"|'[^']*')\s+)*)href=`).ReplaceAllString(result, "<sevk-link${1}href=")
	result = regexp.MustCompile(`(?i)</link>`).ReplaceAllString(result, "</sevk-link>")
	if !strings.Contains(result, "<sevk-email") &&
		!strings.Contains(result, "<email") &&
		!strings.Contains(result, "<mail") {
		result = fmt.Sprintf("<mail><body>%s</body></mail>", result)
	}
	return result
}

func generateFontLinks(fonts []FontConfig) string {
	links := make([]string, 0, len(fonts))
	for _, font := range fonts {
		links = append(links, fmt.Sprintf(`<link href="%s" rel="stylesheet" type="text/css" />`, font.URL))
	}
	return strings.Join(links, "\n")
}
