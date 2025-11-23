package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tidehollow Sculler", NewTidehollowSculler)
}

// NewTidehollowSculler creates a Tidehollow Sculler
// {W}{B} - ARTIFACT CREATURE
func NewTidehollowSculler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tidehollow Sculler")
	card.ManaCost = "{W}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ZOMBIE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
