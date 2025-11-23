package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Marrow Chomper", NewMarrowChomper)
}

// NewMarrowChomper creates a Marrow Chomper
// {3}{B}{G} - CREATURE
func NewMarrowChomper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Marrow Chomper")
	card.ManaCost = "{3}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "LIZARD"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(new DevouredCreaturesCount(2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}