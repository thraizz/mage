package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Song Of The Worldsoul", NewSongOfTheWorldsoul)
}

// NewSongOfTheWorldsoul creates a Song Of The Worldsoul
// {4}{W}{W} - ENCHANTMENT
func NewSongOfTheWorldsoul(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Song Of The Worldsoul")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
