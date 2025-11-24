package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grand Master Of Flowers", NewGrandMasterOfFlowers)
}

// NewGrandMasterOfFlowers creates a Grand Master Of Flowers
// {2}{W}{W} - PLANESWALKER
// Flying, Indestructible
func NewGrandMasterOfFlowers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grand Master Of Flowers")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"BAHAMUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability1)
	return card, nil
}
