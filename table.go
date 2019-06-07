package main

import (
	"fmt"

	"github.com/ajstarks/svgo"
)

type Table struct {
	X, Y     int
	Title    string
	FontSize int
	Cells    [][]Cell
}

func (t *Table) TitleCell() Cell {
	c := Cell{}
	c.SetText(t.Title)
	return c
}

func (t *Table) TitleRuneWidth() int {
	c := t.TitleCell()
	return t.calcWidth(c.RuneWidthN(), 1)
}

func (t *Table) TitleRuneHeight() int {
	c := t.TitleCell()
	return t.calcHeight(c.LineN(), 1)
}

func (t *Table) SetCells(cells [][]string) (err error) {
	h := len(cells)
	if h == 0 {
		return
	}
	w0 := len(cells[0])
	t.Cells = make([][]Cell, h)
	for j := 0; j < h; j++ {
		wi := len(cells[j])
		if wi != w0 {
			err = fmt.Errorf("[SetCells]: len(cells[%d]) got %d, but required %d", j, wi, w0)
			return
		}
		t.Cells[j] = make([]Cell, w0)
		for i := 0; i < w0; i++ {
			c := Cell{}
			c.SetText(cells[j][i])
			t.Cells[j][i] = c
		}
	}
	return
}

func (t *Table) MarginHeight() int {
	return t.FontSize / 2
}

func (t *Table) MarginWidth() int {
	return t.FontSize / 2
}

func (t *Table) RuneWidth() int {
	return t.FontSize * 2 / 3
}

func (t *Table) RuneHeight() int {
	return t.FontSize
}

func (t *Table) CellXLen() (w int) {
	if t.CellYLen() == 0 {
		return 0
	}
	return len(t.Cells[0])
}

func (t *Table) CellYLen() (w int) {
	return len(t.Cells)
}

func (t *Table) ViewMargin() (w int) {
	canvasMargin := 8
	return canvasMargin
}

func (t *Table) ViewWidth() (w int) {
	titleWidth := t.TitleRuneWidth()
	tableWidth := t.Width()
	w = t.X
	if titleWidth > tableWidth {
		w += titleWidth
	} else {
		w += tableWidth
	}
	w += t.ViewMargin() * 2
	return
}
func (t *Table) ViewHeight() (w int) {
	return t.Y + t.TitleRuneHeight() + t.Height() + t.ViewMargin()*2
}

func (t *Table) TableContentWidthPadding() (w int) {
	if t.Width() > t.TitleRuneWidth() {
		return 0
	}
	return (t.TitleRuneWidth() - t.Width()) / 2
}

func (t *Table) Width() (w int) {
	return t.WidthRange(0, t.CellXLen())
}
func (t *Table) calcWidth(charLen, n int) (w int) {
	return charLen*t.RuneWidth() + t.MarginWidth()*2*n
}
func (t *Table) WidthRange(start, end int) (w int) {
	sum := 0
	for i := start; i < end; i++ {
		max := 0
		for j := 0; j < t.CellYLen(); j++ {
			tmp := t.Cells[j][i].RuneWidthN()
			if tmp > max {
				max = tmp
			}
		}
		sum += max
	}
	w = t.calcWidth(sum, end-start)
	return
}

func (t *Table) Height() (h int) {
	return t.HeightRange(0, t.CellYLen())
}
func (t *Table) calcHeight(charLen, n int) (w int) {
	return charLen*t.RuneHeight() + t.MarginHeight()*2*n
}
func (t *Table) HeightRange(start, end int) (h int) {
	sum := 0
	for j := start; j < end; j++ {
		max := 0
		for i := 0; i < t.CellXLen(); i++ {
			tmp := t.Cells[j][i].LineN()
			if tmp > max {
				max = tmp
			}
		}
		sum += max
	}
	h = t.calcHeight(sum, end-start)
	return
}

func (t *Table) Pos(i, j int) (int, int) {
	return t.X + t.WidthRange(0, i), t.Y + t.TitleRuneHeight() + t.HeightRange(0, j)
}

func (t *Table) calTextYOffset(yIndex int, n int) int {
	return (t.RuneHeight()*(yIndex+1) + t.MarginHeight()*(yIndex+1)/n)
}
func (t *Table) DrawText(canvas *svg.SVG, x, y int, texts []string, style string) {
	n := len(texts)
	xPadding := t.TableContentWidthPadding() + t.ViewMargin()
	for i, text := range texts {
		canvas.Text(x+xPadding+t.MarginWidth(), y+t.calTextYOffset(i, n), text, style)
	}
}
func (t *Table) DrawTextAtMiddle(canvas *svg.SVG, x, y int, texts []string, style string) {
	n := len(texts)
	xPadding := t.ViewMargin()
	for i, text := range texts {
		canvas.Text(x+xPadding, y+t.calTextYOffset(i, n), text, style)
	}
}

func (t *Table) DrawTitleText(canvas *svg.SVG) {
	textStyle := fmt.Sprintf("font-family:Monaco;font-size:%dpx;stroke:black;fill:black;text-anchor:middle;", t.FontSize)
	c := t.TitleCell()
	t.DrawTextAtMiddle(canvas, t.ViewWidth()/2, 0, c.Texts(), textStyle)
}

func (t *Table) DrawTableText(canvas *svg.SVG) {
	textStyle := fmt.Sprintf("font-family:Monaco;font-size:%dpx;stroke:black;fill:black", t.FontSize)
	tx, ty := t.Pos(0, 0)
	for i := 0; i < t.CellXLen(); i++ {
		xOffset := t.WidthRange(0, i)
		for j := 0; j < t.CellYLen(); j++ {
			// NOTE: text baseline is bottom
			yOffset := t.HeightRange(0, j)
			t.DrawText(canvas, tx+xOffset, ty+yOffset, t.Cells[j][i].Texts(), textStyle)
		}
	}
}

func (t *Table) DrawVerticalTableLine(canvas *svg.SVG) {
	tx, ty := t.Pos(0, 0)
	height := t.Height()
	// NOTE: draw vertical line(||||)
	xPadding := t.TableContentWidthPadding() + t.ViewMargin()
	for i := 0; i < t.CellXLen()+1; i++ {
		xOffset := t.WidthRange(0, i)
		canvas.Line(tx+xPadding+xOffset, ty, tx+xPadding+xOffset, ty+height)
	}
}

func (t *Table) DrawHorizontalTableLine(canvas *svg.SVG) {
	tx, ty := t.Pos(0, 0)
	width := t.Width()
	// NOTE: draw horizontal line(----)
	headerLineStyle := "stroke:black;stroke-width:4;"
	xPadding := t.TableContentWidthPadding() + t.ViewMargin()
	for j := 0; j < t.CellYLen()+1; j++ {
		yOffset := t.HeightRange(0, j)
		if j == 1 {
			canvas.Line(tx+xPadding, ty+yOffset, tx+xPadding+width, ty+yOffset, headerLineStyle)
		} else {
			canvas.Line(tx+xPadding, ty+yOffset, tx+xPadding+width, ty+yOffset)
		}
	}
}
