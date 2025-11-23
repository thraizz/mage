package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Down Dirty", NewDownDirty)
}

// NewDownDirty creates a Down Dirty
// {3}{B} - SORCERY
func NewDownDirty(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Down Dirty")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewPlayerTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
