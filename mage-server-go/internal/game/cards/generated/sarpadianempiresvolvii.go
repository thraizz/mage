package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sarpadian Empires Vol V I I", NewSarpadianEmpiresVolVII)
}

// NewSarpadianEmpiresVolVII creates a Sarpadian Empires Vol V I I
// {3} - ARTIFACT
func NewSarpadianEmpiresVolVII(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sarpadian Empires Vol V I I")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
