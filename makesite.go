package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"path/filepath"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
}

func createPageFromTextFile(filePath string) Page {
	// Make sure we can read in the file first!
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Get the name of the file without `.txt` at the end.
	// We'll use this later when naming our new HTML file.
	fileNameWithoutExtension := strings.Split(filePath, ".txt")[0]

	// Instantiate a new Page.
	// Populate each field and return the data.
	return Page{
		TextFilePath: filePath,
		TextFileName: fileNameWithoutExtension,
		HTMLPagePath: fileNameWithoutExtension + ".html",
		Content:      string(fileContents),
	}
}

func renderTemplateFromPage(templateFilePath string, page Page) {
	// Create a new template in memory named "template.tmpl".
	// When the template is executed, it will parse template.tmpl,
	// looking for {{ }} where we can inject content.
	t := template.Must(template.New(templateFilePath).ParseFiles(templateFilePath))

	// Create a new, blank HTML file.
	newFile, err := os.Create(page.HTMLPagePath)
	if err != nil {
		panic(err)
	}

	// Executing the template injects the Page instance's data,
	// allowing us to render the content of our text file.
	// Furthermore, upon execution, the rendered template will be
	// saved inside the new file we created earlier.
	t.Execute(newFile, page)
	fmt.Println("Generated File: ", page.HTMLPagePath)
	fileContents, err := ioutil.ReadFile(page.HTMLPagePath)
	fmt.Println(fileContents)
}

// func writeFile(fileName String, text []byte) {
// 	paths := []string{
// 		"template.tmpl",
// 	}
//
// 	t := template.Must(template.New{"template.tmpl"}.ParseFiles(paths...))
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	err =	t.Execute(file, template.HTML{text})
// 	if err != nil{
// 		panic(err)
// 	}
// }

func findTxtFiles(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}


func main() {
	// This flag represents the name of any `.txt` file in the same directory as your program.
	// Run `./makesite --file=latest-post.txt` to test.
	var textFilePath string
	flag.StringVar(&textFilePath, "file", "", "Name or Path to a text file")
	var dir string
	flag.StringVar(&dir, "dir", ".", "Directory")
	flag.Parse()

	// Make sure the `file` flag isn't blank.
	if textFilePath == "" {
		panic("Missing the --file flag! Please supply one.")
	}

	fmt.Println()

	// Read the provided text file and store it's information in a struct.
	newPage := createPageFromTextFile(textFilePath)

	// Use the struct to generate a new HTML page based on the provided template.
	renderTemplateFromPage("template.tmpl", newPage)

	//find all txt files
	files, err := findTxtFiles("/root/", "*.txt")
}
