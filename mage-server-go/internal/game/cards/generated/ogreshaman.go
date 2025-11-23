package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ogre Shaman", NewOgreShaman)
}

// NewOgreShaman creates a Ogre Shaman
// {3}{R}{R} - CREATURE
func NewOgreShaman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ogre Shaman")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddEffect(abilities.NewDamageEffect(2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
