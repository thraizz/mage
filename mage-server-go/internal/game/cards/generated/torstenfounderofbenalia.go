package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Torsten Founder Of Benalia", NewTorstenFounderOfBenalia)
}

// NewTorstenFounderOfBenalia creates a Torsten Founder Of Benalia
// {5}{G}{W} - CREATURE
func NewTorstenFounderOfBenalia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Torsten Founder Of Benalia")
	card.ManaCost = "{5}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                 7, Integer.MAX_VALUE, filter, Put...)
	// card.AddAbility(ability0)
	token1_0, err := token.GetToken("SoldierToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffectAmount(token1_0, 7)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
