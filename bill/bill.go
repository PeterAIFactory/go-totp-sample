package main

import (
	"fmt"
	"os"
)

type bill struct {
	name  string
	items map[string]float64
	tip   float64
}

// make new bills
func newBill(name string) bill {
	b := bill{
		name:  name,
		items: map[string]float64{},
		tip:   0,
	}
	return b
}

func (b bill) format() string {
	fs := "Bill breakdown: \n"
	var total float64 = 0

	for k, v := range b.items {
		fs += fmt.Sprintf("%-25v ...$%.02f\n", k+":", v)
		total += v
	}
	fmt.Println("Tips: ", b.tip)
	fs += fmt.Sprintf("%-25v ...$%.02f\n", "Tip:", b.tip)

	fs += fmt.Sprintf("%-25v ...$%.02f", "Total:", total+b.tip)

	return fs
}

func (b *bill) updateTip(tip float64) {
	b.tip = tip
	// fmt.Println(b.tip)
}

func (b *bill) addItem(name string, price float64) {
	b.items[name] = price
}

func (b *bill) save() {
	data := []byte(b.format())
	err := os.WriteFile("test/bills/"+b.name+".txt", data, 0644)
	if err != nil {
		panic(err)
	}
}
