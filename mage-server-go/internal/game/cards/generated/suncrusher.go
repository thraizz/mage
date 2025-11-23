package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Suncrusher", NewSuncrusher)
}

// NewSuncrusher creates a Suncrusher
// {9} - ARTIFACT CREATURE
func NewSuncrusher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Suncrusher")
	card.ManaCost = "{9}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddEffect(abilities.NewReturnToHandSourceEffect()).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
