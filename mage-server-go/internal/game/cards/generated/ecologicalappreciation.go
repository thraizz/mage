package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ecological Appreciation", NewEcologicalAppreciation)
}

// NewEcologicalAppreciation creates a Ecological Appreciation
// {X}{2}{G} - SORCERY
func NewEcologicalAppreciation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ecological Appreciation")
	card.ManaCost = "{X}{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}