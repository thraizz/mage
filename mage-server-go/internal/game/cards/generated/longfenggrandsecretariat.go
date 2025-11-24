package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Long Feng Grand Secretariat", NewLongFengGrandSecretariat)
}

// NewLongFengGrandSecretariat creates a Long Feng Grand Secretariat
// {1}{B/G}{B/G} - CREATURE
func NewLongFengGrandSecretariat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Long Feng Grand Secretariat")
	card.ManaCost = "{1}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}