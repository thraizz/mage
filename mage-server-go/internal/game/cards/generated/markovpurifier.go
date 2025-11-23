package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Markov Purifier", NewMarkovPurifier)
}

// NewMarkovPurifier creates a Markov Purifier
// {1}{W}{B} - CREATURE
// Lifelink
func NewMarkovPurifier(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Markov Purifier")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "CLERIC"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(1), new Generic...)
	// card.AddAbility(ability1)
	return card, nil
}
