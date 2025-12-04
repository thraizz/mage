package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Treacherous Terrain", NewTreacherousTerrain)
}

// NewTreacherousTerrain creates a Treacherous Terrain
// {6}{R}{G} - SORCERY
func NewTreacherousTerrain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Treacherous Terrain")
	card.ManaCost = "{6}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
