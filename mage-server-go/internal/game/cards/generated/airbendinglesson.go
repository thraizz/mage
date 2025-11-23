package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Airbending Lesson", NewAirbendingLesson)
}

// NewAirbendingLesson creates a Airbending Lesson
// {2}{W} - INSTANT
func NewAirbendingLesson(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Airbending Lesson")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
