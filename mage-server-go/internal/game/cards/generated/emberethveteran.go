package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Embereth Veteran", NewEmberethVeteran)
}

// NewEmberethVeteran creates a Embereth Veteran
// {R} - CREATURE
func NewEmberethVeteran(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Embereth Veteran")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "KNIGHT"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CreateRoleAttachedTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}
