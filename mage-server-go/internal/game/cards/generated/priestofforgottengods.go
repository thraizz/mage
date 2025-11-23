package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Priest Of Forgotten Gods", NewPriestOfForgottenGods)
}

// NewPriestOfForgottenGods creates a Priest Of Forgotten Gods
// {1}{B} - CREATURE
func NewPriestOfForgottenGods(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Priest Of Forgotten Gods")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewLoseLifeEffect(2)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
