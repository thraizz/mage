package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arvinox The Mind Flail", NewArvinoxTheMindFlail)
}

// NewArvinoxTheMindFlail creates a Arvinox The Mind Flail
// {4}{B}{B}{B} - ENCHANTMENT CREATURE
func NewArvinoxTheMindFlail(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arvinox The Mind Flail")
	card.ManaCost = "{4}{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
