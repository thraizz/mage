package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ethereal Champion", NewEtherealChampion)
}

// NewEtherealChampion creates a Ethereal Champion
// {2}{W}{W}{W} - CREATURE
func NewEtherealChampion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ethereal Champion")
	card.ManaCost = "{2}{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
