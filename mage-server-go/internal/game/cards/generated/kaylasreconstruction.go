package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaylas Reconstruction", NewKaylasReconstruction)
}

// NewKaylasReconstruction creates a Kaylas Reconstruction
// {X}{W}{W}{W} - SORCERY
func NewKaylasReconstruction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaylas Reconstruction")
	card.ManaCost = "{X}{W}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}