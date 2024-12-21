package book

type EpubSection struct {
	ID          string `json:"id"`
	HtmlContent string `json:"htmlString"`
}

type Epub struct {
	Sections []EpubSection
}
