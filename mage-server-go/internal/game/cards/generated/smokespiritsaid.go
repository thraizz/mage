package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Smoke Spirits Aid", NewSmokeSpiritsAid)
}

// NewSmokeSpiritsAid creates a Smoke Spirits Aid
// {X}{R} - SORCERY
func NewSmokeSpiritsAid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Smoke Spirits Aid")
	card.ManaCost = "{X}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
