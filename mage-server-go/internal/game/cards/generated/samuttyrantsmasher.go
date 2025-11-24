package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Samut Tyrant Smasher", NewSamutTyrantSmasher)
}

// NewSamutTyrantSmasher creates a Samut Tyrant Smasher
// {2}{R/G}{R/G} - PLANESWALKER
func NewSamutTyrantSmasher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Samut Tyrant Smasher")
	card.ManaCost = "{2}{R/G}{R/G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SAMUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}