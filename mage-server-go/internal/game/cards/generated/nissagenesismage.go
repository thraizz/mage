package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nissa Genesis Mage", NewNissaGenesisMage)
}

// NewNissaGenesisMage creates a Nissa Genesis Mage
// {5}{G}{G} - PLANESWALKER
func NewNissaGenesisMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nissa Genesis Mage")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NISSA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewUntapEffect("untap up to two target creatures and up to two target lands")).
		AddEffect(abilities.NewBoostEffect(5, 5)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(10, Integer.MAX_VALUE, filter, PutCards.BATTLEFIEL...)
	// card.AddAbility(ability1)
	return card, nil
}