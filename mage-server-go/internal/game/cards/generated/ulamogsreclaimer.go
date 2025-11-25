package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ulamogs Reclaimer", NewUlamogsReclaimer)
}

// NewUlamogsReclaimer creates a Ulamogs Reclaimer
// {4}{U} - CREATURE
func NewUlamogsReclaimer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ulamogs Reclaimer")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "PROCESSOR"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: DoIfCostPaid(new ReturnFromGraveyardToHandTargetEffect(), new E...)
	// card.AddAbility(ability0)
	return card, nil
}
