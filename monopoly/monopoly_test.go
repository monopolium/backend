package monopoly

import (
	"testing"
)

func TestTable(t *testing.T) {
	table := &Table{
		Players: Players{},
		Board:   &Board{},
	}

	table.AddPlayer(NewPlayer("1235", "Nick"))
	table.AddPlayer(NewPlayer("123576", "Lala"))
}