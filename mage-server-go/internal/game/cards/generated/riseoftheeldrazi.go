package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rise Of The Eldrazi", NewRiseOfTheEldrazi)
}

// NewRiseOfTheEldrazi creates a Rise Of The Eldrazi
// {9}{C}{C}{C} - SORCERY
func NewRiseOfTheEldrazi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rise Of The Eldrazi")
	card.ManaCost = "{9}{C}{C}{C}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
