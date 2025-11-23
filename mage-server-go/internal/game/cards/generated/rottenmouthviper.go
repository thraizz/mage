package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rottenmouth Viper", NewRottenmouthViper)
}

// NewRottenmouthViper creates a Rottenmouth Viper
// {5}{B} - CREATURE
func NewRottenmouthViper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rottenmouth Viper")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "SNAKE"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
