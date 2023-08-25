package monopoly

const (
	TileTypeStart      = "start"
	TileTypeJail       = "jail"
	TileTypeJackpot    = "jackpot"
	TileTypeArrest     = "arrest"
	TileTypeChance     = "chance"
	TileTypeBank       = "bank"
	TileTypeCompany    = "company"
	TileTypeCash       = "cash"
	TileTypeService    = "service"
	TileTypeAutomotive = "automotive"
)

// Tile represents 1 step on the Board.
type Tile interface {
	Type() string
	Index() int
}

type ObtainableTile interface {
	Tile
	Price() int
}

type BaseTile struct {
	TType  string `json:"type"`
	TIndex int    `json:"index"`
}

func (t BaseTile) Type() string {
	return t.TType
}

func (t BaseTile) Index() int {
	return t.TIndex
}

type CompanyTile struct {
	TType            string      `json:"type"`
	TIndex           int         `json:"index"`
	Group            string      `json:"group"`
	Name             string      `json:"name"`
	Cost             int         `json:"price"`
	UpgradePrice     int         `json:"upgradePrice"`
	BaseIncome       int         `json:"baseIncome"`
	CollectionIncome int         `json:"collectionIncome"`
	LevelIncome      levelIncome `json:"levelIncome"`
}

func (t CompanyTile) Type() string {
	return t.TType
}

func (t CompanyTile) Index() int {
	return t.TIndex
}

func (t CompanyTile) Price() int {
	return t.Cost
}

type AutomotiveTile struct {
	TType          string         `json:"type"`
	TIndex         int            `json:"index"`
	Group          string         `json:"group"`
	Name           string         `json:"name"`
	Cost           int            `json:"price"`
	QuantityIncome quantityIncome `json:"quantityIncome"`
}

func (t AutomotiveTile) Type() string {
	return t.TType
}

func (t AutomotiveTile) Index() int {
	return t.TIndex
}

func (t AutomotiveTile) Price() int {
	return t.Cost
}

type ServiceTile struct {
	TType              string             `json:"type"`
	TIndex             int                `json:"index"`
	Group              string             `json:"group"`
	Name               string             `json:"name"`
	Cost               int                `json:"price"`
	QuantityMultiplier quantityMultiplier `json:"quantityMultiplier"`
}

func (t ServiceTile) Type() string {
	return t.TType
}

func (t ServiceTile) Index() int {
	return t.TIndex
}
func (t ServiceTile) Price() int {
	return t.Cost
}

type quantityMultiplier struct {
	First  int `json:"1"`
	Second int `json:"2"`
}

type quantityIncome struct {
	First  int `json:"1"`
	Second int `json:"2"`
	Third  int `json:"3"`
	Fourth int `json:"4"`
}

type levelIncome struct {
	First  int `json:"1"`
	Second int `json:"2"`
	Third  int `json:"3"`
	Fourth int `json:"4"`
	Five   int `json:"5"`
}
