package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leyline Of The Guildpact", NewLeylineOfTheGuildpact)
}

// NewLeylineOfTheGuildpact creates a Leyline Of The Guildpact
// {G/W}{G/U}{B/G}{R/G} - ENCHANTMENT
func NewLeylineOfTheGuildpact(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leyline Of The Guildpact")
	card.ManaCost = "{G/W}{G/U}{B/G}{R/G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
