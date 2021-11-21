package main

import "github.com/hidetatz/skv"

func main() {
	// tl := skv.NewTransactionLog()
	p := skv.NewPrimary(3000)
	p.Start()
}
