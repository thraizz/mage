package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Academy Journeymage", NewAcademyJourneymage)
}

// NewAcademyJourneymage creates a Academy Journeymage
// {4}{U} - CREATURE
func NewAcademyJourneymage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Academy Journeymage")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
