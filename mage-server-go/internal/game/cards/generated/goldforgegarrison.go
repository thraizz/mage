package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Gold Forge Garrison", NewGoldForgeGarrison)
}

// NewGoldForgeGarrison creates a Gold Forge Garrison
//   - LAND
func NewGoldForgeGarrison(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gold Forge Garrison")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("GoldForgeGarrisonGolemToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
