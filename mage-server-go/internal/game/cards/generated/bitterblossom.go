package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bitterblossom", NewBitterblossom)
}

// NewBitterblossom creates a Bitterblossom
// {1}{B} - KINDRED ENCHANTMENT
func NewBitterblossom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bitterblossom")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"FAERIE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
