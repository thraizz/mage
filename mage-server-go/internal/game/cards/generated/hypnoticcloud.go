package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hypnotic Cloud", NewHypnoticCloud)
}

// NewHypnoticCloud creates a Hypnotic Cloud
// {1}{B} - SORCERY
func NewHypnoticCloud(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hypnotic Cloud")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(3)
	//
	// Targets:
	//   - abilities.NewPlayerTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
