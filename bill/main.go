package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

func createBill() bill {
	reader := bufio.NewReader(os.Stdin)

	name, _ := getInput("Create a new bill name: ", reader)

	b := newBill(name)
	fmt.Println("Create the bill - ", b.name)

	return b
}

func promptOptions(b *bill) {
	reader := bufio.NewReader((os.Stdin))

	opt, _ := getInput("Choose option (a - add item, s - save bill, t -add tips): ", reader)
	switch opt {
	case "a":
		name, _ := getInput("Enter the name of item: ", reader)
		price, _ := getInput("Enter the price of item: ", reader)
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("The price must be a number")
		} else {
			b.addItem(name, priceFloat)
			fmt.Printf("You add %v to the bill.\n", name)
		}
		promptOptions(b)
	case "s":
		b.save()
		fmt.Println("You save the bill: ", b.name)
	case "t":
		tips, _ := getInput("Set the tips: ", reader)
		tipsfloats, err := strconv.ParseFloat(tips, 64)
		if err != nil {
			fmt.Println("The tip must be a number")
		} else {
			b.updateTip(tipsfloats)
			fmt.Printf("You set the tips : %.02f\n", tipsfloats)
		}
		promptOptions(b)
	default:
		fmt.Println("invalid option...")
		promptOptions(b)
	}

}

func main() {
	mybill := createBill()

	promptOptions(&mybill)
	fs := mybill.format()
	fmt.Println(fs)
	//fmt.Println(mybill)
}
