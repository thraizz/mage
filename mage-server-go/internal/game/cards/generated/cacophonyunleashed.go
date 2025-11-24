package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cacophony Unleashed", NewCacophonyUnleashed)
}

// NewCacophonyUnleashed creates a Cacophony Unleashed
// {5}{B}{B} - ENCHANTMENT
// Deathtouch
func NewCacophonyUnleashed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cacophony Unleashed")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"NIGHTMARE", "GOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: DestroyAllEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
