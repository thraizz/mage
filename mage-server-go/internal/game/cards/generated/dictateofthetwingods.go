package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dictate Of The Twin Gods", NewDictateOfTheTwinGods)
}

// NewDictateOfTheTwinGods creates a Dictate Of The Twin Gods
// {3}{R}{R} - ENCHANTMENT
// Flash
func NewDictateOfTheTwinGods(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dictate Of The Twin Gods")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	return card, nil
}
