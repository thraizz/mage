package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Amber Gristle O Maul", NewAmberGristleOMaul)
}

// NewAmberGristleOMaul creates a Amber Gristle O Maul
// {3}{R} - CREATURE
// Haste
func NewAmberGristleOMaul(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Amber Gristle O Maul")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DWARF", "CLERIC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new DrawCardSourceControllerEffec...)
	// card.AddAbility(ability1)
	return card, nil
}
