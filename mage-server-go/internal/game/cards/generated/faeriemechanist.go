package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Faerie Mechanist", NewFaerieMechanist)
}

// NewFaerieMechanist creates a Faerie Mechanist
// {3}{U} - ARTIFACT CREATURE
// Flying
func NewFaerieMechanist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Faerie Mechanist")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 3, 1, StaticFilters.FILTER_CARD_A...)
	// card.AddAbility(ability1)
	return card, nil
}