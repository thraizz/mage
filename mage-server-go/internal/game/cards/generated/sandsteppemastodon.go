package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sandsteppe Mastodon", NewSandsteppeMastodon)
}

// NewSandsteppeMastodon creates a Sandsteppe Mastodon
// {5}{G}{G} - CREATURE
// Reach
func NewSandsteppeMastodon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sandsteppe Mastodon")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEPHANT"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
