package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warren Weirding", NewWarrenWeirding)
}

// NewWarrenWeirding creates a Warren Weirding
// {1}{B} - KINDRED SORCERY
func NewWarrenWeirding(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warren Weirding")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"GOBLIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
