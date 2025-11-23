package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scene Of The Crime", NewSceneOfTheCrime)
}

// NewSceneOfTheCrime creates a Scene Of The Crime
//   - ARTIFACT LAND
func NewSceneOfTheCrime(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scene Of The Crime")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "LAND"}
	card.Subtypes = []string{"CLUE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	return card, nil
}
