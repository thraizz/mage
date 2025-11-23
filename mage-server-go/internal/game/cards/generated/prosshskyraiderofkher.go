package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Prossh Skyraider Of Kher", NewProsshSkyraiderOfKher)
}

// NewProsshSkyraiderOfKher creates a Prossh Skyraider Of Kher
// {3}{B}{R}{G} - CREATURE
// Flying
func NewProsshSkyraiderOfKher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prossh Skyraider Of Kher")
	card.ManaCost = "{3}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("KherKeepKoboldToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, ManaSpentToCastCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
