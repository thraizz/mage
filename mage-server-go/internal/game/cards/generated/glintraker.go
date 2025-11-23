package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Glint Raker", NewGlintRaker)
}

// NewGlintRaker creates a Glint Raker
// {3}{U} - CREATURE
// Flying
func NewGlintRaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glint Raker")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAKE"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(GreatestAmongPermanentsValue.MANAVALUE_CONTROLLED_ARTIFACTS, StaticValue.get(0))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                         SavedDamageValue.MANY, 1,...)
	// card.AddAbility(ability2)
	return card, nil
}
