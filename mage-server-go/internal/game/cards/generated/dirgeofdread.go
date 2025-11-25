package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Dirge Of Dread", NewDirgeOfDread)
}

// NewDirgeOfDread creates a Dirge Of Dread
// {2}{B} - SORCERY
func NewDirgeOfDread(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dirge Of Dread")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("FearAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect("FearAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect("FearAbility", effects.DurationEndOfTurn)).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
