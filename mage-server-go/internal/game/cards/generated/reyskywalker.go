package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rey Skywalker", NewReySkywalker)
}

// NewReySkywalker creates a Rey Skywalker
// {2}{U}{U} - PLANESWALKER
func NewReySkywalker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rey Skywalker")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewUntapEffect()).
		AddEffect(abilities.NewTapEffect()).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
