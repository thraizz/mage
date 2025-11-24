package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Lorehold Excavation", NewLoreholdExcavation)
}

// NewLoreholdExcavation creates a Lorehold Excavation
// {R}{W} - ENCHANTMENT
func NewLoreholdExcavation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lorehold Excavation")
	card.ManaCost = "{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("Spirit32Token")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 1, true, false)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}