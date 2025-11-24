package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maelstrom Muse", NewMaelstromMuse)
}

// NewMaelstromMuse creates a Maelstrom Muse
// {1}{U}{U/R}{R} - CREATURE
// Flying
func NewMaelstromMuse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maelstrom Muse")
	card.ManaCost = "{1}{U}{U/R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DJINN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}