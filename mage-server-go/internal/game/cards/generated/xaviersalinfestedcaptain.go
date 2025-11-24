package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Xavier Sal Infested Captain", NewXavierSalInfestedCaptain)
}

// NewXavierSalInfestedCaptain creates a Xavier Sal Infested Captain
// {B}{G}{U} - CREATURE
func NewXavierSalInfestedCaptain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Xavier Sal Infested Captain")
	card.ManaCost = "{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "FUNGUS", "PIRATE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}