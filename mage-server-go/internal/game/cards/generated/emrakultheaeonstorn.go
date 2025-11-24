package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Emrakul The Aeons Torn", NewEmrakulTheAeonsTorn)
}

// NewEmrakulTheAeonsTorn creates a Emrakul The Aeons Torn
// {15} - CREATURE
// Flying
func NewEmrakulTheAeonsTorn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Emrakul The Aeons Torn")
	card.ManaCost = "{15}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "15"
	card.Toughness = "15"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}