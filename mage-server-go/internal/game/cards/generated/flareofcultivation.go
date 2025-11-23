package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flare Of Cultivation", NewFlareOfCultivation)
}

// NewFlareOfCultivation creates a Flare Of Cultivation
// {1}{G}{G} - SORCERY
func NewFlareOfCultivation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flare Of Cultivation")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
