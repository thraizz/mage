package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Attendant Of Vraska", NewAttendantOfVraska)
}

// NewAttendantOfVraska creates a Attendant Of Vraska
// {1}{B}{G} - CREATURE
func NewAttendantOfVraska(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Attendant Of Vraska")
	card.ManaCost = "{1}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(SourcePermanentPowerValue.NOT_NEGATIVE)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
