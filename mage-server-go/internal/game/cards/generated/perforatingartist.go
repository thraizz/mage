package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Perforating Artist", NewPerforatingArtist)
}

// NewPerforatingArtist creates a Perforating Artist
// {1}{B}{R} - CREATURE
// Deathtouch
func NewPerforatingArtist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Perforating Artist")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEVIL"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}
