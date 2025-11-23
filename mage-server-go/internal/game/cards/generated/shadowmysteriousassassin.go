package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shadow Mysterious Assassin", NewShadowMysteriousAssassin)
}

// NewShadowMysteriousAssassin creates a Shadow Mysterious Assassin
// {2}{B} - CREATURE
// Deathtouch
func NewShadowMysteriousAssassin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shadow Mysterious Assassin")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ASSASSIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}
