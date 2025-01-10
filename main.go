package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Knetic/govaluate"
)

var entry *widget.Entry
var filter = []string{"+", "-", "*", "/", "%", "."}
var lastInputIsFlag bool

func main() {
	a := app.New()
	w := a.NewWindow("Calculator")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(450, 300))

	entry = widget.NewEntry()
	entry.MultiLine = true

	digits := []string{
		"7", "8", "9", "*",
		"4", "5", "6", "-",
		"1", "2", "3", "+",
	}
	var digitBtns []fyne.CanvasObject
	digitBtns = append(digitBtns, widget.NewButton("c", func() {
		entry.SetText("")
		entry.Refresh()
	}))
	digitBtns = append(digitBtns, widget.NewButton("+/-", sign()))
	digitBtns = append(digitBtns, widget.NewButton("()", input("()")))
	digitBtns = append(digitBtns, widget.NewButton("/", input("/")))

	for _, v := range digits {
		val := v
		digitBtns = append(digitBtns, widget.NewButton(val, input(val)))
	}

	buts := container.New(layout.NewGridLayout(4), digitBtns...)

	equal := widget.NewButton("=", equals())

	b := container.New(
		layout.NewGridLayout(2),
		widget.NewButton("0", input("0")),
		container.New(layout.NewGridLayout(2), widget.NewButton(".", input(".")), equal))

	w.SetContent(container.New(layout.NewVBoxLayout(), entry, buts, b))
	w.ShowAndRun()
}

func percent() func() {
	return func() {

	}
}

func sign() func() {
	return func() {
		if strings.Contains(entry.Text, ".") {
			value, err := strconv.ParseFloat(entry.Text, 64)
			if err != nil {
				return
			}
			value = -value
			entry.Text = fmt.Sprint(value)
		} else {
			value, err := strconv.ParseInt(entry.Text, 10, 64)
			if err != nil {
				return
			}
			value = -value
			entry.Text = fmt.Sprint(value)
		}

		entry.Refresh()
	}
}

func equals() func() {
	return func() {
		lines := strings.Split(entry.Text, "\n")
		if len(lines) == 0 {
			return
		}

		line := lines[len(lines)-1]
		line = strings.Trim(line, "+/x")
		expr, _ := govaluate.NewEvaluableExpression(line)
		result, _ := expr.Evaluate(nil)
		line += "=\n"
		line += fmt.Sprint(result)
		entry.Text = line
		entry.Refresh()
	}
}

func input(val string) func() {
	return func() {
		var thisFlag bool
		for _, v := range filter {
			if v == val {
				thisFlag = true
			}
		}

		if thisFlag && lastInputIsFlag {
			return
		}

		lastInputIsFlag = thisFlag
		entry.SetText(entry.Text + val)
		entry.Refresh()
	}
}
