package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urza Academy Headmaster", NewUrzaAcademyHeadmaster)
}

// NewUrzaAcademyHeadmaster creates a Urza Academy Headmaster
// {W}{U}{B}{R}{G} - PLANESWALKER
func NewUrzaAcademyHeadmaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urza Academy Headmaster")
	card.ManaCost = "{W}{U}{B}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"URZA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
