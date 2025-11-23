package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Garruk The Veil Cursed", NewGarrukTheVeilCursed)
}

// NewGarrukTheVeilCursed creates a Garruk The Veil Cursed
//  - PLANESWALKER
func NewGarrukTheVeilCursed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Garruk The Veil Cursed")
	card.ManaCost = ""
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GARRUK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(new WolfTokenWithDeathtouch())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new SearchLibraryPutInHandEffect(...)
	// card.AddAbility(ability1)
	return card, nil
}