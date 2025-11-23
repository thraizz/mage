package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Marit Lages Slumber", NewMaritLagesSlumber)
}

// NewMaritLagesSlumber creates a Marit Lages Slumber
// {1}{U} - ENCHANTMENT
func NewMaritLagesSlumber(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Marit Lages Slumber")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY", "SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new MaritLa...)
	// card.AddAbility(ability1)
	return card, nil
}
