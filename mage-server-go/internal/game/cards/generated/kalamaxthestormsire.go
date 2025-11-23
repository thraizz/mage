package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kalamax The Stormsire", NewKalamaxTheStormsire)
}

// NewKalamaxTheStormsire creates a Kalamax The Stormsire
// {1}{G}{U}{R} - CREATURE
func NewKalamaxTheStormsire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kalamax The Stormsire")
	card.ManaCost = "{1}{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "DINOSAUR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
