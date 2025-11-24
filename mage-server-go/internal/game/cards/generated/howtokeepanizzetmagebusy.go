package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("How To Keep An Izzet Mage Busy", NewHowToKeepAnIzzetMageBusy)
}

// NewHowToKeepAnIzzetMageBusy creates a How To Keep An Izzet Mage Busy
// {U/R} - SORCERY
func NewHowToKeepAnIzzetMageBusy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "How To Keep An Izzet Mage Busy")
	card.ManaCost = "{U/R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnToHandSourceEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
