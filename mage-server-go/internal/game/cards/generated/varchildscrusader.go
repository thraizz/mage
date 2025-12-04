package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Varchilds Crusader", NewVarchildsCrusader)
}

// NewVarchildsCrusader creates a Varchilds Crusader
// {3}{R} - CREATURE
func NewVarchildsCrusader(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Varchilds Crusader")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CantBeBlockedByCreaturesSourceEffect()
	//   - SacrificeSourceEffect()
	//
	// Costs:
	//   - AddManaCost("{0}")
	// card.AddAbility(ability0)
	return card, nil
}
