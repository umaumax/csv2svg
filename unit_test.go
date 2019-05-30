package main

import (
	"testing"
)

func TestCellRuneWidthNSuccess(t *testing.T) {
	c := Cell{}
	ioData := []struct {
		in  string
		out int
	}{
		{"", 0},
		{"a", 1},
		{"a\nab", 2},
		{"a\nab\nabc", 3},
	}
	for _, data := range ioData {
		c.SetText(data.in)
		got := c.RuneWidthN()
		if got != data.out {
			t.Fatalf("Input='%s', Output=(got:'%d', req:'%d')", data.in, got, data.out)
		}
	}
}

func TestTableWidthSuccess(t *testing.T) {
	fontSize := 16
	table := Table{
		X: 0, Y: 0, FontSize: fontSize,
	}
	ioData := []struct {
		in  [][]string
		out int
	}{
		{[][]string{{"a", "ab", "abc"}}, table.calcWidth(1+2+3, 3)},
		{[][]string{{"a\nabcd", "ab", "abc"}}, table.calcWidth(4+2+3, 3)},
		{[][]string{{"a", "ab", "abc"}, {"ab", "ab", "abc"}}, table.calcWidth(2+2+3, 3)},
		{[][]string{{"犬", "dog"}}, table.calcWidth(2+3, 2)},
	}
	for _, data := range ioData {
		table.SetCells(data.in)
		got := table.Width()
		if got != data.out {
			t.Fatalf("Input='%s', Output=(got:'%d', req:'%d')", data.in, got, data.out)
		}
	}
}
