package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Revival Revenge", NewRevivalRevenge)
}

// NewRevivalRevenge creates a Revival Revenge
// {W/B}{W/B} - SORCERY
func NewRevivalRevenge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Revival Revenge")
	card.ManaCost = "{W/B}{W/B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
