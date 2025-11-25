package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wowzer The Aspirational", NewWowzerTheAspirational)
}

// NewWowzerTheAspirational creates a Wowzer The Aspirational
// {C}{W}{U}{B}{R}{G}{S} - CREATURE
func NewWowzerTheAspirational(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wowzer The Aspirational")
	card.ManaCost = "{C}{W}{U}{B}{R}{G}{S}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Supertypes = []string{"LEGENDARY", "SNOW"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
