package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shivan Harvest", NewShivanHarvest)
}

// NewShivanHarvest creates a Shivan Harvest
// {1}{R} - ENCHANTMENT
func NewShivanHarvest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shivan Harvest")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
