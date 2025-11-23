package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Guidelight Synergist", NewGuidelightSynergist)
}

// NewGuidelightSynergist creates a Guidelight Synergist
// {3}{W} - ARTIFACT CREATURE
// Flying
func NewGuidelightSynergist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Guidelight Synergist")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT", "ARTIFICER"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(ArtifactYouControlCount.instance, StaticValue.get(0))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
