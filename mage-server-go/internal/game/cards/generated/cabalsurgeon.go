package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cabal Surgeon", NewCabalSurgeon)
}

// NewCabalSurgeon creates a Cabal Surgeon
// {2}{B}{B} - CREATURE
func NewCabalSurgeon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cabal Surgeon")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "MINION"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
