package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Treacherous Urge", NewTreacherousUrge)
}

// NewTreacherousUrge creates a Treacherous Urge
// {4}{B} - INSTANT
func NewTreacherousUrge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Treacherous Urge")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("sacrifice " + card.getName(), source.getControlle...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
