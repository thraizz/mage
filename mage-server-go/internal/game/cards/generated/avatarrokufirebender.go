package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Avatar Roku Firebender", NewAvatarRokuFirebender)
}

// NewAvatarRokuFirebender creates a Avatar Roku Firebender
// {3}{R}{R}{R} - CREATURE
func NewAvatarRokuFirebender(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Avatar Roku Firebender")
	card.ManaCost = "{3}{R}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(3, 0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
