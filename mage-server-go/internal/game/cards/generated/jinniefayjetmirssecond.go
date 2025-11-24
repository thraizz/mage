package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jinnie Fay Jetmirs Second", NewJinnieFayJetmirsSecond)
}

// NewJinnieFayJetmirsSecond creates a Jinnie Fay Jetmirs Second
// {R/G}{G}{G/W} - CREATURE
func NewJinnieFayJetmirsSecond(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jinnie Fay Jetmirs Second")
	card.ManaCost = "{R/G}{G}{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}