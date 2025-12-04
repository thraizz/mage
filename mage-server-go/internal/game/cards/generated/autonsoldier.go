package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Auton Soldier", NewAutonSoldier)
}

// NewAutonSoldier creates a Auton Soldier
// {4}{U}{U} - ARTIFACT CREATURE
func NewAutonSoldier(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Auton Soldier")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ALIEN", "SOLDIER"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                         StaticFilters.FILTER_PERM...)
	// card.AddAbility(ability0)
	return card, nil
}
