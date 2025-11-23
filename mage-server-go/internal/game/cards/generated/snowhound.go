package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Snow Hound", NewSnowHound)
}

// NewSnowHound creates a Snow Hound
// {2}{W} - CREATURE
func NewSnowHound(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Snow Hound")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DOG"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddEffect(abilities.NewReturnToHandSourceEffect()).
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
