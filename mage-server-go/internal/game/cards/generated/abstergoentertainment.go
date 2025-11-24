package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abstergo Entertainment", NewAbstergoEntertainment)
}

// NewAbstergoEntertainment creates a Abstergo Entertainment
//  - LAND
func NewAbstergoEntertainment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abstergo Entertainment")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect()
	//
	// Costs:
	//   - AddManaCost("{3}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}