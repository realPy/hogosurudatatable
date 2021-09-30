package main

import (
	"hogosurudatatable"
	"strconv"

	"github.com/realPy/hogosuru"
	"github.com/realPy/hogosuru/document"
	"github.com/realPy/hogosuru/documentfragment"
	"github.com/realPy/hogosuru/event"
	"github.com/realPy/hogosuru/hogosurudebug"
	"github.com/realPy/hogosuru/htmlanchorelement"
	"github.com/realPy/hogosuru/htmltablecellelement"
	"github.com/realPy/hogosuru/htmltableelement"
	"github.com/realPy/hogosuru/htmltablesectionelement"
	htmltemplatelement "github.com/realPy/hogosuru/htmltemplateelement"
	"github.com/realPy/hogosuru/node"
	"github.com/realPy/hogosuru/promise"
)

type GlobalContainer struct {
	parentNode node.Node
	DataTable  hogosurudatatable.DataTable
}

var template htmltemplatelement.HtmlTemplateElement

type CustomDataTable struct {
	data map[string]string
}

func (c CustomDataTable) Table(d document.Document) (htmltableelement.HtmlTableElement, error) {
	var t htmltableelement.HtmlTableElement
	var err error
	if t, err = htmltableelement.New(d); hogosuru.AssertErr(err) {
		t.Style_().SetProperty("min-width", "100%")
	}

	return t, err
}

func (c CustomDataTable) Thead(d document.Document) (htmltablesectionelement.HtmlTableSectionElement, error) {
	var thead htmltablesectionelement.HtmlTableSectionElement
	var err error

	if thead, err = htmltablesectionelement.NewTHead(d); hogosuru.AssertErr(err) {
		thead.SetID("customthead")
	}

	return thead, err
}

func (c CustomDataTable) Tbody(d document.Document) (htmltablesectionelement.HtmlTableSectionElement, error) {
	var tbody htmltablesectionelement.HtmlTableSectionElement
	var err error

	if tbody, err = htmltablesectionelement.NewTBody(d); hogosuru.AssertErr(err) {

	}

	return tbody, err
}

func (c CustomDataTable) Columns() int {

	return 3
}
func (c CustomDataTable) Rows() int {
	return 53
}

func (c CustomDataTable) MaxRowsByPage() int {
	return 10
}

func (c CustomDataTable) LoadData() *promise.Promise {
	/*
		p, _ := promise.SetTimeout(func() (interface{}, error) {
			println("Data loaded")
			return nil, nil
		}, 500)
	return &p*/
	return nil
}

func (c CustomDataTable) Cell(d document.Document, indexRow, indexColumn int) (htmltablecellelement.HtmlTableCellElement, error) {

	var td htmltablecellelement.HtmlTableCellElement
	var err error

	if td, err = htmltablecellelement.NewTd(d); hogosuru.AssertErr(err) {
		if indexRow >= 0 {
			td.SetTextContent(strconv.Itoa(indexColumn) + " " + strconv.Itoa(indexRow))
		} else {
			td.SetTextContent("...")
		}

	}
	return td, err
}

func (c CustomDataTable) Head(d document.Document, indexColumn int) (htmltablecellelement.HtmlTableCellElement, error) {
	var th htmltablecellelement.HtmlTableCellElement
	var err error

	if fragment, err := template.Content(); hogosuru.AssertErr(err) {
		if cth, err := fragment.GetElementById("custom-table-th"); hogosuru.AssertErr(err) {
			if clone, err := d.ImportNode(cth.Node, true); hogosuru.AssertErr(err) {
				if t, ok := clone.(htmltablecellelement.HtmlTableCellElement); ok {
					th = t

					if div, err := t.QuerySelector("#custom-table-content"); hogosuru.AssertErr(err) {

						switch indexColumn {
						case 0:
							div.SetTextContent("Index")
						case 1:
							div.SetTextContent("Value")
						case 2:
							div.SetTextContent("....")
						}

					}

				}
			}
		}
	}

	return th, err
}

func (w *GlobalContainer) OnLoad(d document.Document, n node.Node, route string) (*promise.Promise, []hogosuru.Rendering) {

	w.parentNode = n
	htmltemplatelement.GetInterface()
	documentfragment.GetInterface()
	htmltablecellelement.GetInterface()
	htmlanchorelement.GetInterface()

	if elem, err := d.GetElementById("custom-table"); hogosuru.AssertErr(err) {

		if elem, err := elem.Discover(); hogosuru.AssertErr(err) {

			if t, ok := elem.(htmltemplatelement.HtmlTemplateElement); ok {
				template = t
			}
		}

	}

	//if no promise return we dont wait all childs to append
	w.DataTable.Data = CustomDataTable{}
	w.DataTable.InitialCurrentPage = 0
	w.DataTable.FieldEmptyLine = true

	if a, err := d.GetElementById("jump-right"); hogosuru.AssertErr(err) {
		if a, err := a.Discover(); hogosuru.AssertErr(err) {

			if alink, ok := a.(htmlanchorelement.HtmlAnchorElement); ok {
				alink.OnClick(func(e event.Event) {

					w.DataTable.Jump(w.DataTable.CurrentPage() + 1)

					e.PreventDefault()
				})

			}
		}

	}

	if a, err := d.GetElementById("jump-left"); hogosuru.AssertErr(err) {
		if a, err := a.Discover(); hogosuru.AssertErr(err) {

			if alink, ok := a.(htmlanchorelement.HtmlAnchorElement); ok {
				alink.OnClick(func(e event.Event) {

					w.DataTable.Jump(w.DataTable.CurrentPage() - 1)

					e.PreventDefault()
				})

			}
		}

	}

	return nil, []hogosuru.Rendering{&w.DataTable}
}

func (w *GlobalContainer) Node(r hogosuru.Rendering) node.Node {

	return w.parentNode
}

func (w *GlobalContainer) OnEndChildRendering(r hogosuru.Rendering) {

	if r == &w.DataTable {
		if d, err := document.New(); hogosuru.AssertErr(err) {

			if elem, err := d.GetElementById("container"); hogosuru.AssertErr(err) {
				elem.AppendChild(r.Node(r))
			}

		}
	}

}

func (w *GlobalContainer) OnEndChildsRendering() {

	//w.parentNode.AppendChild(tree)
}

func (w *GlobalContainer) OnUnload() {

}

func main() {

	hogosurudebug.EnableDebug()
	hogosuru.Router().DefaultRendering(&GlobalContainer{})
	hogosuru.Router().Start(hogosuru.HASHROUTE)
	ch := make(chan struct{})
	<-ch

}
