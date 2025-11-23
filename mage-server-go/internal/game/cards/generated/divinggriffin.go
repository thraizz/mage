package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Diving Griffin", NewDivingGriffin)
}

// NewDivingGriffin creates a Diving Griffin
// {1}{W}{W} - CREATURE
// Flying, Vigilance
func NewDivingGriffin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Diving Griffin")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GRIFFIN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
