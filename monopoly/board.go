package monopoly

import "log"

// Board represents monopoly game board.
type Board struct {
	Body [40]Tile `json:"body"`
}

func (b *Board) Init(j map[string]interface{}) {
	body, ok := j["body"].([]interface{})
	if !ok {
		panic("body is not ok")
	}
	finalBody := [40]Tile{}
	for _, t := range body {
		tile, ok := t.(map[string]interface{})
		if !ok {
			log.Println("tile is not tile")
			continue
		}
		tileType := tile["type"].(string)
		tileIndex := int(tile["id"].(float64))

		if tileType == TileTypeStart ||
			tileType == TileTypeJail ||
			tileType == TileTypeJackpot ||
			tileType == TileTypeArrest ||
			tileType == TileTypeChance ||
			tileType == TileTypeBank {

			finalBody[tileIndex] = BaseTile{
				TType:  tileType,
				TIndex: tileIndex,
			}

		} else if tileType == TileTypeCompany {
			finalBody[tileIndex] = CompanyTile{
				TType:            tileType,
				TIndex:           tileIndex,
				Group:            tile["group"].(string),
				Name:             tile["name"].(string),
				Cost:             int(tile["price"].(float64)),
				UpgradePrice:     0,
				BaseIncome:       0,
				CollectionIncome: 0,
				LevelIncome:      levelIncome{},
			}

		} else if tileType == TileTypeAutomotive {
			finalBody[tileIndex] = AutomotiveTile{
				TType:          tileType,
				TIndex:         tileIndex,
				Group:          tile["group"].(string),
				Name:           tile["name"].(string),
				Cost:           int(tile["price"].(float64)),
				QuantityIncome: quantityIncome{},
			}
		} else if tileType == TileTypeService {
			finalBody[tileIndex] = ServiceTile{
				TType:              tileType,
				TIndex:             tileIndex,
				Group:              tile["group"].(string),
				Name:               tile["name"].(string),
				Cost:               int(tile["price"].(float64)),
				QuantityMultiplier: quantityMultiplier{},
			}
		}
	}

	b.Body = finalBody
}
