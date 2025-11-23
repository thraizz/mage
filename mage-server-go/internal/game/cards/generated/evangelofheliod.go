package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Evangel Of Heliod", NewEvangelOfHeliod)
}

// NewEvangelOfHeliod creates a Evangel Of Heliod
// {4}{W}{W} - CREATURE
func NewEvangelOfHeliod(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Evangel Of Heliod")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("SoldierToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, DevotionCount.W)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
