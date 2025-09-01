package doc

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	TerraformRegistryURL = "https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest"
	DocumentationURL     = "https://registry.terraform.io/providers/ionos-cloud/ionoscloud/latest/docs"
	PrivacyPolicyURL     = "https://www.ionos.com/terms-gtc/terms-privacy/"
	ImprintURL           = "https://www.ionos.de/impressum"
	ResourcesDocsDir     = "resources"
	DataSourcesDocsDir   = "data-sources"
)

type docInfo struct {
	name string
	path string
}

type productInfo struct {
	resources   []docInfo
	dataSources []docInfo
}

func GenerateSummary(docsDir string) error {
	f, err := os.Create(filepath.Join(docsDir, "summary.md"))
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	// Write the first part of the summary.md file
	_, err = buf.WriteString(fmt.Sprintf("# Table of contents\n\n* [Introduction](./../README.md)\n* [Changelog](./../CHANGELOG.md)\n\n## Terraform Registry\n\n* [Terraform Registry](%s)\n* [Documentation](%s)\n\n## Legal\n\n* [Privacy policy](%s)\n* [Imprint](%s)\n\n\n## API\n", TerraformRegistryURL, DocumentationURL, PrivacyPolicyURL, ImprintURL))
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}

	// Write the second part of the summary.md file, the 'products' part.
	// Iterate over the docs, extract the name of the resource/data-source (server, lan, etc.),
	// the subcategory (Compute Engine, API Gateway, etc.) and the path to the doc file.
	// Sort alphabetically and write the content in summary.md file.
	productNames := []string{"Compute Engine", "API Gateway", "DBaaS - PgSQL"}
	pathToResourcesDocsDir := filepath.Join(docsDir, ResourcesDocsDir)
	pathToDataSourcesDocsDir := filepath.Join(docsDir, DataSourcesDocsDir)

	// TODO -- Maybe extract the info for the resources and data sources in the same time using
	// threads.
	// Extract info about resources.
	docs, err := os.ReadDir(pathToResourcesDocsDir)
	productsMap := make(map[productName]productDocs)
	for _, doc := range docs {
		var product productName
		docName := doc.Name()
		docPath := filepath.Join(pathToResourcesDocsDir, docName)
		product = productName(productNames[rand.Intn(len(productNames))])

		p, ok := productsMap[]
		if !ok {
			p = productDocs{}
			p.resources = make([]docInfo, 0)
			p.dataSources = make([]docInfo, 0)
			p.resources = append(p.resources, docInfo{name: docName, path: docPath})
			productsMap[productName(product)] = p
		} else {
			p.resources = append(p.resources, docFile{name: name, path: path})
			productMap[productName] = p
		}
	}

	// Extract info about data-sources.
	docs, err = os.ReadDir(pathToDataSourcesDocsDir)
	for _, doc := range docs {
		name := doc.Name()
		path := filepath.Join(pathToResourcesDocsDir, name)
		productName := productNames[rand.Intn(len(productNames))]

		p := productMap[productName]
		p.dataSources = append(p.dataSources, docFile{name: name, path: path})
		productMap[productName] = p
	}
	return nil
}
