package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Saheeli The Suns Brilliance", NewSaheeliTheSunsBrilliance)
}

// NewSaheeliTheSunsBrilliance creates a Saheeli The Suns Brilliance
// {U}{R} - CREATURE
func NewSaheeliTheSunsBrilliance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Saheeli The Suns Brilliance")
	card.ManaCost = "{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
