package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Seal Of Strength", NewSealOfStrength)
}

// NewSealOfStrength creates a Seal Of Strength
// {G} - ENCHANTMENT
func NewSealOfStrength(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Seal Of Strength")
	card.ManaCost = "{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewBoostEffect(3, 3)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
