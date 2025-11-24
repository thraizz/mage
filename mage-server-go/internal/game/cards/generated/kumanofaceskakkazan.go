package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kumano Faces Kakkazan", NewKumanoFacesKakkazan)
}

// NewKumanoFacesKakkazan creates a Kumano Faces Kakkazan
// {R} - ENCHANTMENT
func NewKumanoFacesKakkazan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kumano Faces Kakkazan")
	card.ManaCost = "{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}