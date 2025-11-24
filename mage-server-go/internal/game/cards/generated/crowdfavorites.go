package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crowd Favorites", NewCrowdFavorites)
}

// NewCrowdFavorites creates a Crowd Favorites
// {6}{W} - CREATURE
func NewCrowdFavorites(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crowd Favorites")
	card.ManaCost = "{6}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: TapTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
