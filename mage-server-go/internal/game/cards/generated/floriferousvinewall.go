package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Floriferous Vinewall", NewFloriferousVinewall)
}

// NewFloriferousVinewall creates a Floriferous Vinewall
// {1}{G} - CREATURE
// Defender
func NewFloriferousVinewall(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Floriferous Vinewall")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "0"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 6, 1, StaticFilters.FILTER_CARD_L...)
	// card.AddAbility(ability1)
	return card, nil
}