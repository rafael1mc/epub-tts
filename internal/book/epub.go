package book

type EpubSection struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	HtmlContent string `json:"htmlString"`
}

type Epub struct {
	Name     string
	Toc      map[string]string
	Sections []EpubSection
}
