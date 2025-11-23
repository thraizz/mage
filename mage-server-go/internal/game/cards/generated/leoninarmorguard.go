package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leonin Armorguard", NewLeoninArmorguard)
}

// NewLeoninArmorguard creates a Leonin Armorguard
// {2}{G}{W} - CREATURE
func NewLeoninArmorguard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leonin Armorguard")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1, filter, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
