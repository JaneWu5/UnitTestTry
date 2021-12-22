package gowiki

import (
	"os"
)

// Page A wiki consists of a series of interconnected pages, each of which has a title and a body (the page content).
// Here, we define Page as a struct with two fields representing the title and body.
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := "./resources/data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := "./resources/data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

//func main() {
//	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
//	err := p1.save()
//	if err != nil {
//		fmt.Errorf("save err: %v", err)
//		return
//	}
//	p2, _ := loadPage("TestPage")
//	fmt.Println(string(p2.Body))
//}
