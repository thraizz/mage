package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fodder Launch", NewFodderLaunch)
}

// NewFodderLaunch creates a Fodder Launch
// {3}{B} - KINDRED SORCERY
func NewFodderLaunch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fodder Launch")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"GOBLIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-5, -5)).
		AddEffect(abilities.NewDamageEffect(5)).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
