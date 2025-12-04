package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Afterlife From The Loam", NewAfterlifeFromTheLoam)
}

// NewAfterlifeFromTheLoam creates a Afterlife From The Loam
// {5}{B}{B}{B} - SORCERY
func NewAfterlifeFromTheLoam(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Afterlife From The Loam")
	card.ManaCost = "{5}{B}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
