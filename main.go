package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Transaction struct {
	ID          int
	Type        string
	Amount      float64
	Category    string
	Description string
}
type Tracker struct {
	transactions []Transaction
	nextID       int
}

// Добавляем операцию после проверки суммы и категории.
func (tracker *Tracker) add(kind string, amount float64, category, description string) error {
	if amount <= 0 {
		return errors.New("сумма должна быть больше нуля")
	}
	if strings.TrimSpace(category) == "" {
		return errors.New("категория не может быть пустой")
	}
	tracker.nextID = tracker.nextID + 1
	tracker.transactions = append(tracker.transactions, Transaction{tracker.nextID, kind, amount, category, description})
	return nil
}

func (tracker *Tracker) remove(id int) error {
	for index, item := range tracker.transactions {
		if item.ID == id {
			tracker.transactions = append(tracker.transactions[:index], tracker.transactions[index+1:]...)
			return nil
		}
	}
	return errors.New("операция не найдена")
}

// Отдельно считаем доходы и расходы.
func (tracker *Tracker) totals() (float64, float64) {
	income, expense := 0.0, 0.0
	for _, item := range tracker.transactions {
		if item.Type == "доход" {
			income = income + item.Amount
		} else {
			expense = expense + item.Amount
		}
	}
	return income, expense
}

func (tracker *Tracker) expenseByCategory() map[string]float64 {
	result := map[string]float64{}
	for _, item := range tracker.transactions {
		if item.Type == "расход" {
			result[item.Category] = result[item.Category] + item.Amount
		}
	}
	return result
}

func show(items []Transaction, category string) {
	found := false
	for _, item := range items {
		if category == "" || strings.EqualFold(item.Category, category) {
			fmt.Printf("#%d | %s | %s | %.2f тг | %s\n", item.ID, item.Type, item.Category, item.Amount, item.Description)
			found = true
		}
	}
	if !found {
		fmt.Println("Операций нет")
	}
}

var input = bufio.NewScanner(os.Stdin)

func read(prompt string) string {
	fmt.Print(prompt)
	if !input.Scan() {
		return "0"
	}
	return strings.TrimSpace(input.Text())
}

func main() {
	tracker := Tracker{}
	for {
		fmt.Println()
		fmt.Println("1. Добавить доход\n2. Добавить расход\n3. История операций\n4. Текущий баланс\n5. Расходы по категориям\n6. Поиск по категории\n7. Удалить операцию\n0. Выход")
		choice := read("Выберите действие: ")
		switch choice {
		case "1", "2":
			kind := "доход"
			if choice == "2" {
				kind = "расход"
			}
			amount, inputError := strconv.ParseFloat(strings.ReplaceAll(read("Сумма: "), ",", "."), 64)
			if inputError != nil {
				fmt.Println("Неверная сумма")
				continue
			}
			category := read("Категория: ")
			description := read("Описание: ")
			if inputError = tracker.add(kind, amount, category, description); inputError != nil {
				fmt.Println(inputError)
			} else {
				fmt.Println("Операция добавлена")
			}
		case "3":
			show(tracker.transactions, "")
		case "4":
			income, expense := tracker.totals()
			fmt.Printf("Доходы: %.2f тг\nРасходы: %.2f тг\nБаланс: %.2f тг\n", income, expense, income-expense)
		case "5":
			for category, amount := range tracker.expenseByCategory() {
				fmt.Printf("%s: %.2f тг\n", category, amount)
			}
		case "6":
			show(tracker.transactions, read("Категория: "))
		case "7":
			id, inputError := strconv.Atoi(read("ID операции: "))
			if inputError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			if inputError = tracker.remove(id); inputError != nil {
				fmt.Println(inputError)
			} else {
				fmt.Println("Операция удалена")
			}
		case "0":
			return
		default:
			fmt.Println("Неизвестный пункт")
		}
	}
}
