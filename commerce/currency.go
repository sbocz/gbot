package commerce

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Currency int

const currencySymbol = "𝔻"

var (
	printer = message.NewPrinter(language.English)
)

func (c Currency) String() string {
	return printer.Sprintf("%d%s", int(c), currencySymbol)
}
