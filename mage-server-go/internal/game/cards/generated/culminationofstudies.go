package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Culmination Of Studies", NewCulminationOfStudies)
}

// NewCulminationOfStudies creates a Culmination Of Studies
// {X}{U}{R} - SORCERY
func NewCulminationOfStudies(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Culmination Of Studies")
	card.ManaCost = "{X}{U}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}