package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lifecraft Awakening", NewLifecraftAwakening)
}

// NewLifecraftAwakening creates a Lifecraft Awakening
// {X}{G} - INSTANT
func NewLifecraftAwakening(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lifecraft Awakening")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}