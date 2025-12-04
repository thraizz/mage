package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sepulchral Primordial", NewSepulchralPrimordial)
}

// NewSepulchralPrimordial creates a Sepulchral Primordial
// {5}{B}{B} - CREATURE
func NewSepulchralPrimordial(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sepulchral Primordial")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: SepulchralPrimordialEffect()
	// card.AddAbility(ability0)
	return card, nil
}
