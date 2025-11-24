package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Davros Dalek Creator", NewDavrosDalekCreator)
}

// NewDavrosDalekCreator creates a Davros Dalek Creator
// {1}{U}{B}{R} - ARTIFACT CREATURE
func NewDavrosDalekCreator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Davros Dalek Creator")
	card.ManaCost = "{1}{U}{B}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ALIEN", "SCIENTIST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
