package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Terra Stomper", NewTerraStomper)
}

// NewTerraStomper creates a Terra Stomper
// {3}{G}{G}{G} - CREATURE
// Trample
func NewTerraStomper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Terra Stomper")
	card.ManaCost = "{3}{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
