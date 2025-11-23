package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abattoir Ghoul", NewAbattoirGhoul)
}

// NewAbattoirGhoul creates a Abattoir Ghoul
// {3}{B} - CREATURE
// FirstStrike
func NewAbattoirGhoul(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abattoir Ghoul")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	return card, nil
}
