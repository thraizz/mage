package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Trade Federation Battleship", NewTradeFederationBattleship)
}

// NewTradeFederationBattleship creates a Trade Federation Battleship
// {3}{W}{U}{B} - ARTIFACT CREATURE
func NewTradeFederationBattleship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trade Federation Battleship")
	card.ManaCost = "{3}{W}{U}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"STARSHIP"}
	card.Power = "0"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1, filter, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
