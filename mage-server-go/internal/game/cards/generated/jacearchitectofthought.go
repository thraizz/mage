package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jace Architect Of Thought", NewJaceArchitectOfThought)
}

// NewJaceArchitectOfThought creates a Jace Architect Of Thought
// {2}{U}{U} - PLANESWALKER
func NewJaceArchitectOfThought(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jace Architect Of Thought")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JACE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-1, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
