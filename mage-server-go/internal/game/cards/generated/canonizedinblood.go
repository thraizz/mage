package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Canonized In Blood", NewCanonizedInBlood)
}

// NewCanonizedInBlood creates a Canonized In Blood
// {1}{B} - ENCHANTMENT
func NewCanonizedInBlood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Canonized In Blood")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("VampireDemonToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
