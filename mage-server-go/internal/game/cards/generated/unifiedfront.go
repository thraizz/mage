package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Unified Front", NewUnifiedFront)
}

// NewUnifiedFront creates a Unified Front
// {3}{W} - SORCERY
func NewUnifiedFront(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unified Front")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("KorAllyToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, ColorsOfManaSpentToCastCount.getInstance())).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, ColorsOfManaSpentToCastCount.getInstance())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
