package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Burning Tree Bloodscale", NewBurningTreeBloodscale)
}

// NewBurningTreeBloodscale creates a Burning Tree Bloodscale
// {2}{R}{G} - CREATURE
func NewBurningTreeBloodscale(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Burning Tree Bloodscale")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CantBeBlockedByTargetSourceEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - MustBeBlockedByTargetSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
