// This script is designed to precompile all our documentation template files,
// so they are ready to then be passed onto the official terraform
// documentation generation tool:
// https://github.com/hashicorp/terraform-plugin-docs
//
// The reason why we need to precompile the templates is because they use the
// {{ define "..." }} syntax to allow us to reuse Markdown content across
// multiple files, and so we need this script to expose those separate files as
// if they were defined within a single template file.
//
// The reason we need the official terraform documentation generator tool is
// because it helps to pick up missing content that we've not defined by
// reflecting over the schemas of our resources and data sources.
//
// THE COMPLETE DOCUMENTATION STEPS ARE:
//
// 1. acquire all the templates
// 2. render the templates into Markdown
// 3. write the file output to a temp directory and still use .tmpl extension
// 4. append to each rendered .tmpl file the template code needed by tfplugindocs (e.g. {{ .SchemaMarkdown | trimspace }})
// 5. copy the index.md file (which requires no pre-compiling) to the temporary directory so tfplugindocs can include it
// 6. rename repo /templates/{data-sources/resources} directories to avoid being overwritten by next step
// 7. move contents of temp directory (i.e. data-sources/resources) into repo /templates directory
// 8. run tfplugindocs generate function to output final documentation to /docs directory
// 9. replace /templates/{data-sources/resources} directories with their backed up equivalents
//
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jameswoolfenden/terraform-provider-bridgecrew/version"
)

// PageData represents a type of service and enables a template to render
// different content depending on the service type.
type PageData struct {
	ServiceType     string
	ProviderVersion string
}

// Page represents a template page to be rendered.
//
// name: the {{ define "..." }} in each template.
// path: the path to write rendered output to.
// Data: context specific information (e.g. vcl vs wasm).
//
// Data is public as it's called via the template processing logic.
type Page struct {
	name string
	path string
	Data PageData
}

func main() {
	tfplugindocsPath := flag.String("tfplugindocsPath", "bin/tfplugindocs", "location where tfplugindocs is installed")
	flag.Parse()

	if !tfPluginDocsExists(*tfplugindocsPath) {
		log.Fatalf("tfplugindocs not found at '%s' - have you run the Makefile?", *tfplugindocsPath)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tmplDir := baseDir + "/templates"

	tempDir, err := ioutil.TempDir("", "precompile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	var dataPages = []Page{
		{
			name: "data_source_justifications",
			path: tempDir + "/data-sources/justifications.md.tmpl",
		},
		{
			name: "data_source_tags",
			path: tempDir + "/data-sources/tags.md.tmpl",
		},
		{
			name: "data_source_tag",
			path: tempDir + "/data-sources/tag.md.tmpl",
		},
		{
			name: "data_source_mappings",
			path: tempDir + "/data-sources/mappings.md.tmpl",
		},
		{
			name: "data_source_integrations",
			path: tempDir + "/data-sources/integrations.md.tmpl",
		},
		{
			name: "data_source_users",
			path: tempDir + "/data-sources/users.md.tmpl",
		},
		{
			name: "data_source_apitokens",
			path: tempDir + "/data-sources/apitokens.md.tmpl",
		},
		{
			name: "data_source_apitokens_customer",
			path: tempDir + "/data-sources/apitokens_customer.md.tmpl",
		},
		{
			name: "data_source_policies",
			path: tempDir + "/data-sources/policies.md.tmpl",
		},
		{
			name: "data_source_repositories",
			path: tempDir + "/data-sources/repositories.md.tmpl",
		},
		{
			name: "data_source_repository_branches",
			path: tempDir + "/data-sources/repository_branches.md.tmpl",
		},
		{
			name: "data_source_suppressions",
			path: tempDir + "/data-sources/suppressions.md.tmpl",
		},
		{
			name: "data_source_authors",
			path: tempDir + "/data-sources/authors.md.tmpl",
		},
		{
			name: "data_source_incidents",
			path: tempDir + "/data-sources/incidents.md.tmpl",
		},
		{
			name: "data_source_incidents_preset",
			path: tempDir + "/data-sources/incidents_preset.md.tmpl",
		},
		{
			name: "data_source_incidents_info",
			path: tempDir + "/data-sources/incidents_info.md.tmpl",
		},
		{
			name: "data_source_organisation",
			path: tempDir + "/data-sources/organisation.md.tmpl",
		},
	}

	var resourcePages = []Page{
		{
			name: "resource_policy",
			path: tempDir + "/resources/policy.md.tmpl",
		},
		{
			name: "resource_simple_policy",
			path: tempDir + "/resources/simple_policy.md.tmpl",
		},
		{
			name: "resource_complex_policy",
			path: tempDir + "/resources/complex_policy.md.tmpl",
		},
		{
			name: "resource_tag",
			path: tempDir + "/resources/tag.md.tmpl",
		},
	}

	var indexPages = []Page{
		{
			name: "index",
			path: tempDir + "/index.md.tmpl",
			Data: PageData{
				ProviderVersion: strings.Replace(version.ProviderVersion, "v", "", 1),
			},
		},
	}

	pages := append(append(indexPages, resourcePages...), dataPages...)

	renderPages(getTemplate(tmplDir), pages)

	appendSyntaxToFiles(tempDir)

	backupTemplatesDir(tmplDir)

	replaceTemplatesDir(tmplDir, tempDir)

	runTFPluginDocs(*tfplugindocsPath)

	replaceTemplatesDir(tmplDir, tmplDir+"-backup")
}

func tfPluginDocsExists(tfplugindocsLocation string) bool {
	stat, err := os.Stat(tfplugindocsLocation)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

// getTemplate walks the templates' directory, filtering non-tmpl extension
// files, and parsing all the templates found (ensuring they must parse).
func getTemplate(tmplDir string) *template.Template {
	var templateFiles []string
	err := filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tmpl" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Print(err)
	}

	myTemplate, err := template.ParseFiles(templateFiles...)
	if err != nil {
		log.Fatalf("Error parsing template files: %s", err)
	}
	return myTemplate
}

// renderPages iterates over the given pages and renders each element.
func renderPages(t *template.Template, pages []Page) {
	for _, p := range pages {
		renderPage(t, p)
	}
}

// renderPage creates a new file based on the page information given, and
// renders the associated template for that page.
func renderPage(t *template.Template, p Page) {
	basePath := filepath.Dir(p.path)
	err := makeDirectoryIfNotExists(basePath)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(p.path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = t.ExecuteTemplate(f, p.name, p)
	if err != nil {
		panic(err)
	}
}

// makeDirectoryIfNotExists asserts whether a directory exists and makes it
// if not. Returns nil if exists or successfully made.
func makeDirectoryIfNotExists(path string) error {
	fi, err := os.Stat(path)
	switch {
	case err == nil && fi.IsDir():
		return nil
	case err == nil && !fi.IsDir():
		return fmt.Errorf("%s already exists as a regular file", path)
	case os.IsNotExist(err):
		return os.MkdirAll(path, 0750)
	case err != nil:
		return err
	}

	return nil
}

// appendSyntaxToFiles walks the temporary directory finding all the rendered
// Markdown files we generated and proceeds to append the required template
// syntax that the tfplugindocs tool needs.
func appendSyntaxToFiles(tempDir string) {
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tmpl" {
			// open file for appending and in write-only mode
			f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := f.Write([]byte("{{ .SchemaMarkdown | trimspace }}\n")); err != nil {
				f.Close() // ignore error; Write error takes precedence
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Print(err)
	}
}

// backupTemplatesDir renames the /templates directory.
//
// We do this so that we can create a new /templates directory and move the
// contents of the temporary directory into the new location. Thus allowing
// tfplugindocs to be run within the root of the Terraform provider repo.
func backupTemplatesDir(tmplDir string) {

	// os.RemoveAll(tmplDir+"-backup")
	// err := CopyDir(tmplDir, tmplDir+"-backup")
	err := os.Rename(tmplDir, tmplDir+"-backup")
	if err != nil {
		log.Fatal(err)
	}
}

// replaceTemplatesDir removes the template directory from the repo and moves
// the temporary directory to where the template one would have been.
func replaceTemplatesDir(tmplDir string, tempDir string) {
	err := os.RemoveAll(tmplDir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(tempDir, tmplDir)
	if err != nil {
		log.Fatal(err)
	}
}

// runTFPluginDocs executes the tfplugindocs binary which generates
// documentation Markdown files from our terraform code, while also utilizing
// any templates we have defined in the /templates directory.
//
// NOTE: it is presumed that the /templates directory that is referenced will
// consist of precompiled templates and that the original untouched templates
// will still exist in the /templates-backup directory ready to be restored
// once the /docs content has been generated.
func runTFPluginDocs(tfplugindocsLocation string) {
	cmd := exec.Command(tfplugindocsLocation, "generate")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
//goland:noinspection GoUnusedExportedFunction
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
