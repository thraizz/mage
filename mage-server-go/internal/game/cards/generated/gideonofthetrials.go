package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gideon Of The Trials", NewGideonOfTheTrials)
}

// NewGideonOfTheTrials creates a Gideon Of The Trials
// {1}{W}{W} - PLANESWALKER
// Indestructible
func NewGideonOfTheTrials(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gideon Of The Trials")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GIDEON", "HUMAN", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	return card, nil
}
