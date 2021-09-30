package hogosurudatatable

import (
	"errors"

	"github.com/realPy/hogosuru"
	"github.com/realPy/hogosuru/document"
	"github.com/realPy/hogosuru/hogosurudebug"
	"github.com/realPy/hogosuru/htmltablecellelement"
	"github.com/realPy/hogosuru/htmltableelement"
	"github.com/realPy/hogosuru/htmltablerowelement"
	"github.com/realPy/hogosuru/htmltablesectionelement"
	"github.com/realPy/hogosuru/node"
	"github.com/realPy/hogosuru/promise"
)

//go:generate go run cmd/csscompact.go hogosurudatatable

//DataTableBind interface
type DataTableBind interface {
	Columns() int
	Table(d document.Document) (htmltableelement.HtmlTableElement, error)
	Thead(d document.Document) (htmltablesectionelement.HtmlTableSectionElement, error)
	Head(d document.Document, indexColumn int) (htmltablecellelement.HtmlTableCellElement, error)

	Rows() int

	Tbody(d document.Document) (htmltablesectionelement.HtmlTableSectionElement, error)
	Cell(d document.Document, indexRow, indexColumn int) (htmltablecellelement.HtmlTableCellElement, error)

	MaxRowsByPage() int

	LoadData() *promise.Promise
}

type DataTable struct {
	Data               DataTableBind
	InitialCurrentPage int
	FieldEmptyLine     bool

	parentNode  node.Node
	container   htmltableelement.HtmlTableElement
	loadPromise *promise.Promise
	tbody       htmltablesectionelement.HtmlTableSectionElement
	nbcolumns   int
	currentpage int
}

func (t *DataTable) OnLoad(d document.Document, n node.Node, route string) (*promise.Promise, []hogosuru.Rendering) {
	var err error
	t.parentNode = n
	var p *promise.Promise = nil

	if t.container.Empty() {

		if t.container, err = t.Data.Table(d); err == nil {
			if t.Data != nil {

				t.nbcolumns = t.Data.Columns()

				if thead, err := t.Data.Thead(d); err == nil {
					if !thead.Empty() {
						if tr, err := htmltablerowelement.New(d); hogosuru.AssertErr(err) {
							for i := 0; i < t.nbcolumns; i++ {

								if th, err := t.Data.Head(d, i); err == nil {
									tr.AppendChild(th.Node)
								}
							}

							thead.AppendChild(tr.Node)
						}
						t.container.AppendChild(thead.Node)
					}
				}
				var err error
				if t.tbody, err = t.Data.Tbody(d); err == nil {

					if !t.tbody.Empty() {
						t.container.AppendChild(t.tbody.Node)
					}

				}

			}

		}
	}

	return p, nil
}

func (t *DataTable) OnEndChildRendering(r hogosuru.Rendering) {

}

func (t *DataTable) CurrentPage() int {

	return t.currentpage
}

func (t *DataTable) Jump(page int) {

	t.refreshData(page)
}

func (t *DataTable) refreshData(page int) {

	if t.loadPromise == nil {
		if d, err := document.New(); hogosuru.AssertErr(err) {
			p, _ := promise.New(func() (interface{}, error) {

				waitingData := t.Data.LoadData()
				var err error

				if waitingData != nil {
					_, err = waitingData.Await()
				}

				if waitingData == nil && err == nil {
					var nbrow int = t.Data.Rows()
					var nbRowByPage = t.Data.MaxRowsByPage()
					var begin, end, currentpage int

					if page < 0 {
						page = 0
					}
					if page > (nbrow / nbRowByPage) {
						currentpage = nbrow / nbRowByPage
					} else {
						currentpage = page
					}

					t.currentpage = currentpage

					begin = nbRowByPage * currentpage

					if (begin + nbRowByPage) < nbrow {

						end = begin + nbRowByPage
					} else {
						end = nbrow
					}

					for r, err := t.tbody.FirstChild(); err == nil; r, err = t.tbody.FirstChild() {
						t.tbody.RemoveChild(r)
					}

					for i := begin; i < end; i++ {

						if tr, err := htmltablerowelement.New(d); hogosuru.AssertErr(err) {
							for j := 0; j < t.nbcolumns; j++ {

								if td, err := t.Data.Cell(d, i, j); err == nil {
									if !td.Empty() {
										tr.AppendChild(td.Node)
									}

								}
							}

							t.tbody.AppendChild(tr.Node)
						}

					}

					if t.FieldEmptyLine && (begin-end) < nbRowByPage {
						for i := 0; i < nbRowByPage-(end-begin); i++ {

							if tr, err := htmltablerowelement.New(d); hogosuru.AssertErr(err) {
								for j := 0; j < t.nbcolumns; j++ {

									if td, err := t.Data.Cell(d, -1, j); err == nil {
										if !td.Empty() {
											tr.AppendChild(td.Node)
										}

									}
								}
								tr.Style_().SetProperty("visibility", "hidden")
								t.tbody.AppendChild(tr.Node)
							}

						}

					}
				}
				return nil, nil
			})
			t.loadPromise = &p
			t.loadPromise.Finally(func() {
				t.loadPromise = nil
			})

		}
	} else {

		hogosurudebug.AssertDebug(errors.New("A loading is in progress please cancel it"))
	}

}

func (t *DataTable) OnEndChildsRendering() {

	t.parentNode.AppendChild(t.container.Node)
	t.refreshData(t.InitialCurrentPage)

}

func (t *DataTable) Node(r hogosuru.Rendering) node.Node {

	return t.container.Node
}

func (t *DataTable) OnUnload() {

	p, _ := t.container.ParentNode()

	p.RemoveChild(t.container.Node)

}
