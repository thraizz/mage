package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tears Of Rage", NewTearsOfRage)
}

// NewTearsOfRage creates a Tears Of Rage
// {2}{R}{R} - INSTANT
func NewTearsOfRage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tears Of Rage")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("Sacrifice those creatures at the beginning of the...)
	// card.AddAbility(ability0)
	return card, nil
}
