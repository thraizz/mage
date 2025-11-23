package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Proven Combatant", NewProvenCombatant)
}

// NewProvenCombatant creates a Proven Combatant
// {U} - CREATURE
func NewProvenCombatant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Proven Combatant")
	card.ManaCost = "{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
