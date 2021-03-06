// +build js

package DOM

import (
	"honnef.co/go/js/dom"
	"strings"
)

func (el *Element) GetAttr(name string) (result string, err error) {
	return el.el.GetAttribute(name), nil
}

func (dom *DOM) GetAttrOf(name string, css string) (result string, err error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return
	}

	return input.GetAttr(name)
}

func (el *Element) GetText() (result string, err error) {
	return el.el.TextContent(), nil
}

func (el *Element) GetStringValue() (result string, err error) {
	if strings.ToUpper(el.el.TagName()) == "TEXTAREA" {
		return el.el.(*dom.HTMLTextAreaElement).Value, nil
	}

	return el.el.(*dom.HTMLInputElement).Value, nil
}

func (dom *DOM) GetStringValueOf(css string) (result string, err error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return
	}

	return input.GetStringValue()
}

func (el *Element) GetBytesValue() (result []byte, err error) {
	r, err := el.GetStringValue()
	return []byte(r), err
}

func (dom *DOM) GetBytesValueOf(css string) (result []byte, err error) {
	input, err := dom.SelectFirstElement(css)
	if err != nil {
		return
	}

	return input.GetBytesValue()
}
