package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Daretti Scrap Savant", NewDarettiScrapSavant)
}

// NewDarettiScrapSavant creates a Daretti Scrap Savant
// {3}{R} - PLANESWALKER
func NewDarettiScrapSavant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daretti Scrap Savant")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DARETTI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: DarettiSacrificeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
