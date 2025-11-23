package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kethis The Hidden Hand", NewKethisTheHiddenHand)
}

// NewKethisTheHiddenHand creates a Kethis The Hidden Hand
// {W}{B}{G} - CREATURE
func NewKethisTheHiddenHand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kethis The Hidden Hand")
	card.ManaCost = "{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
