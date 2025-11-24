package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aggressive Biomancy", NewAggressiveBiomancy)
}

// NewAggressiveBiomancy creates a Aggressive Biomancy
// {X}{X}{G}{U} - SORCERY
func NewAggressiveBiomancy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aggressive Biomancy")
	card.ManaCost = "{X}{X}{G}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
