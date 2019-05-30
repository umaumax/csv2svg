package main

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

type Cell struct {
	texts []string
}

func (c *Cell) SetText(text string) (lineN int) {
	// NOTE: trim prefix and suffix space
	c.texts = strings.Split(strings.TrimSpace(strings.Replace(text, "\r\n", "\n", -1)), "\n")
	lineN = len(c.texts)
	return
}
func (c *Cell) Texts() []string {
	return c.texts
}
func (c *Cell) LineN() (lineN int) {
	lineN = len(c.texts)
	return
}
func (c *Cell) RuneWidthN() (runeWidthN int) {
	max := 0
	for _, text := range c.texts {
		tmp := runewidth.StringWidth(text)
		if tmp > max {
			max = tmp
		}
	}
	runeWidthN = max
	return
}
