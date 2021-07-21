package models

import (
	"fmt"
	"gowiki/consts"
	"io/ioutil"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) Save() error {
	filename := fmt.Sprintf("%s/%s.txt", consts.ViewDir, p.Title)

	return ioutil.WriteFile(filename, p.Body, 0600)
}

