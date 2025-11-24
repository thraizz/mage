package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Monastery Messenger", NewMonasteryMessenger)
}

// NewMonasteryMessenger creates a Monastery Messenger
// {2/U}{2/R}{2/W} - CREATURE
// Flying, Vigilance
func NewMonasteryMessenger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Monastery Messenger")
	card.ManaCost = "{2/U}{2/R}{2/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "SCOUT"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
