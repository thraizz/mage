package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gleancrawler", NewGleancrawler)
}

// NewGleancrawler creates a Gleancrawler
// {3}{B/G}{B/G}{B/G} - CREATURE
// Trample
func NewGleancrawler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gleancrawler")
	card.ManaCost = "{3}{B/G}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "HORROR"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
