package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Wall Of Kelp", NewWallOfKelp)
}

// NewWallOfKelp creates a Wall Of Kelp
// {U}{U} - CREATURE
// Defender
func NewWallOfKelp(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wall Of Kelp")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PLANT", "WALL"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("KelpToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
