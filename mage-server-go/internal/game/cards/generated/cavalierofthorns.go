package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cavalier Of Thorns", NewCavalierOfThorns)
}

// NewCavalierOfThorns creates a Cavalier Of Thorns
// {2}{G}{G}{G} - CREATURE
// Reach
func NewCavalierOfThorns(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cavalier Of Thorns")
	card.ManaCost = "{2}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "KNIGHT"}
	card.Power = "5"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                 5, 1, StaticFilters.FILTER_CARD_L...)
	// card.AddAbility(ability1)
	return card, nil
}
