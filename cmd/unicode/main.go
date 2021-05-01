//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	for _, arg := range flag.Args() {
		err := processFile(arg)
		if err != nil {
			log.Fatalf("%s: %s", arg, err)
		}
	}
}

type Region struct {
	maxWidth int
	lines    [][]rune
}

func (r *Region) String() string {
	var b strings.Builder

	for row, line := range r.lines {
		if row > 0 {
			b.WriteRune('\n')
		}
		for _, r := range line {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (r *Region) Width() int {
	return r.maxWidth
}

func (r *Region) Height() int {
	return len(r.lines)
}

func (r *Region) Get(row, col int) rune {
	if row < 0 || row >= len(r.lines) {
		return 0
	}
	if col < 0 || col >= len(r.lines[row]) {
		return 0
	}
	return r.lines[row][col]
}

func (r *Region) Set(row, col int, ch rune) {
	if row < 0 || row >= len(r.lines) {
		return
	}
	if col < 0 || col >= len(r.lines[row]) {
		return
	}
	r.lines[row][col] = ch
}

func NewRegion(input []byte) *Region {
	var lines [][]rune
	var width int

	for _, line := range strings.Split(string(input), "\n") {
		runes := []rune(line)
		if len(runes) > width {
			width = len(runes)
		}
		lines = append(lines, runes)
	}
	return &Region{
		maxWidth: width,
		lines:    lines,
	}
}

const (
	FlagUp int = 1 << iota
	FlagDown
	FlagLeft
	FlagRight
)

var lineDrawing = []rune{
	'+',    //
	0x2576, // 	               Up
	0x2577, // 	          Down
	0x2502, // 	          Down Up
	0x2574, // 	     Left
	0x2518, // 	     Left      Up
	0x2510, // 	     Left Down
	0x2524, // 	     Left Down Up
	0x2576, // Right
	0x2514, // Right           Up
	0x250C, // Right      Down
	0x251C, // Right      Down Up
	0x2500, // Right Left
	0x2534, // Right Left      Up
	0x252C, // Right Left Down
	0x253C, // Right Left Down Up
}

var lineDrawingRound = []rune{
	'+',    //
	0x2576, // 	               Up
	0x2577, // 	          Down
	0x2502, // 	          Down Up
	0x2574, // 	     Left
	0x256F, // 	     Left      Up
	0x256E, // 	     Left Down
	0x2524, // 	     Left Down Up
	0x2576, // Right
	0x2570, // Right           Up
	0x256D, // Right      Down
	0x251C, // Right      Down Up
	0x2500, // Right Left
	0x2534, // Right Left      Up
	0x252C, // Right Left Down
	0x253C, // Right Left Down Up
}

func isLine(r rune) bool {
	switch r {
	case '|', '-', '+', '*':
		return true

	default:
		return 0x2500 <= r && r <= 0x257F
	}
}

func processFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	region := NewRegion(data)
	for row := 0; row < region.Height(); row++ {
		for col := 0; col < region.Width(); col++ {
			ch := region.Get(row, col)
			switch ch {
			case '|':
				region.Set(row, col, 0x2502)
			case '-':
				region.Set(row, col, 0x2500)
			case '+', '*':
				var index int
				if isLine(region.Get(row-1, col)) {
					index |= FlagUp
				}
				if isLine(region.Get(row+1, col)) {
					index |= FlagDown
				}
				if isLine(region.Get(row, col-1)) {
					index |= FlagLeft
				}
				if isLine(region.Get(row, col+1)) {
					index |= FlagRight
				}
				if ch == '+' {
					region.Set(row, col, lineDrawing[index])
				} else {
					region.Set(row, col, lineDrawingRound[index])
				}
			}
		}
	}

	fmt.Print(region.String())

	return nil
}
