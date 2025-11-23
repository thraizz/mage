package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unnatural Growth", NewUnnaturalGrowth)
}

// NewUnnaturalGrowth creates a Unnatural Growth
// {1}{G}{G}{G}{G} - ENCHANTMENT
func NewUnnaturalGrowth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unnatural Growth")
	card.ManaCost = "{1}{G}{G}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(permanent.getPower().getValue(), permanent.getToughness().getValue())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
