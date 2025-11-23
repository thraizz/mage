package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Song Of Earendil", NewSongOfEarendil)
}

// NewSongOfEarendil creates a Song Of Earendil
// {3}{G}{U} - ENCHANTMENT
func NewSongOfEarendil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Song Of Earendil")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
