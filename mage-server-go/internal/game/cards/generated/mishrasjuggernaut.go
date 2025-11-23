package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mishras Juggernaut", NewMishrasJuggernaut)
}

// NewMishrasJuggernaut creates a Mishras Juggernaut
// {5} - ARTIFACT CREATURE
// Trample
func NewMishrasJuggernaut(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mishras Juggernaut")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"JUGGERNAUT"}
	card.Power = "5"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
