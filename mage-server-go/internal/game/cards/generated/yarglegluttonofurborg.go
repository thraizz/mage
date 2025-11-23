package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yargle Glutton Of Urborg", NewYargleGluttonOfUrborg)
}

// NewYargleGluttonOfUrborg creates a Yargle Glutton Of Urborg
// {4}{B} - CREATURE
func NewYargleGluttonOfUrborg(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yargle Glutton Of Urborg")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "9"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
