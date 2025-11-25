package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Battlefield Butcher", NewBattlefieldButcher)
}

// NewBattlefieldButcher creates a Battlefield Butcher
// {2}{B} - CREATURE
func NewBattlefieldButcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Battlefield Butcher")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LoseLifeOpponentsEffect()
	//
	// Costs:
	//   - AddManaCost("{5}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
