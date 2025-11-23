package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leonardo Worldly Warrior", NewLeonardoWorldlyWarrior)
}

// NewLeonardoWorldlyWarrior creates a Leonardo Worldly Warrior
// {7}{W} - CREATURE
// DoubleStrike
func NewLeonardoWorldlyWarrior(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leonardo Worldly Warrior")
	card.ManaCost = "{7}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MUTANT", "NINJA", "TURTLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability0)
	return card, nil
}
