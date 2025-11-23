package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shagrat Loot Bearer", NewShagratLootBearer)
}

// NewShagratLootBearer creates a Shagrat Loot Bearer
// {2}{B}{R} - CREATURE
func NewShagratLootBearer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shagrat Loot Bearer")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORC", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
