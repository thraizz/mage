package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dreamshaper Shaman", NewDreamshaperShaman)
}

// NewDreamshaperShaman creates a Dreamshaper Shaman
// {5}{R} - ENCHANTMENT CREATURE
func NewDreamshaperShaman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dreamshaper Shaman")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"MINOTAUR", "SHAMAN"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new RevealCardsFromLibraryUntilEf...)
	// card.AddAbility(ability0)
	return card, nil
}
