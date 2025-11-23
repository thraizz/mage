package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elderscale Wurm", NewElderscaleWurm)
}

// NewElderscaleWurm creates a Elderscale Wurm
// {4}{G}{G}{G} - CREATURE
// Trample
func NewElderscaleWurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elderscale Wurm")
	card.ManaCost = "{4}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
