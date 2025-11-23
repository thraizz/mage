package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kuro Pitlord", NewKuroPitlord)
}

// NewKuroPitlord creates a Kuro Pitlord
// {6}{B}{B}{B} - CREATURE
func NewKuroPitlord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kuro Pitlord")
	card.ManaCost = "{6}{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(-1, -1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
