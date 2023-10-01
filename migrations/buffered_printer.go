package migrations

import (
	"fmt"
	"strings"
)

type bufferedPrinter struct {
	builder strings.Builder
}

func (b *bufferedPrinter) PrintUpMigration(text string) {
	b.builder.WriteString("UP MIGRATION:" + text + "\n")
}

func (b *bufferedPrinter) PrintDownMigration(text string) {
	b.builder.WriteString("DOWN MIGRATION:" + text + "\n")
}

func (b *bufferedPrinter) PrintError(text string) {
	b.builder.WriteString("ERROR:" + text + "\n")
}

func (b *bufferedPrinter) PrintSuccess(text string) {
	b.builder.WriteString("SUCCESS:" + text + "\n")
}

func (b *bufferedPrinter) PrintMigrations(date string, onFS string, inDB string) {
	var db, fs string

	if onFS != "" {
		fs = fmt.Sprintf("fs:%s", onFS)
	}

	if inDB != "" {
		db = fmt.Sprintf("db: %s", inDB)
	}

	result := fmt.Sprintf("%s   |   %s / %s", date, fs, db)
	b.builder.WriteString("ALL MIGRATIONS:" + result + "\n")
}

func (b *bufferedPrinter) SetNoColor(color bool) {}

func (b *bufferedPrinter) GetAllPrints() string {
	return b.builder.String()
}

func newBufferedPrinter() *bufferedPrinter {
	return &bufferedPrinter{}
}
