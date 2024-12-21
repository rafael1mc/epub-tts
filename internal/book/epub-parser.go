package book

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
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

	book := Epub{
		Sections: []EpubSection{},
	}

	// Get content in order of the spine
	for _, spineItem := range packageData.Spine {
		manifestItem := findManifestItem(packageData.Manifest, spineItem.IDRef)
		if manifestItem != nil {
			content, err := readFileFromZip(r, filepath.Join(basePath, manifestItem.Href))
			if err != nil {
				return Epub{}, err
			}

			book.Sections = append(book.Sections, EpubSection{
				ID:          manifestItem.ID,
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
