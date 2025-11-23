package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Mmmenon The Right Hand", NewMmmenonTheRightHand)
}

// NewMmmenonTheRightHand creates a Mmmenon The Right Hand
// {3}{U}{U} - CREATURE
// Flying
func NewMmmenonTheRightHand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mmmenon The Right Hand")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"JELLYFISH", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new ConditionalColoredManaAbility(Mana.BlueMana(1), new MmmenonTheRightHandManaBuilder()))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}