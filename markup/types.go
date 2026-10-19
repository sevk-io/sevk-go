package markup

type FontConfig struct {
	ID   string
	Name string
	URL  string
}

type EmailHeadSettings struct {
	Title       string
	PreviewText string
	Styles      string
	Fonts       []FontConfig
	Lang        string
	Dir         string
}

type ParsedEmailContent struct {
	Body         string
	HeadSettings EmailHeadSettings
}
