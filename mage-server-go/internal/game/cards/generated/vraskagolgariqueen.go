package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vraska Golgari Queen", NewVraskaGolgariQueen)
}

// NewVraskaGolgariQueen creates a Vraska Golgari Queen
// {2}{B}{G} - PLANESWALKER
func NewVraskaGolgariQueen(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vraska Golgari Queen")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VRASKA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}