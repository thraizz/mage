package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chalice Of Death", NewChaliceOfDeath)
}

// NewChaliceOfDeath creates a Chalice Of Death
//   - ARTIFACT
func NewChaliceOfDeath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chalice Of Death")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewLoseLifeEffect(5)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
