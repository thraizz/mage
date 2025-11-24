package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ureni The Song Unending", NewUreniTheSongUnending)
}

// NewUreniTheSongUnending creates a Ureni The Song Unending
// {5}{G}{U}{R} - CREATURE
// Flying
func NewUreniTheSongUnending(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ureni The Song Unending")
	card.ManaCost = "{5}{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}