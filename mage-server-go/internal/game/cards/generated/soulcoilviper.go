package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Soulcoil Viper", NewSoulcoilViper)
}

// NewSoulcoilViper creates a Soulcoil Viper
// {2}{B} - CREATURE
func NewSoulcoilViper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soulcoil Viper")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldWithCounterTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
