package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ordruun Mentor", NewOrdruunMentor)
}

// NewOrdruunMentor creates a Ordruun Mentor
// {2}{R/W} - CREATURE
func NewOrdruunMentor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ordruun Mentor")
	card.ManaCost = "{2}{R/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MINOTAUR", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
