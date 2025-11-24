package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Tawnos Solemn Survivor", NewTawnosSolemnSurvivor)
}

// NewTawnosSolemnSurvivor creates a Tawnos Solemn Survivor
// {1}{U} - CREATURE
func NewTawnosSolemnSurvivor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tawnos Solemn Survivor")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}