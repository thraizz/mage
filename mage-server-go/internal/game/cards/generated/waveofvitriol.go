package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wave Of Vitriol", NewWaveOfVitriol)
}

// NewWaveOfVitriol creates a Wave Of Vitriol
// {5}{G}{G} - SORCERY
func NewWaveOfVitriol(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wave Of Vitriol")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
