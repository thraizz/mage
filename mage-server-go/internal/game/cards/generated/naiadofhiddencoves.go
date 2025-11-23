package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Naiad Of Hidden Coves", NewNaiadOfHiddenCoves)
}

// NewNaiadOfHiddenCoves creates a Naiad Of Hidden Coves
// {2}{U} - ENCHANTMENT CREATURE
func NewNaiadOfHiddenCoves(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Naiad Of Hidden Coves")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"NYMPH"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
