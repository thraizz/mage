package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Fire Nation Archers", NewFireNationArchers)
}

// NewFireNationArchers creates a Fire Nation Archers
// {3}{R} - CREATURE
// Reach
func NewFireNationArchers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fire Nation Archers")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARCHER"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("SoldierRedToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
