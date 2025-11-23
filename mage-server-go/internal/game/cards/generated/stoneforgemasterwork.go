package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stoneforge Masterwork", NewStoneforgeMasterwork)
}

// NewStoneforgeMasterwork creates a Stoneforge Masterwork
// {1} - ARTIFACT
func NewStoneforgeMasterwork(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stoneforge Masterwork")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(countEnchantments, countEnchantments)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
