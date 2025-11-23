package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Ajani Adversary Of Tyrants", NewAjaniAdversaryOfTyrants)
}

// NewAjaniAdversaryOfTyrants creates a Ajani Adversary Of Tyrants
// {2}{W}{W} - PLANESWALKER
func NewAjaniAdversaryOfTyrants(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ajani Adversary Of Tyrants")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"AJANI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
