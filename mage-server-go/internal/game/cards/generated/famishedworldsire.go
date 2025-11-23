package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Famished Worldsire", NewFamishedWorldsire)
}

// NewFamishedWorldsire creates a Famished Worldsire
// {5}{G}{G}{G} - CREATURE
func NewFamishedWorldsire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Famished Worldsire")
	card.ManaCost = "{5}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LEVIATHAN"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
