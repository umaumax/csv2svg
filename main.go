package main

import (
	"encoding/csv"
	"flag"
	"path/filepath"
	// "fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ajstarks/svgo"
)

func GenarateTable(outputFile string, title string, x int, y int, fontSize int, table [][]string, styles ...string) (canvas *svg.SVG, err error) {
	t := Table{
		X: x, Y: y, FontSize: fontSize,
		Title: title,
	}
	t.SetCells(table)

	file, err := os.Create(outputFile)
	if err != nil {
		return
	}
	defer file.Close()
	canvas = svg.New(file)

	canvas.Start(t.ViewWidth(), t.ViewHeight())
	// NOTE: enable group style
	if len(styles) > 0 {
		canvas.Gstyle(styles[0])
	}
	t.DrawTitleText(canvas)
	t.DrawVerticalTableLine(canvas)
	t.DrawHorizontalTableLine(canvas)
	t.DrawTableText(canvas)

	if len(styles) > 0 {
		canvas.Gend()
	}
	canvas.End()
	return
}

var (
	tsvFlag bool
	// TODO: enable noHeaderFlag
	noHeaderFlag bool
	fontSizeFlag int
	titleFlag    string
	// TODO: add verbose flag
)

func init() {
	flag.BoolVar(&tsvFlag, "tsv", false, "use tsv or not")
	flag.BoolVar(&noHeaderFlag, "no-header", false, "don't use header")
	flag.IntVar(&fontSizeFlag, "font-size", 16, "font size (px)")
	flag.StringVar(&titleFlag, "title", "", "title")
}

func main() {
	flag.Parse()

	// NOTE: default input file is input pipe
	var inputFiles []string
	if flag.NArg() == 0 {
		inputFiles = append(inputFiles, os.Stdin.Name())
	} else {
		inputFiles = append(inputFiles, flag.Args()...)
	}
	for _, inputFile := range inputFiles {
		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		outputFile := os.Stdout.Name()

		ext := filepath.Ext(inputFile)
		if inputFile != os.Stdin.Name() {
			outputFile = inputFile[0:len(inputFile)-len(ext)] + ".svg"
		}

		log.Println(inputFile)
		// TODO: 下記を参考に複数の文字コードに対応?
		// [Go言語でCSVの読み書き\(sjis、euc、utf8対応\) \- Qiita]( https://qiita.com/kesuzuki/items/202cc58db3fd1763c095 )
		reader := csv.NewReader(file)
		// NOTE: ダブルクオートを厳密にチェックしない
		// TODO: write test
		reader.LazyQuotes = true
		if tsvFlag || ext == "tsv" {
			reader.Comma = '\t'
		}
		var table [][]string
		var line []string
		// TODO: 先頭行が，#始まりの場合に，headerとみなす独自拡張?
		// defaultでこれを有効にするが，無効にするflagを設定する
		for {
			line, err = reader.Read()
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			log.Println(line)
			for i, text := range line {
				line[i] = strings.TrimSpace(strings.Replace(text, "\\n", "\n", -1))
			}
			table = append(table, line)
		}
		// TODO: enable x,y flag
		x := 0
		y := 0
		canvas, err := GenarateTable(outputFile, titleFlag, x, y, fontSizeFlag, table, "fill:blue;stroke:black;storke-width:2;")
		if err != nil {
			log.Fatal(err)
		}
		canvas.End()
	}
	return
}
