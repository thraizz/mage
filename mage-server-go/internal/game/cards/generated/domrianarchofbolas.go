package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Domri Anarch Of Bolas", NewDomriAnarchOfBolas)
}

// NewDomriAnarchOfBolas creates a Domri Anarch Of Bolas
// {1}{R}{G} - PLANESWALKER
func NewDomriAnarchOfBolas(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Domri Anarch Of Bolas")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOMRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
