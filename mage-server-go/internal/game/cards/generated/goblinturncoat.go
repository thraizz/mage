package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goblin Turncoat", NewGoblinTurncoat)
}

// NewGoblinTurncoat creates a Goblin Turncoat
// {1}{B} - CREATURE
func NewGoblinTurncoat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goblin Turncoat")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "MERCENARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
