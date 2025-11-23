package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gideons Battle Cry", NewGideonsBattleCry)
}

// NewGideonsBattleCry creates a Gideons Battle Cry
// {2}{W}{W} - SORCERY
func NewGideonsBattleCry(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gideons Battle Cry")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
