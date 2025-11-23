package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Wookiee Raidleader", NewWookieeRaidleader)
}

// NewWookieeRaidleader creates a Wookiee Raidleader
// {2}{R/G}{R/G} - CREATURE
func NewWookieeRaidleader(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wookiee Raidleader")
	card.ManaCost = "{2}{R/G}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WOOKIEE", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
