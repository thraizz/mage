package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sharkey Tyrant Of The Shire", NewSharkeyTyrantOfTheShire)
}

// NewSharkeyTyrantOfTheShire creates a Sharkey Tyrant Of The Shire
// {2}{U}{B} - CREATURE
func NewSharkeyTyrantOfTheShire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sharkey Tyrant Of The Shire")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
