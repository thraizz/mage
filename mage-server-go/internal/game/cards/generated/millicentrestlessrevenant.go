package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Millicent Restless Revenant", NewMillicentRestlessRevenant)
}

// NewMillicentRestlessRevenant creates a Millicent Restless Revenant
// {5}{W}{U} - CREATURE
// Flying
func NewMillicentRestlessRevenant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Millicent Restless Revenant")
	card.ManaCost = "{5}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
