package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Green Goblin Revenant", NewGreenGoblinRevenant)
}

// NewGreenGoblinRevenant creates a Green Goblin Revenant
// {3}{B}{R} - CREATURE
// Flying, Deathtouch
func NewGreenGoblinRevenant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Green Goblin Revenant")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "HUMAN", "VILLAIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	// card.AddAbility(ability2)
	return card, nil
}
