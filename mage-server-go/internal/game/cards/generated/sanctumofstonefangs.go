package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanctum Of Stone Fangs", NewSanctumOfStoneFangs)
}

// NewSanctumOfStoneFangs creates a Sanctum Of Stone Fangs
// {1}{B} - ENCHANTMENT
func NewSanctumOfStoneFangs(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanctum Of Stone Fangs")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SHRINE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
