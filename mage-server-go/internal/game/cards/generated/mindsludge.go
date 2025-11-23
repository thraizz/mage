package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mind Sludge", NewMindSludge)
}

// NewMindSludge creates a Mind Sludge
// {4}{B} - SORCERY
func NewMindSludge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mind Sludge")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(new PermanentsOnBattlefieldCount(filter))
	//
	// Targets:
	//   - abilities.NewPlayerTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
