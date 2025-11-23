package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Claim Fame", NewClaimFame)
}

// NewClaimFame creates a Claim Fame
// {B} - SORCERY
func NewClaimFame(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Claim Fame")
	card.ManaCost = "{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//
	// Targets:
	//   - abilities.NewCreatureTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
