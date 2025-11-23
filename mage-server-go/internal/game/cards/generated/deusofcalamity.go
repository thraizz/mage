package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deus Of Calamity", NewDeusOfCalamity)
}

// NewDeusOfCalamity creates a Deus Of Calamity
// {R/G}{R/G}{R/G}{R/G}{R/G} - CREATURE
// Trample
func NewDeusOfCalamity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deus Of Calamity")
	card.ManaCost = "{R/G}{R/G}{R/G}{R/G}{R/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "AVATAR"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
