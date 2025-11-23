package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Liliana Of The Dark Realms", NewLilianaOfTheDarkRealms)
}

// NewLilianaOfTheDarkRealms creates a Liliana Of The Dark Realms
// {2}{B}{B} - PLANESWALKER
func NewLilianaOfTheDarkRealms(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Liliana Of The Dark Realms")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LILIANA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInHandEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
