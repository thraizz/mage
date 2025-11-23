package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Totem", NewPhyrexianTotem)
}

// NewPhyrexianTotem creates a Phyrexian Totem
// {3} - ARTIFACT
// Trample
func NewPhyrexianTotem(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Totem")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"PHYREXIAN", "HORROR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(StaticFilters.FILTER_PERMANENTS, SavedDamageValue....)
	// card.AddAbility(ability2)
	return card, nil
}
