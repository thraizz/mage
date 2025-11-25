package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Morsel Theft", NewMorselTheft)
}

// NewMorselTheft creates a Morsel Theft
// {2}{B}{B} - KINDRED SORCERY
func NewMorselTheft(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Morsel Theft")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"ROGUE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(3)).
		AddEffect(abilities.NewLoseLifeEffect(3)).
		AddEffect(abilities.NewGainLifeEffect(3)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
