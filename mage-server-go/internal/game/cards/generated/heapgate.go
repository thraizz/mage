package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Heap Gate", NewHeapGate)
}

// NewHeapGate creates a Heap Gate
//   - LAND
func NewHeapGate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Heap Gate")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"GATE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
