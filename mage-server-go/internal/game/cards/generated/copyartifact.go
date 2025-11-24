package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Copy Artifact", NewCopyArtifact)
}

// NewCopyArtifact creates a Copy Artifact
// {1}{U} - ENCHANTMENT
func NewCopyArtifact(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Copy Artifact")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(new FilterArtifactPermanent(), new CardTypeCopyApp...)
	// card.AddAbility(ability0)
	return card, nil
}
