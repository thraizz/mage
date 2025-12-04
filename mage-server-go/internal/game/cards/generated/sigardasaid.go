package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sigardas Aid", NewSigardasAid)
}

// NewSigardasAid creates a Sigardas Aid
// {W} - ENCHANTMENT
func NewSigardasAid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sigardas Aid")
	card.ManaCost = "{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldControlledTriggeredAbility
	//   - Effect: SigardasAidEffect()
	// card.AddAbility(ability0)
	return card, nil
}
