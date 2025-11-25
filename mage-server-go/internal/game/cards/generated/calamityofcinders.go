package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Calamity Of Cinders", NewCalamityOfCinders)
}

// NewCalamityOfCinders creates a Calamity Of Cinders
// {5}{R}{R} - SORCERY
func NewCalamityOfCinders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Calamity Of Cinders")
	card.ManaCost = "{5}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(6, filter)
	// card.AddAbility(ability0)
	return card, nil
}
