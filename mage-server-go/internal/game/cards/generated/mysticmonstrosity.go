package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Mystic Monstrosity", NewMysticMonstrosity)
}

// NewMysticMonstrosity creates a Mystic Monstrosity
//   - ARTIFACT CREATURE
func NewMysticMonstrosity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mystic Monstrosity")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "5"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("AnyColorManaAbility", effects.DurationWhileOnBattlefield)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
