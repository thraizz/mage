package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kami Of Twisted Reflection", NewKamiOfTwistedReflection)
}

// NewKamiOfTwistedReflection creates a Kami Of Twisted Reflection
// {1}{U}{U} - CREATURE
func NewKamiOfTwistedReflection(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kami Of Twisted Reflection")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
