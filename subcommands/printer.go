package subcommands

import (
	"fmt"
)

const (
	colorReset = "\033[0m"

	colorError   = "\033[31m"
	colorSuccess = "\033[32m"
	colorUp      = "\033[36m"
	colorDown    = "\033[35m"
)

type Printer interface {
	PrintUpMigration(text string)
	PrintDownMigration(text string)
	PrintError(text string)
	PrintSuccess(text string)

	SetNoColor(color bool)
}

type ImplPrinter struct {
	NoColor bool
}

func (p *ImplPrinter) PrintUpMigration(text string) {
	if p.NoColor {
		fmt.Println(text)
		return
	}

	fmt.Println(colorUp, "⏫ ", text, colorReset)
}

func (p *ImplPrinter) PrintDownMigration(text string) {
	if p.NoColor {
		fmt.Println(text)
		return
	}

	fmt.Println(colorDown, "⏬ ", text, colorReset)
}

func (p *ImplPrinter) PrintError(text string) {
	if p.NoColor {
		fmt.Println(text)
		return
	}

	fmt.Println(colorError, "❌ ", text, colorReset)

}

func (p *ImplPrinter) PrintSuccess(text string) {
	if p.NoColor {
		fmt.Println(text)
		return
	}

	fmt.Println(colorSuccess, "✔️ ", text, colorReset)
}

func (p *ImplPrinter) SetNoColor(color bool) {
	p.NoColor = color
}
