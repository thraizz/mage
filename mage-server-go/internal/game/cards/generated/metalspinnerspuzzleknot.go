package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Metalspinners Puzzleknot", NewMetalspinnersPuzzleknot)
}

// NewMetalspinnersPuzzleknot creates a Metalspinners Puzzleknot
// {2} - ARTIFACT
func NewMetalspinnersPuzzleknot(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Metalspinners Puzzleknot")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewLoseLifeEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
