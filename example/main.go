package main

import (
	"fmt"
	"strconv"

	"github.com/realPy/hogosurudatatable"
	"github.com/realPy/hogosurupagination"

	"github.com/realPy/hogosuru"
	"github.com/realPy/hogosuru/document"
	"github.com/realPy/hogosuru/documentfragment"
	"github.com/realPy/hogosuru/event"
	"github.com/realPy/hogosuru/hogosurudebug"
	"github.com/realPy/hogosuru/htmlanchorelement"
	"github.com/realPy/hogosuru/htmlelement"
	"github.com/realPy/hogosuru/htmltablecellelement"
	"github.com/realPy/hogosuru/htmltableelement"
	"github.com/realPy/hogosuru/htmltablesectionelement"
	"github.com/realPy/hogosuru/htmltemplateelement"
	"github.com/realPy/hogosuru/node"
	"github.com/realPy/hogosuru/promise"
)

type GlobalContainer struct {
	parentNode node.Node
	DataTable  hogosurudatatable.DataTable
	pagination hogosurupagination.Pagination
	page       int
}

var template htmltemplateelement.HtmlTemplateElement

func (g *GlobalContainer) Table(d document.Document) (htmltableelement.HtmlTableElement, error) {
	var t htmltableelement.HtmlTableElement
	var err error
	if t, err = htmltableelement.New(d); hogosuru.AssertErr(err) {
		t.Style_().SetProperty("min-width", "100%")
	}

	return t, err
}

func (g *GlobalContainer) Thead(d document.Document) (htmltablesectionelement.HtmlTableSectionElement, error) {
	var thead htmltablesectionelement.HtmlTableSectionElement
	var err error

	if thead, err = htmltablesectionelement.NewTHead(d); hogosuru.AssertErr(err) {
		thead.SetID("customthead")
	}

	return thead, err
}

func (g *GlobalContainer) Tbody(d document.Document) (htmltablesectionelement.HtmlTableSectionElement, error) {
	var tbody htmltablesectionelement.HtmlTableSectionElement
	var err error

	if tbody, err = htmltablesectionelement.NewTBody(d); hogosuru.AssertErr(err) {

	}

	return tbody, err
}

func (g *GlobalContainer) Columns() int {

	return 3
}
func (g *GlobalContainer) Rows() int {
	return 53
}

func (g *GlobalContainer) MaxRowsByPage() int {
	return 10
}

func (g *GlobalContainer) LoadData() *promise.Promise {

	/*
		p, _ := promise.SetTimeout(func() (interface{}, error) {
			println("Data loaded")
			return nil, nil
		}, 500)
	return &p*/
	return nil
}

func (g *GlobalContainer) Cell(d document.Document, indexRow, indexColumn int) (htmltablecellelement.HtmlTableCellElement, error) {

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

func (g *GlobalContainer) Head(d document.Document, indexColumn int) (htmltablecellelement.HtmlTableCellElement, error) {
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
	htmltemplateelement.GetInterface()
	documentfragment.GetInterface()
	htmltablecellelement.GetInterface()
	htmlanchorelement.GetInterface()

	if elem, err := d.GetElementById("custom-table"); hogosuru.AssertErr(err) {

		if elem, err := elem.Discover(); hogosuru.AssertErr(err) {

			if t, ok := elem.(htmltemplateelement.HtmlTemplateElement); ok {
				template = t
			}
		}

	}

	//if no promise return we dont wait all childs to append
	w.DataTable.Data = w
	w.DataTable.InitialCurrentPage = 0
	w.DataTable.FieldEmptyLine = true
	w.page = 6
	w.pagination.IDPatternElem = "item-pattern"
	w.pagination.IDTemplate = "pagination-tpl"

	w.pagination.OnConfigureItem = func(elem htmlelement.HtmlElement, page int) {

		if link, err := elem.QuerySelector("#link-pattern"); hogosuru.AssertErr(err) {

			if aobj, err := link.Discover(); hogosuru.AssertErr(err) {

				if a, ok := aobj.(htmlanchorelement.HtmlAnchorElement); ok {

					if page >= 0 {

						a.SetTextContent(fmt.Sprintf("%d", page+1))

						a.OnClick(func(e event.Event) {
							w.pagination.Select(elem, page)

							w.DataTable.Jump(page)
							e.PreventDefault()

						})
					} else {
						a.SetTextContent("...")
					}

				}
			}

		}

		w.pagination.OnSelectItem = func(elem htmlelement.HtmlElement) {
			class, _ := elem.ClassName()
			elem.SetClassName(class + " selected")
		}

		/*
			if page >= 0 {

				a.SetTextContent(fmt.Sprintf("%d", page+1))
				a.OnClick(func(e event.Event) {
					w.pagination.Select(page)

					w.DataTable.Jump(page)
					e.PreventDefault()

				})
			} else {
				a.SetTextContent("...")
			}
		*/

	}

	return nil, []hogosuru.Rendering{&w.DataTable, &w.pagination}
}

func (w *GlobalContainer) Node(r hogosuru.Rendering) node.Node {

	if r == &w.pagination {
		if d, err := document.New(); hogosuru.AssertErr(err) {

			if elem, err := d.GetElementById("pagination"); hogosuru.AssertErr(err) {
				return elem.Node
			}

		}
	}

	if r == &w.DataTable {
		if d, err := document.New(); hogosuru.AssertErr(err) {

			if elem, err := d.GetElementById("container"); hogosuru.AssertErr(err) {
				return elem.Node
			}

		}
	}

	return w.parentNode
}

func (w *GlobalContainer) OnEndChildRendering(r hogosuru.Rendering) {

}

func (w *GlobalContainer) OnEndChildsRendering() {
	w.pagination.SetMax(w.page)
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
