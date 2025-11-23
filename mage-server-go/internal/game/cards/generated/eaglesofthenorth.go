package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eagles Of The North", NewEaglesOfTheNorth)
}

// NewEaglesOfTheNorth creates a Eagles Of The North
// {5}{W} - CREATURE
// Flying
func NewEaglesOfTheNorth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eagles Of The North")
	card.ManaCost = "{5}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
