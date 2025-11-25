package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haunting Hymn", NewHauntingHymn)
}

// NewHauntingHymn creates a Haunting Hymn
// {4}{B}{B} - INSTANT
func NewHauntingHymn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haunting Hymn")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(4)
	//   - DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
