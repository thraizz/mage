package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mirage Mirror", NewMirageMirror)
}

// NewMirageMirror creates a Mirage Mirror
// {3} - ARTIFACT
func NewMirageMirror(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mirage Mirror")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - MirageMirrorCopyEffect()
	// card.AddAbility(ability0)
	return card, nil
}
