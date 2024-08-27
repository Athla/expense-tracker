package main

import "github.com/charmbracelet/bubbles/table"

type AddModel struct{}
type SeeModel struct{}
type ManageModel struct{}

type Parcel struct {
	Expense  Expense
	Parcelas int
}

var columns = []table.Column{
	{Title: "Name", Width: 10},
	{Title: "Value", Width: 10},
	{Title: "Tag", Width: 15},
	{Title: "Type", Width: 10},
}
var rows = []table.Row{
	{"Netflix", "55.90", "Entretenimento", "Mensal"},
	{"Gympass", "49.90", "Entretenimento", "Mensal"},
	{"Espetinho", "35.00", "Comida", "Avulso"},
}
var expenses = []Expense{
	{Name: "Netflix", Value: 55.90, Tag: "Entretenimento", Type: "Mensal"},
	{Name: "Gympass", Value: 49.90, Tag: "Entretenimento", Type: "Mensal"},
	{Name: "Espetinho", Value: 35.00, Tag: "Comida", Type: "Avulso"},
}
