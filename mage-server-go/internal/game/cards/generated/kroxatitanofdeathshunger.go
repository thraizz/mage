package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kroxa Titan Of Deaths Hunger", NewKroxaTitanOfDeathsHunger)
}

// NewKroxaTitanOfDeathsHunger creates a Kroxa Titan Of Deaths Hunger
// {B}{R} - CREATURE
func NewKroxaTitanOfDeathsHunger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kroxa Titan Of Deaths Hunger")
	card.ManaCost = "{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "GIANT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}