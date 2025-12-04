package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirror Of Life Trapping", NewMirrorOfLifeTrapping)
}

// NewMirrorOfLifeTrapping creates a Mirror Of Life Trapping
// {4} - ARTIFACT
func NewMirrorOfLifeTrapping(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirror Of Life Trapping")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
