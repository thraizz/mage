package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dawnhart Rejuvenator", NewDawnhartRejuvenator)
}

// NewDawnhartRejuvenator creates a Dawnhart Rejuvenator
// {3}{G} - CREATURE
func NewDawnhartRejuvenator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dawnhart Rejuvenator")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARLOCK"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
