package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Tocasia Dig Site Mentor", NewTocasiaDigSiteMentor)
}

// NewTocasiaDigSiteMentor creates a Tocasia Dig Site Mentor
// {1}{G}{W}{U} - CREATURE
func NewTocasiaDigSiteMentor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tocasia Dig Site Mentor")
	card.ManaCost = "{1}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSurveilEffect(1)).
		AddEffect(abilities.NewGrantAbilityEffect(new SimpleActivatedAbility( new SurveilEffect(1, false), new TapSourceCost() ))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}