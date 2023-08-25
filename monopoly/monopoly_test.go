package monopoly

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"
)

func TestTable(t *testing.T) {

	var rawMap map[string]interface{}
	f, err := os.ReadFile("map.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(f, &rawMap)
	if err != nil {
		t.Fatal(err)
	}

	b := &Board{}

	b.Init(rawMap)

	table := NewTable(b)

	nick := NewPlayer("1235", "Nick")

	go func() {
		for {
			select {
			case event := <-nick.Updates():
				log.Printf("%s got event %v", nick.Name, event)
				switch event.EventType {
				case EventTypePlayerJoined:
					eData, ok := event.Data.(EventPlayerJoined)
					if !ok {
						t.Errorf("%s got wrong event data", EventTypePlayerJoined)
					}
					log.Println("Player joined", eData)
				}
			}
		}
	}()

	table.AddPlayer(nick)
	table.AddPlayer(NewPlayer("123576", "Lala"))

	table.Start()
	time.Sleep(20 * time.Second)
}
