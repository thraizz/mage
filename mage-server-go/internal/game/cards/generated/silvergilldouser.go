package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Silvergill Douser", NewSilvergillDouser)
}

// NewSilvergillDouser creates a Silvergill Douser
// {1}{U} - CREATURE
func NewSilvergillDouser(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Silvergill Douser")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewBoostEffect(xValue, StaticValue.get(0))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
