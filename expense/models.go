package expense

type Expense struct {
	Type        string
	Name        string
	Value       string
	Description string
}

type Expenser interface {
	Create(e Expense, db Database) error
	Read(id string) (bool, error)
	Update(id string) (bool, error)
	Delete(id string)
}
