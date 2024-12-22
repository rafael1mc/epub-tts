package book

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"
)

// Container structure to parse container.xml
type Container struct {
	Rootfiles []Rootfile `xml:"rootfiles>rootfile"`
}

// Rootfile structure for OPF reference
type Rootfile struct {
	FullPath string `xml:"full-path,attr"`
}

// Package structure to parse OPF file
type Package struct {
	Manifest []Item    `xml:"manifest>item"`
	Spine    []Itemref `xml:"spine>itemref"`
}

// Item structure for manifest items
type Item struct {
	ID   string `xml:"id,attr"`
	Href string `xml:"href,attr"`
}

// Itemref structure for spine items
type Itemref struct {
	IDRef string `xml:"idref,attr"`
}

// NavPoint represents a navigation point in the EPUB toc.ncx file
type NavPoint struct {
	Text string `xml:"navLabel>text"`
	// Src  string `xml:"content>src,attr"` // this doesn;t work
	Src          string     `xml:"content,attr"`
	SubNavPoints []NavPoint `xml:"navPoint"`
}

func (n *NavPoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var aux struct {
		Text    string `xml:"navLabel>text"`
		Content struct {
			Src string `xml:"src,attr"`
		} `xml:"content"`
		SubNavPoints []NavPoint `xml:"navPoint"`
	}
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}

	n.Text = aux.Text
	n.Src = aux.Content.Src
	n.SubNavPoints = aux.SubNavPoints

	return nil
}

// NCX represents the structure of the toc.ncx file
type NCX struct {
	NavMap []NavPoint `xml:"navMap>navPoint"`
}

func ParseEpub(epubPath string) (Epub, error) {
	// Open the ePUB file as a zip archive
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return Epub{}, err
	}
	defer r.Close()

	// Read the container.xml to locate the OPF file
	container, err := readContainer(r)
	if err != nil {
		return Epub{}, err
	}

	// Parse the OPF file
	packageData, err := readOPF(r, container.Rootfiles[0].FullPath)
	if err != nil {
		return Epub{}, err
	}

	basePath := extractBasePath(container.Rootfiles[0].FullPath)
	tocFileName := findTocFileName(packageData.Manifest)

	// Parse table of contents
	tableOfContents, err := extractTableOfContents(r, basePath, tocFileName)
	if err != nil {
		fmt.Println("Failed to parse table of contents", err)
	}

	book := Epub{
		Toc:      map[string]string{},
		Sections: []EpubSection{},
	}

	// Get content in order of the spine
	for _, spineItem := range packageData.Spine {
		manifestItem := findManifestItem(packageData.Manifest, spineItem.IDRef)
		if manifestItem != nil {
			currFile := filepath.Join(basePath, manifestItem.Href)
			content, err := readFileFromZip(r, currFile)
			if err != nil {
				return Epub{}, err
			}

			title := tableOfContents[manifestItem.Href]
			if title == "" {
				title = tableOfContents[currFile]
			}

			book.Sections = append(book.Sections, EpubSection{
				ID:          manifestItem.ID,
				Title:       title,
				HtmlContent: string(content),
			})
		}
	}

	return book, nil
}

// readContainer reads and parses the container.xml
func readContainer(r *zip.ReadCloser) (*Container, error) {
	content, err := readFileFromZip(r, "META-INF/container.xml")
	if err != nil {
		return nil, err
	}

	var container Container
	if err := xml.Unmarshal(content, &container); err != nil {
		return nil, err
	}

	return &container, nil
}

// readOPF reads and parses the OPF file
func readOPF(r *zip.ReadCloser, opfPath string) (*Package, error) {
	content, err := readFileFromZip(r, opfPath)
	if err != nil {
		return nil, err
	}

	var packageData Package
	if err := xml.Unmarshal(content, &packageData); err != nil {
		return nil, err
	}

	return &packageData, nil
}

// readFileFromZip extracts a file's content from the zip archive
func readFileFromZip(r *zip.ReadCloser, name string) ([]byte, error) {
	for _, file := range r.File {
		if file.Name == name {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("file not found: %s", name)
}

func extractTableOfContents(
	r *zip.ReadCloser,
	basePath string,
	tocFileName string,
) (map[string]string, error) {
	navPoints, err := parseNCX(r, basePath, tocFileName)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	addNavPoints(result, navPoints)

	return result, nil
}

func addNavPoints(m map[string]string, navPoints []NavPoint) {
	for _, v := range navPoints {
		parsedSrc, err := url.Parse(v.Src)
		if err != nil {
			// TODO: log
			m[v.Src] = v.Text
			continue
		}

		parsedSrc.Fragment = ""
		parsedSrc.RawQuery = ""
		src := parsedSrc.String()

		m[src] = v.Text

		if len(v.SubNavPoints) > 0 {
			addNavPoints(m, v.SubNavPoints)
		}
	}
}

func parseNCX(r *zip.ReadCloser, basePath string, tocFileName string) ([]NavPoint, error) {
	ncxContent, err := readFileFromZip(r, filepath.Join(basePath, tocFileName))
	if err != nil {
		return nil, err
	}

	// Parse the NCX XML
	var ncx NCX
	err = xml.Unmarshal(ncxContent, &ncx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return ncx.NavMap, nil
}

// findManifestItem finds a manifest item by ID
func findManifestItem(manifest []Item, id string) *Item {
	for _, item := range manifest {
		if item.ID == id {
			return &item
		}
	}

	return nil
}

func extractBasePath(fullPath string) string {
	parsedFullPath := strings.Split(fullPath, "/")
	var fullPathBase string
	if len(parsedFullPath) > 1 {
		fullPathBase = strings.Join(parsedFullPath[:len(parsedFullPath)-1], "/")
	}

	return fullPathBase
}

func findTocFileName(manifestItems []Item) string {
	for _, v := range manifestItems {
		if strings.Contains(v.ID, "ncx") &&
			strings.Contains(v.Href, "ncx") {
			return v.Href
		}
	}

	return "toc.ncx" // TODO look for any ncx file inside the whole zip
}
