package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Anzrag The Quake Mole", NewAnzragTheQuakeMole)
}

// NewAnzragTheQuakeMole creates a Anzrag The Quake Mole
// {2}{R}{G} - CREATURE
func NewAnzragTheQuakeMole(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Anzrag The Quake Mole")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MOLE", "GOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
