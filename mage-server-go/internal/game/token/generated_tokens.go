package token

// This file contains predefined token types.
// Auto-generated from Java token implementations.

// ========================================
// NewATATToken
// ========================================

// NewATATToken creates a token.
func NewATATToken() *Token {
	tok := NewToken("ATATToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ATAT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewAetherbornToken
// ========================================

// NewAetherbornToken creates a X/X black Aetherborn creature token, where X is the amount of {E} paid this way.
func NewAetherbornToken() *Token {
	tok := NewToken("AetherbornToken", "X/X black Aetherborn creature token, where X is the amount of {E} paid this way")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("AETHERBORN")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewAjanisPridemateToken
// ========================================

// NewAjanisPridemateToken creates a token.
func NewAjanisPridemateToken() *Token {
	tok := NewToken("AjanisPridemateToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewAkroanSoldierToken
// ========================================

// NewAkroanSoldierToken creates a 1/1 red Soldier creature token with haste.
func NewAkroanSoldierToken() *Token {
	tok := NewToken("AkroanSoldierToken", "1/1 red Soldier creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewAlien00Token
// ========================================

// NewAlien00Token creates a 0/0 blue Alien creature token.
func NewAlien00Token() *Token {
	tok := NewToken("Alien00Token", "0/0 blue Alien creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewAlienAngelToken
// ========================================

// NewAlienAngelToken creates a token.
func NewAlienAngelToken() *Token {
	tok := NewToken("AlienAngelToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.AddSubtype("ANGEL")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	tok.AddAbility("first strike")
	return tok
}

// ========================================
// NewAlienInsectToken
// ========================================

// NewAlienInsectToken creates a 1/1 green and white Alien Insect creature token with flying.
func NewAlienInsectToken() *Token {
	tok := NewToken("AlienInsectToken", "1/1 green and white Alien Insect creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewAlienRhinoToken
// ========================================

// NewAlienRhinoToken creates a 4/4 white Alien Rhino creature token.
func NewAlienRhinoToken() *Token {
	tok := NewToken("AlienRhinoToken", "4/4 white Alien Rhino creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.AddSubtype("RHINO")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewAlienSalamanderToken
// ========================================

// NewAlienSalamanderToken creates a 2/2 green Alien Salamander creature token with islandwalk.
func NewAlienSalamanderToken() *Token {
	tok := NewToken("AlienSalamanderToken", "2/2 green Alien Salamander creature token with islandwalk")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.AddSubtype("SALAMANDER")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewAlienToken
// ========================================

// NewAlienToken creates a 2/2 white Alien creature token.
func NewAlienToken() *Token {
	tok := NewToken("AlienToken", "2/2 white Alien creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewAlienWarriorToken
// ========================================

// NewAlienWarriorToken creates a 2/2 red Alien Warrior creature token.
func NewAlienWarriorToken() *Token {
	tok := NewToken("AlienWarriorToken", "2/2 red Alien Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALIEN")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewAllyToken
// ========================================

// NewAllyToken creates a 1/1 white Ally creature token.
func NewAllyToken() *Token {
	tok := NewToken("AllyToken", "1/1 white Ally creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ALLY")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewAngel33Token
// ========================================

// NewAngel33Token creates a 3/3 white Angel creature token with flying.
func NewAngel33Token() *Token {
	tok := NewToken("Angel33Token", "3/3 white Angel creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewAngelToken
// ========================================

// NewAngelToken creates a 4/4 white Angel creature token with flying.
func NewAngelToken() *Token {
	tok := NewToken("AngelToken", "4/4 white Angel creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewAngelVigilanceToken
// ========================================

// NewAngelVigilanceToken creates a 4/4 white Angel creature token with flying and vigilance.
func NewAngelVigilanceToken() *Token {
	tok := NewToken("AngelVigilanceToken", "4/4 white Angel creature token with flying and vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewAngelWarriorToken
// ========================================

// NewAngelWarriorToken creates a 4/4 white Angel Warrior creature token with flying.
func NewAngelWarriorToken() *Token {
	tok := NewToken("AngelWarriorToken", "4/4 white Angel Warrior creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewAngelWarriorVigilanceToken
// ========================================

// NewAngelWarriorVigilanceToken creates a 4/4 white Angel Warrior creature token with flying and vigilance.
func NewAngelWarriorVigilanceToken() *Token {
	tok := NewToken("AngelWarriorVigilanceToken", "4/4 white Angel Warrior creature token with flying and vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewAngeloToken
// ========================================

// NewAngeloToken creates a Angelo, a legendary 1/1 green and white Dog creature token.
func NewAngeloToken() *Token {
	tok := NewToken("AngeloToken", "Angelo, a legendary 1/1 green and white Dog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewAnotherSpiritToken
// ========================================

// NewAnotherSpiritToken creates a 3/3 white Spirit creature token with flying.
func NewAnotherSpiritToken() *Token {
	tok := NewToken("AnotherSpiritToken", "3/3 white Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewApeToken
// ========================================

// NewApeToken creates a 3/3 green Ape creature token.
func NewApeToken() *Token {
	tok := NewToken("ApeToken", "3/3 green Ape creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("APE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewArchitectOfTheUntamedBeastToken
// ========================================

// NewArchitectOfTheUntamedBeastToken creates a 6/6 colorless Beast artifact creature token.
func NewArchitectOfTheUntamedBeastToken() *Token {
	tok := NewToken("ArchitectOfTheUntamedBeastToken", "6/6 colorless Beast artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetPowerToughness(6, 6)
	return tok
}

// ========================================
// NewArtifactWallToken
// ========================================

// NewArtifactWallToken creates a 0/4 colorless Wall artifact creature token with defender.
func NewArtifactWallToken() *Token {
	tok := NewToken("ArtifactWallToken", "0/4 colorless Wall artifact creature token with defender")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetPowerToughness(0, 4)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewAshiokNightmareMuseToken
// ========================================

// NewAshiokNightmareMuseToken creates a token.
func NewAshiokNightmareMuseToken() *Token {
	tok := NewToken("AshiokNightmareMuseToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("NIGHTMARE")
	tok.SetColor(Color{Blue: true, Black: true})
	tok.SetPowerToughness(2, 3)
	return tok
}

// ========================================
// NewAshiokWickedManipulatorNightmareToken
// ========================================

// NewAshiokWickedManipulatorNightmareToken creates a token.
func NewAshiokWickedManipulatorNightmareToken() *Token {
	tok := NewToken("AshiokWickedManipulatorNightmareToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("NIGHTMARE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewAshnodZombieToken
// ========================================

// NewAshnodZombieToken creates a 3/3 colorless Zombie artifact creature token.
func NewAshnodZombieToken() *Token {
	tok := NewToken("AshnodZombieToken", "3/3 colorless Zombie artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewAssassinMenaceToken
// ========================================

// NewAssassinMenaceToken creates a 1/1 black Assassin creature token with menace.
func NewAssassinMenaceToken() *Token {
	tok := NewToken("AssassinMenaceToken", "1/1 black Assassin creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ASSASSIN")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewAssassinToken
// ========================================

// NewAssassinToken creates a token.
func NewAssassinToken() *Token {
	tok := NewToken("AssassinToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ASSASSIN")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewAssemblyWorkerToken
// ========================================

// NewAssemblyWorkerToken creates a 2/2 colorless Assembly-Worker artifact creature token.
func NewAssemblyWorkerToken() *Token {
	tok := NewToken("AssemblyWorkerToken", "2/2 colorless Assembly-Worker artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ASSEMBLY_WORKER")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewAtlaPalaniToken
// ========================================

// NewAtlaPalaniToken creates a 0/1 green Egg creature token with defender.
func NewAtlaPalaniToken() *Token {
	tok := NewToken("AtlaPalaniToken", "0/1 green Egg creature token with defender")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("EGG")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewAvacynToken
// ========================================

// NewAvacynToken creates a Avacyn, a legendary 8/8 white Angel creature token with flying, vigilance, and indestructible.
func NewAvacynToken() *Token {
	tok := NewToken("AvacynToken", "Avacyn, a legendary 8/8 white Angel creature token with flying, vigilance, and indestructible")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(8, 8)
	tok.AddAbility("flying")
	tok.AddAbility("vigilance")
	tok.AddAbility("indestructible")
	return tok
}

// ========================================
// NewAvatarToken
// ========================================

// NewAvatarToken creates a token.
func NewAvatarToken() *Token {
	tok := NewToken("AvatarToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("AVATAR")
	tok.SetColor(Color{White: true})
	return tok
}

// ========================================
// NewBadgerToken
// ========================================

// NewBadgerToken creates a 3/3 green Badger creature token.
func NewBadgerToken() *Token {
	tok := NewToken("BadgerToken", "3/3 green Badger creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BADGER")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewBananaToken
// ========================================

// NewBananaToken creates a token.
func NewBananaToken() *Token {
	tok := NewToken("BananaToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewBaruFistOfKrosaToken
// ========================================

// NewBaruFistOfKrosaToken creates a X/X green Wurm creature token, where X is the number of lands you control.
func NewBaruFistOfKrosaToken() *Token {
	tok := NewToken("BaruFistOfKrosaToken", "X/X green Wurm creature token, where X is the number of lands you control")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewBat21Token
// ========================================

// NewBat21Token creates a 2/1 black Bat creature token with flying.
func NewBat21Token() *Token {
	tok := NewToken("Bat21Token", "2/1 black Bat creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBatToken
// ========================================

// NewBatToken creates a 1/1 black Bat creature token with flying.
func NewBatToken() *Token {
	tok := NewToken("BatToken", "1/1 black Bat creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBearToken
// ========================================

// NewBearToken creates a 2/2 green Bear creature token.
func NewBearToken() *Token {
	tok := NewToken("BearToken", "2/2 green Bear creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewBearsCompanionBearToken
// ========================================

// NewBearsCompanionBearToken creates a 4/4 green Bear creature token.
func NewBearsCompanionBearToken() *Token {
	tok := NewToken("BearsCompanionBearToken", "4/4 green Bear creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewBeastToken
// ========================================

// NewBeastToken creates a 3/3 green Beast creature token.
func NewBeastToken() *Token {
	tok := NewToken("BeastToken", "3/3 green Beast creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewBeastieToken
// ========================================

// NewBeastieToken creates a token.
func NewBeastieToken() *Token {
	tok := NewToken("BeastieToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewBeauToken
// ========================================

// NewBeauToken creates a token.
func NewBeauToken() *Token {
	tok := NewToken("BeauToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OX")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewBelzenlokClericToken
// ========================================

// NewBelzenlokClericToken creates a 0/1 black Cleric creature token.
func NewBelzenlokClericToken() *Token {
	tok := NewToken("BelzenlokClericToken", "0/1 black Cleric creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CLERIC")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewBelzenlokDemonToken
// ========================================

// NewBelzenlokDemonToken creates a token.
func NewBelzenlokDemonToken() *Token {
	tok := NewToken("BelzenlokDemonToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(6, 6)
	tok.AddAbility("flying")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewBiogenicOozeToken
// ========================================

// NewBiogenicOozeToken creates a 2/2 green Ooze creature token.
func NewBiogenicOozeToken() *Token {
	tok := NewToken("BiogenicOozeToken", "2/2 green Ooze creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewBirdIllusionToken
// ========================================

// NewBirdIllusionToken creates a 1/1 blue Bird Illusion creature token with flying.
func NewBirdIllusionToken() *Token {
	tok := NewToken("BirdIllusionToken", "1/1 blue Bird Illusion creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBirdSoldierToken
// ========================================

// NewBirdSoldierToken creates a 1/1 white Bird Soldier creature token with flying.
func NewBirdSoldierToken() *Token {
	tok := NewToken("BirdSoldierToken", "1/1 white Bird Soldier creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBirdToken
// ========================================

// NewBirdToken creates a 1/1 white Bird creature token with flying.
func NewBirdToken() *Token {
	tok := NewToken("BirdToken", "1/1 white Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBirdVigilanceToken
// ========================================

// NewBirdVigilanceToken creates a 1/1 blue Bird creature token with flying and vigilance.
func NewBirdVigilanceToken() *Token {
	tok := NewToken("BirdVigilanceToken", "1/1 blue Bird creature token with flying and vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewBlack22BirdToken
// ========================================

// NewBlack22BirdToken creates a 2/2 black Bird creature token with flying.
func NewBlack22BirdToken() *Token {
	tok := NewToken("Black22BirdToken", "2/2 black Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBlackAstartesWarriorToken
// ========================================

// NewBlackAstartesWarriorToken creates a 2/2 black Astartes Warrior creature tokens with menace.
func NewBlackAstartesWarriorToken() *Token {
	tok := NewToken("BlackAstartesWarriorToken", "2/2 black Astartes Warrior creature tokens with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ASTARTES")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewBlackBirdToken
// ========================================

// NewBlackBirdToken creates a token.
func NewBlackBirdToken() *Token {
	tok := NewToken("BlackBirdToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBlackGreenWormToken
// ========================================

// NewBlackGreenWormToken creates a 1/1 black and green Worm creature token.
func NewBlackGreenWormToken() *Token {
	tok := NewToken("BlackGreenWormToken", "1/1 black and green Worm creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WORM")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewBlackWizardToken
// ========================================

// NewBlackWizardToken creates a token.
func NewBlackWizardToken() *Token {
	tok := NewToken("BlackWizardToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewBloodAvatarToken
// ========================================

// NewBloodAvatarToken creates a token.
func NewBloodAvatarToken() *Token {
	tok := NewToken("BloodAvatarToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("AVATAR")
	tok.SetColor(Color{Black: true, Red: true})
	tok.SetPowerToughness(3, 6)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewBloodToken
// ========================================

// NewBloodToken creates a Blood token.
func NewBloodToken() *Token {
	tok := NewToken("BloodToken", "Blood token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("BLOOD")
	return tok
}

// ========================================
// NewBlueBirdToken
// ========================================

// NewBlueBirdToken creates a 1/1 blue Bird creature token with flying.
func NewBlueBirdToken() *Token {
	tok := NewToken("BlueBirdToken", "1/1 blue Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewBlueHorrorToken
// ========================================

// NewBlueHorrorToken creates a token.
func NewBlueHorrorToken() *Token {
	tok := NewToken("BlueHorrorToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Blue: true, Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewBoar2Token
// ========================================

// NewBoar2Token creates a 2/2 green Boar creature token.
func NewBoar2Token() *Token {
	tok := NewToken("Boar2Token", "2/2 green Boar creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BOAR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewBoar3Token
// ========================================

// NewBoar3Token creates a 3/1 green Boar creature token.
func NewBoar3Token() *Token {
	tok := NewToken("Boar3Token", "3/1 green Boar creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BOAR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 1)
	return tok
}

// ========================================
// NewBoarToken
// ========================================

// NewBoarToken creates a 3/3 green Boar creature token.
func NewBoarToken() *Token {
	tok := NewToken("BoarToken", "3/3 green Boar creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BOAR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewBooToken
// ========================================

// NewBooToken creates a Boo, a legendary 1/1 red Hamster creature token with trample and haste.
func NewBooToken() *Token {
	tok := NewToken("BooToken", "Boo, a legendary 1/1 red Hamster creature token with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HAMSTER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewBrainiacToken
// ========================================

// NewBrainiacToken creates a 1/1 red Brainiac creature token.
func NewBrainiacToken() *Token {
	tok := NewToken("BrainiacToken", "1/1 red Brainiac creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BRAINIAC")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewBreedingPitThrullToken
// ========================================

// NewBreedingPitThrullToken creates a 0/1 black Thrull creature token.
func NewBreedingPitThrullToken() *Token {
	tok := NewToken("BreedingPitThrullToken", "0/1 black Thrull creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("THRULL")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewBrokenVisageSpiritToken
// ========================================

// NewBrokenVisageSpiritToken creates a token.
func NewBrokenVisageSpiritToken() *Token {
	tok := NewToken("BrokenVisageSpiritToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewBrudicladTelchorMyrToken
// ========================================

// NewBrudicladTelchorMyrToken creates a 2/1 blue Phyrexian Myr artifact creature token.
func NewBrudicladTelchorMyrToken() *Token {
	tok := NewToken("BrudicladTelchorMyrToken", "2/1 blue Phyrexian Myr artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("MYR")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewButterflyToken
// ========================================

// NewButterflyToken creates a 1/1 green Insect creature token with flying named Butterfly.
func NewButterflyToken() *Token {
	tok := NewToken("ButterflyToken", "1/1 green Insect creature token with flying named Butterfly")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewCallTheSkyBreakerElementalToken
// ========================================

// NewCallTheSkyBreakerElementalToken creates a 5/5 blue and red Elemental creature token with flying.
func NewCallTheSkyBreakerElementalToken() *Token {
	tok := NewToken("CallTheSkyBreakerElementalToken", "5/5 blue and red Elemental creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Blue: true, Red: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewCamaridToken
// ========================================

// NewCamaridToken creates a 1/1 blue Camarid creature tokens.
func NewCamaridToken() *Token {
	tok := NewToken("CamaridToken", "1/1 blue Camarid creature tokens")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAMARID")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewCaribouToken
// ========================================

// NewCaribouToken creates a 0/1 white Caribou creature token.
func NewCaribouToken() *Token {
	tok := NewToken("CaribouToken", "0/1 white Caribou creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CARIBOU")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewCarnivoreToken
// ========================================

// NewCarnivoreToken creates a 3/1 red Beast creature token.
func NewCarnivoreToken() *Token {
	tok := NewToken("CarnivoreToken", "3/1 red Beast creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	return tok
}

// ========================================
// NewCarrionBlackInsectToken
// ========================================

// NewCarrionBlackInsectToken creates a 0/1 black Insect creature token.
func NewCarrionBlackInsectToken() *Token {
	tok := NewToken("CarrionBlackInsectToken", "0/1 black Insect creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewCatBeastToken
// ========================================

// NewCatBeastToken creates a 2/2 white Cat Beast creature token.
func NewCatBeastToken() *Token {
	tok := NewToken("CatBeastToken", "2/2 white Cat Beast creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewCatBirdToken
// ========================================

// NewCatBirdToken creates a 1/1 white Cat Bird creature token with flying.
func NewCatBirdToken() *Token {
	tok := NewToken("CatBirdToken", "1/1 white Cat Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewCatHasteToken
// ========================================

// NewCatHasteToken creates a 2/2 green Cat creature token with haste.
func NewCatHasteToken() *Token {
	tok := NewToken("CatHasteToken", "2/2 green Cat creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewCatSoldierCreatureToken
// ========================================

// NewCatSoldierCreatureToken creates a 1/1 white Cat Soldier creature token with vigilance.
func NewCatSoldierCreatureToken() *Token {
	tok := NewToken("CatSoldierCreatureToken", "1/1 white Cat Soldier creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewCatToken
// ========================================

// NewCatToken creates a 2/2 white Cat creature token.
func NewCatToken() *Token {
	tok := NewToken("CatToken", "2/2 white Cat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewCatWarrior21Token
// ========================================

// NewCatWarrior21Token creates a 2/1 white Cat Warrior creature token.
func NewCatWarrior21Token() *Token {
	tok := NewToken("CatWarrior21Token", "2/1 white Cat Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewCatWarriorToken
// ========================================

// NewCatWarriorToken creates a 2/2 green Cat Warrior creature token with forestwalk.
func NewCatWarriorToken() *Token {
	tok := NewToken("CatWarriorToken", "2/2 green Cat Warrior creature token with forestwalk")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewCentaurEnchantmentCreatureToken
// ========================================

// NewCentaurEnchantmentCreatureToken creates a 3/3 green Centaur enchantment creature token.
func NewCentaurEnchantmentCreatureToken() *Token {
	tok := NewToken("CentaurEnchantmentCreatureToken", "3/3 green Centaur enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("CENTAUR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewCentaurToken
// ========================================

// NewCentaurToken creates a 3/3 green Centaur creature token.
func NewCentaurToken() *Token {
	tok := NewToken("CentaurToken", "3/3 green Centaur creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CENTAUR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewChainersTormentNightmareToken
// ========================================

// NewChainersTormentNightmareToken creates a X/X black Nightmare Horror creature token.
func NewChainersTormentNightmareToken() *Token {
	tok := NewToken("ChainersTormentNightmareToken", "X/X black Nightmare Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("NIGHTMARE")
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewCherubaelToken
// ========================================

// NewCherubaelToken creates a Cherubael, a legendary 4/4 black Demon creature token with flying.
func NewCherubaelToken() *Token {
	tok := NewToken("CherubaelToken", "Cherubael, a legendary 4/4 black Demon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewChocoboToken
// ========================================

// NewChocoboToken creates a token.
func NewChocoboToken() *Token {
	tok := NewToken("ChocoboToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewCitizenGreenWhiteToken
// ========================================

// NewCitizenGreenWhiteToken creates a 1/1 green and white Citizen creature token.
func NewCitizenGreenWhiteToken() *Token {
	tok := NewToken("CitizenGreenWhiteToken", "1/1 green and white Citizen creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CITIZEN")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewCitizenToken
// ========================================

// NewCitizenToken creates a 1/1 white Citizen creature token.
func NewCitizenToken() *Token {
	tok := NewToken("CitizenToken", "1/1 white Citizen creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CITIZEN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewCloudSpriteToken
// ========================================

// NewCloudSpriteToken creates a token.
func NewCloudSpriteToken() *Token {
	tok := NewToken("CloudSpriteToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewClownRobotToken
// ========================================

// NewClownRobotToken creates a 1/1 white Clown Robot artifact creature token.
func NewClownRobotToken() *Token {
	tok := NewToken("ClownRobotToken", "1/1 white Clown Robot artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CLOWN")
	tok.AddSubtype("ROBOT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewClueArtifactToken
// ========================================

// NewClueArtifactToken creates a Clue token.
func NewClueArtifactToken() *Token {
	tok := NewToken("ClueArtifactToken", "Clue token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("CLUE")
	return tok
}

// ========================================
// NewCommodoreGuffToken
// ========================================

// NewCommodoreGuffToken creates a token.
func NewCommodoreGuffToken() *Token {
	tok := NewToken("CommodoreGuffToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewConstruct2Token
// ========================================

// NewConstruct2Token creates a 2/2 colorless Construct artifact creature token.
func NewConstruct2Token() *Token {
	tok := NewToken("Construct2Token", "2/2 colorless Construct artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewConstruct4Token
// ========================================

// NewConstruct4Token creates a 4/4 colorless Construct artifact creature token.
func NewConstruct4Token() *Token {
	tok := NewToken("Construct4Token", "4/4 colorless Construct artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewConstructRedToken
// ========================================

// NewConstructRedToken creates a 3/1 red Construct artifact creature token with haste.
func NewConstructRedToken() *Token {
	tok := NewToken("ConstructRedToken", "3/1 red Construct artifact creature token with haste")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewConstructToken
// ========================================

// NewConstructToken creates a 1/1 colorless Construct artifact creature token.
func NewConstructToken() *Token {
	tok := NewToken("ConstructToken", "1/1 colorless Construct artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewConsumingBlobOozeToken
// ========================================

// NewConsumingBlobOozeToken creates a token.
func NewConsumingBlobOozeToken() *Token {
	tok := NewToken("ConsumingBlobOozeToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewCordycepsInfectedToken
// ========================================

// NewCordycepsInfectedToken creates a 1/1 black Fungus Zombie creature token named Cordyceps Infected.
func NewCordycepsInfectedToken() *Token {
	tok := NewToken("CordycepsInfectedToken", "1/1 black Fungus Zombie creature token named Cordyceps Infected")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FUNGUS")
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewCorpseweftZombieToken
// ========================================

// NewCorpseweftZombieToken creates a X/X black Zombie Horror creature token, where X is twice the number of cards exiled this way.
func NewCorpseweftZombieToken() *Token {
	tok := NewToken("CorpseweftZombieToken", "X/X black Zombie Horror creature token, where X is twice the number of cards exiled this way")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewCorruptedZendikonOozeToken
// ========================================

// NewCorruptedZendikonOozeToken creates a 3/3 black Ooze creature.
func NewCorruptedZendikonOozeToken() *Token {
	tok := NewToken("CorruptedZendikonOozeToken", "3/3 black Ooze creature")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewCrabToken
// ========================================

// NewCrabToken creates a 0/3 blue Crab creature token.
func NewCrabToken() *Token {
	tok := NewToken("CrabToken", "0/3 blue Crab creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CRAB")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 3)
	return tok
}

// ========================================
// NewCragflameToken
// ========================================

// NewCragflameToken creates a token.
func NewCragflameToken() *Token {
	tok := NewToken("CragflameToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("EQUIPMENT")
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewCreatureToken
// ========================================

// NewCreatureToken creates a token.
func NewCreatureToken() *Token {
	tok := NewToken("CreatureToken", "")
	tok.AddCardType(CardTypeCreature)
	return tok
}

// ========================================
// NewCrestedSunmareToken
// ========================================

// NewCrestedSunmareToken creates a 5/5 white Horse creature token.
func NewCrestedSunmareToken() *Token {
	tok := NewToken("CrestedSunmareToken", "5/5 white Horse creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORSE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewCribSwapShapeshifterWhiteToken
// ========================================

// NewCribSwapShapeshifterWhiteToken creates a 1/1 colorless Shapeshifter creature token with changeling.
func NewCribSwapShapeshifterWhiteToken() *Token {
	tok := NewToken("CribSwapShapeshifterWhiteToken", "1/1 colorless Shapeshifter creature token with changeling")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHAPESHIFTER")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewCursedRoleToken
// ========================================

// NewCursedRoleToken creates a Cursed Role token.
func NewCursedRoleToken() *Token {
	tok := NewToken("CursedRoleToken", "Cursed Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	return tok
}

// ========================================
// NewCustomIllusionToken
// ========================================

// NewCustomIllusionToken creates a X/X blue Illusion creature token.
func NewCustomIllusionToken() *Token {
	tok := NewToken("CustomIllusionToken", "X/X blue Illusion creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	return tok
}

// ========================================
// NewDalekToken
// ========================================

// NewDalekToken creates a 3/3 black Dalek artifact creature token with menace.
func NewDalekToken() *Token {
	tok := NewToken("DalekToken", "3/3 black Dalek artifact creature token with menace")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DALEK")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewDarettiConstructToken
// ========================================

// NewDarettiConstructToken creates a 1/1 colorless Construct artifact creature token with defender.
func NewDarettiConstructToken() *Token {
	tok := NewToken("DarettiConstructToken", "1/1 colorless Construct artifact creature token with defender")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewDarkstarToken
// ========================================

// NewDarkstarToken creates a Darkstar, a legendary 2/2 white and black Dog creature token.
func NewDarkstarToken() *Token {
	tok := NewToken("DarkstarToken", "Darkstar, a legendary 2/2 white and black Dog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewDaxosSpiritToken
// ========================================

// NewDaxosSpiritToken creates a token.
func NewDaxosSpiritToken() *Token {
	tok := NewToken("DaxosSpiritToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewDeadlyGrubInsectToken
// ========================================

// NewDeadlyGrubInsectToken creates a 6/1 green Insect creature token with shroud.
func NewDeadlyGrubInsectToken() *Token {
	tok := NewToken("DeadlyGrubInsectToken", "6/1 green Insect creature token with shroud")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(6, 1)
	return tok
}

// ========================================
// NewDeathpactAngelToken
// ========================================

// NewDeathpactAngelToken creates a token.
func NewDeathpactAngelToken() *Token {
	tok := NewToken("DeathpactAngelToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CLERIC")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewDeathtouchRatToken
// ========================================

// NewDeathtouchRatToken creates a 1/1 black Rat creature token with deathtouch.
func NewDeathtouchRatToken() *Token {
	tok := NewToken("DeathtouchRatToken", "1/1 black Rat creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewDeathtouchSnakeToken
// ========================================

// NewDeathtouchSnakeToken creates a 1/1 green Snake creature token with deathtouch.
func NewDeathtouchSnakeToken() *Token {
	tok := NewToken("DeathtouchSnakeToken", "1/1 green Snake creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAKE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewDefenderPlantToken
// ========================================

// NewDefenderPlantToken creates a 0/2 green Plant creature token with defender.
func NewDefenderPlantToken() *Token {
	tok := NewToken("DefenderPlantToken", "0/2 green Plant creature token with defender")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PLANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 2)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewDemon33Token
// ========================================

// NewDemon33Token creates a 3/3 black Demon creature token.
func NewDemon33Token() *Token {
	tok := NewToken("Demon33Token", "3/3 black Demon creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewDemon66Token
// ========================================

// NewDemon66Token creates a 6/6 black Demon creature token with flying.
func NewDemon66Token() *Token {
	tok := NewToken("Demon66Token", "6/6 black Demon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(6, 6)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDemonBerserkerToken
// ========================================

// NewDemonBerserkerToken creates a 2/3 red Demon Berserker creature token with menace.
func NewDemonBerserkerToken() *Token {
	tok := NewToken("DemonBerserkerToken", "2/3 red Demon Berserker creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.AddSubtype("BERSERKER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 3)
	return tok
}

// ========================================
// NewDemonFlyingToken
// ========================================

// NewDemonFlyingToken creates a X/X black Demon creature token with flying.
func NewDemonFlyingToken() *Token {
	tok := NewToken("DemonFlyingToken", "X/X black Demon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDemonToken
// ========================================

// NewDemonToken creates a 5/5 black Demon creature token with flying.
func NewDemonToken() *Token {
	tok := NewToken("DemonToken", "5/5 black Demon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDeserterToken
// ========================================

// NewDeserterToken creates a 0/1 white Deserter creature token.
func NewDeserterToken() *Token {
	tok := NewToken("DeserterToken", "0/1 white Deserter creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DESERTER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewDetectiveToken
// ========================================

// NewDetectiveToken creates a 2/2 white and blue Detective creature token.
func NewDetectiveToken() *Token {
	tok := NewToken("DetectiveToken", "2/2 white and blue Detective creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DETECTIVE")
	tok.SetColor(Color{White: true, Blue: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewDevastatingSummonsElementalToken
// ========================================

// NewDevastatingSummonsElementalToken creates a X/X red Elemental creature.
func NewDevastatingSummonsElementalToken() *Token {
	tok := NewToken("DevastatingSummonsElementalToken", "X/X red Elemental creature")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	return tok
}

// ========================================
// NewDevilToken
// ========================================

// NewDevilToken creates a token.
func NewDevilToken() *Token {
	tok := NewToken("DevilToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEVIL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewDinDragonToken
// ========================================

// NewDinDragonToken creates a 4/4 red Dinosaur Dragon creature token with flying.
func NewDinDragonToken() *Token {
	tok := NewToken("DinDragonToken", "4/4 red Dinosaur Dragon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDinOfTheFireherdToken
// ========================================

// NewDinOfTheFireherdToken creates a 5/5 black and red Elemental creature.
func NewDinOfTheFireherdToken() *Token {
	tok := NewToken("DinOfTheFireherdToken", "5/5 black and red Elemental creature")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Black: true, Red: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewDinosaur31Token
// ========================================

// NewDinosaur31Token creates a 3/1 red Dinosaur creature token.
func NewDinosaur31Token() *Token {
	tok := NewToken("Dinosaur31Token", "3/1 red Dinosaur creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	return tok
}

// ========================================
// NewDinosaurBeastToken
// ========================================

// NewDinosaurBeastToken creates a X/X green Dinosaur Beast creature token with trample.
func NewDinosaurBeastToken() *Token {
	tok := NewToken("DinosaurBeastToken", "X/X green Dinosaur Beast creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewDinosaurCatToken
// ========================================

// NewDinosaurCatToken creates a 2/2 red and white Dinosaur Cat creature token.
func NewDinosaurCatToken() *Token {
	tok := NewToken("DinosaurCatToken", "2/2 red and white Dinosaur Cat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.AddSubtype("CAT")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewDinosaurEggToken
// ========================================

// NewDinosaurEggToken creates a 0/1 green Dinosaur Egg creature token.
func NewDinosaurEggToken() *Token {
	tok := NewToken("DinosaurEggToken", "0/1 green Dinosaur Egg creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.AddSubtype("EGG")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewDinosaurFlyingHasteToken
// ========================================

// NewDinosaurFlyingHasteToken creates a 2/2 red and white Dinosaur creature token with flying and haste.
func NewDinosaurFlyingHasteToken() *Token {
	tok := NewToken("DinosaurFlyingHasteToken", "2/2 red and white Dinosaur creature token with flying and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewDinosaurHasteToken
// ========================================

// NewDinosaurHasteToken creates a 1/1 red Dinosaur creature token with haste.
func NewDinosaurHasteToken() *Token {
	tok := NewToken("DinosaurHasteToken", "1/1 red Dinosaur creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewDinosaurToken
// ========================================

// NewDinosaurToken creates a 3/3 green Dinosaur creature token with trample.
func NewDinosaurToken() *Token {
	tok := NewToken("DinosaurToken", "3/3 green Dinosaur creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewDinosaurVanillaToken
// ========================================

// NewDinosaurVanillaToken creates a 3/3 green Dinosaur creature token.
func NewDinosaurVanillaToken() *Token {
	tok := NewToken("DinosaurVanillaToken", "3/3 green Dinosaur creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewDinosaurXXToken
// ========================================

// NewDinosaurXXToken creates a X/X green Dinosaur creature token with trample.
func NewDinosaurXXToken() *Token {
	tok := NewToken("DinosaurXXToken", "X/X green Dinosaur creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{Green: true})
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewDjinnMonkToken
// ========================================

// NewDjinnMonkToken creates a 2/2 blue Djinn Monk creature token with flying.
func NewDjinnMonkToken() *Token {
	tok := NewToken("DjinnMonkToken", "2/2 blue Djinn Monk creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DJINN")
	tok.AddSubtype("MONK")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDjinnToken
// ========================================

// NewDjinnToken creates a 5/5 colorless Djinn artifact creature token with flying.
func NewDjinnToken() *Token {
	tok := NewToken("DjinnToken", "5/5 colorless Djinn artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DJINN")
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDogIllusionToken
// ========================================

// NewDogIllusionToken creates a token.
func NewDogIllusionToken() *Token {
	tok := NewToken("DogIllusionToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewDogVigilanceToken
// ========================================

// NewDogVigilanceToken creates a 3/1 green Dog creature token with vigilance.
func NewDogVigilanceToken() *Token {
	tok := NewToken("DogVigilanceToken", "3/1 green Dog creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewDokaiWeaverofLifeToken
// ========================================

// NewDokaiWeaverofLifeToken creates a X/X green Elemental creature token, where X is the number of lands you control.
func NewDokaiWeaverofLifeToken() *Token {
	tok := NewToken("DokaiWeaverofLifeToken", "X/X green Elemental creature token, where X is the number of lands you control")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewDoomedArtisanToken
// ========================================

// NewDoomedArtisanToken creates a token.
func NewDoomedArtisanToken() *Token {
	tok := NewToken("DoomedArtisanToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SCULPTURE")
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewDorotheasRetributionSpiritToken
// ========================================

// NewDorotheasRetributionSpiritToken creates a 4/4 white Spirit creature token with flying.
func NewDorotheasRetributionSpiritToken() *Token {
	tok := NewToken("DorotheasRetributionSpiritToken", "4/4 white Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDoublestrikeSamuraiToken
// ========================================

// NewDoublestrikeSamuraiToken creates a 2/2 white Samurai creature token with double strike..
func NewDoublestrikeSamuraiToken() *Token {
	tok := NewToken("DoublestrikeSamuraiToken", "2/2 white Samurai creature token with double strike.")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SAMURAI")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("double strike")
	return tok
}

// ========================================
// NewDragonBroodmotherDragonToken
// ========================================

// NewDragonBroodmotherDragonToken creates a 1/1 red and green Dragon creature token with flying and devour 2.
func NewDragonBroodmotherDragonToken() *Token {
	tok := NewToken("DragonBroodmotherDragonToken", "1/1 red and green Dragon creature token with flying and devour 2")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDragonEggDragonToken
// ========================================

// NewDragonEggDragonToken creates a token.
func NewDragonEggDragonToken() *Token {
	tok := NewToken("DragonEggDragonToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDragonElementalToken
// ========================================

// NewDragonElementalToken creates a 4/4 red Dragon Elemental creature token with flying and prowess.
func NewDragonElementalToken() *Token {
	tok := NewToken("DragonElementalToken", "4/4 red Dragon Elemental creature token with flying and prowess")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDragonIllusionToken
// ========================================

// NewDragonIllusionToken creates a X/X red Dragon Illusion creature token with flying and haste.
func NewDragonIllusionToken() *Token {
	tok := NewToken("DragonIllusionToken", "X/X red Dragon Illusion creature token with flying and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Red: true})
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewDragonMenaceAndStealArtifactToken
// ========================================

// NewDragonMenaceAndStealArtifactToken creates a token.
func NewDragonMenaceAndStealArtifactToken() *Token {
	tok := NewToken("DragonMenaceAndStealArtifactToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Black: true, Red: true})
	tok.SetPowerToughness(6, 6)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDragonSpiritToken
// ========================================

// NewDragonSpiritToken creates a 5/5 red Dragon Spirit creature token with flying.
func NewDragonSpiritToken() *Token {
	tok := NewToken("DragonSpiritToken", "5/5 red Dragon Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDragonToken
// ========================================

// NewDragonToken creates a 4/4 red Dragon creature token with flying.
func NewDragonToken() *Token {
	tok := NewToken("DragonToken", "4/4 red Dragon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDrakeToken
// ========================================

// NewDrakeToken creates a 2/2 blue Drake creature token with flying.
func NewDrakeToken() *Token {
	tok := NewToken("DrakeToken", "2/2 blue Drake creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAKE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewDroidToken
// ========================================

// NewDroidToken creates a 1/1 colorless Droid creature token.
func NewDroidToken() *Token {
	tok := NewToken("DroidToken", "1/1 colorless Droid creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DROID")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewDroneToken
// ========================================

// NewDroneToken creates a token.
func NewDroneToken() *Token {
	tok := NewToken("DroneToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRONE")
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewDuneBroodNephilimToken
// ========================================

// NewDuneBroodNephilimToken creates a 1/1 colorless Sand creature token.
func NewDuneBroodNephilimToken() *Token {
	tok := NewToken("DuneBroodNephilimToken", "1/1 colorless Sand creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SAND")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewDwarfBerserkerToken
// ========================================

// NewDwarfBerserkerToken creates a 2/1 red Dwarf Berserker creature token.
func NewDwarfBerserkerToken() *Token {
	tok := NewToken("DwarfBerserkerToken", "2/1 red Dwarf Berserker creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DWARF")
	tok.AddSubtype("BERSERKER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewDwarfToken
// ========================================

// NewDwarfToken creates a 1/1 red Dwarf creature token.
func NewDwarfToken() *Token {
	tok := NewToken("DwarfToken", "1/1 red Dwarf creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DWARF")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewEdgarMarkovToken
// ========================================

// NewEdgarMarkovToken creates a 1/1 black Vampire creature token.
func NewEdgarMarkovToken() *Token {
	tok := NewToken("EdgarMarkovToken", "1/1 black Vampire creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewEdgarMarkovsCoffinVampireToken
// ========================================

// NewEdgarMarkovsCoffinVampireToken creates a 1/1 white and black Vampire creature token with lifelink.
func NewEdgarMarkovsCoffinVampireToken() *Token {
	tok := NewToken("EdgarMarkovsCoffinVampireToken", "1/1 white and black Vampire creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewEldraziAngelToken
// ========================================

// NewEldraziAngelToken creates a 4/4 colorless Eldrazi Angel creature token with flying and vigilance.
func NewEldraziAngelToken() *Token {
	tok := NewToken("EldraziAngelToken", "4/4 colorless Eldrazi Angel creature token with flying and vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewEldraziAnnihilatorToken
// ========================================

// NewEldraziAnnihilatorToken creates a 7/7 colorless Eldrazi creature token with annihilator 1.
func NewEldraziAnnihilatorToken() *Token {
	tok := NewToken("EldraziAnnihilatorToken", "7/7 colorless Eldrazi creature token with annihilator 1")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELDRAZI")
	tok.SetPowerToughness(7, 7)
	return tok
}

// ========================================
// NewEldraziHorrorToken
// ========================================

// NewEldraziHorrorToken creates a 3/2 colorless Eldrazi Horror creature token.
func NewEldraziHorrorToken() *Token {
	tok := NewToken("EldraziHorrorToken", "3/2 colorless Eldrazi Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELDRAZI")
	tok.AddSubtype("HORROR")
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewEldraziScionToken
// ========================================

// NewEldraziScionToken creates a token.
func NewEldraziScionToken() *Token {
	tok := NewToken("EldraziScionToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELDRAZI")
	tok.AddSubtype("SCION")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewEldraziSliverToken
// ========================================

// NewEldraziSliverToken creates a token.
func NewEldraziSliverToken() *Token {
	tok := NewToken("EldraziSliverToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELDRAZI")
	tok.AddSubtype("SLIVER")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewEldraziSpawnToken
// ========================================

// NewEldraziSpawnToken creates a token.
func NewEldraziSpawnToken() *Token {
	tok := NewToken("EldraziSpawnToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELDRAZI")
	tok.AddSubtype("SPAWN")
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewEldraziToken
// ========================================

// NewEldraziToken creates a 10/10 colorless Eldrazi creature token.
func NewEldraziToken() *Token {
	tok := NewToken("EldraziToken", "10/10 colorless Eldrazi creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELDRAZI")
	tok.SetPowerToughness(10, 10)
	return tok
}

// ========================================
// NewElemental11BlueRedToken
// ========================================

// NewElemental11BlueRedToken creates a 1/1 blue and red Elemental creature token.
func NewElemental11BlueRedToken() *Token {
	tok := NewToken("Elemental11BlueRedToken", "1/1 blue and red Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Blue: true, Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewElemental11HasteToken
// ========================================

// NewElemental11HasteToken creates a 1/1 red Elemental creature token with haste.
func NewElemental11HasteToken() *Token {
	tok := NewToken("Elemental11HasteToken", "1/1 red Elemental creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewElemental21TrampleHasteToken
// ========================================

// NewElemental21TrampleHasteToken creates a 2/1 red Elemental creature token with trample and haste.
func NewElemental21TrampleHasteToken() *Token {
	tok := NewToken("Elemental21TrampleHasteToken", "2/1 red Elemental creature token with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewElemental31TrampleHasteToken
// ========================================

// NewElemental31TrampleHasteToken creates a 3/1 red Elemental creature token with trample and haste.
func NewElemental31TrampleHasteToken() *Token {
	tok := NewToken("Elemental31TrampleHasteToken", "3/1 red Elemental creature token with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewElemental44GreenToken
// ========================================

// NewElemental44GreenToken creates a 4/4 green Elemental creature token.
func NewElemental44GreenToken() *Token {
	tok := NewToken("Elemental44GreenToken", "4/4 green Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewElemental44Token
// ========================================

// NewElemental44Token creates a 4/4 blue and red Elemental creature token.
func NewElemental44Token() *Token {
	tok := NewToken("Elemental44Token", "4/4 blue and red Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Blue: true, Red: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewElemental44WUToken
// ========================================

// NewElemental44WUToken creates a 4/4 white and blue Elemental creature token.
func NewElemental44WUToken() *Token {
	tok := NewToken("Elemental44WUToken", "4/4 white and blue Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{White: true, Blue: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewElementalAllColorsToken
// ========================================

// NewElementalAllColorsToken creates a 2/2 Elemental creature token that's all colors.
func NewElementalAllColorsToken() *Token {
	tok := NewToken("ElementalAllColorsToken", "2/2 Elemental creature token that's all colors")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{White: true, Blue: true, Black: true, Red: true, Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewElementalCatToken
// ========================================

// NewElementalCatToken creates a 1/1 red Elemental Cat creature token.
func NewElementalCatToken() *Token {
	tok := NewToken("ElementalCatToken", "1/1 red Elemental Cat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.AddSubtype("CAT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewElementalShamanToken
// ========================================

// NewElementalShamanToken creates a 3/1 red Elemental Shaman creature token.
func NewElementalShamanToken() *Token {
	tok := NewToken("ElementalShamanToken", "3/1 red Elemental Shaman creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.AddSubtype("SHAMAN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewElementalXXGreenToken
// ========================================

// NewElementalXXGreenToken creates a X/X green Elemental creature token.
func NewElementalXXGreenToken() *Token {
	tok := NewToken("ElementalXXGreenToken", "X/X green Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewElephant55Token
// ========================================

// NewElephant55Token creates a 5/5 green Elephant creature token.
func NewElephant55Token() *Token {
	tok := NewToken("Elephant55Token", "5/5 green Elephant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEPHANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewElephantResurgenceToken
// ========================================

// NewElephantResurgenceToken creates a token.
func NewElephantResurgenceToken() *Token {
	tok := NewToken("ElephantResurgenceToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEPHANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewElephantToken
// ========================================

// NewElephantToken creates a 3/3 green Elephant creature token.
func NewElephantToken() *Token {
	tok := NewToken("ElephantToken", "3/3 green Elephant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEPHANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewElfDruidToken
// ========================================

// NewElfDruidToken creates a token.
func NewElfDruidToken() *Token {
	tok := NewToken("ElfDruidToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELF")
	tok.AddSubtype("DRUID")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewElfKnightToken
// ========================================

// NewElfKnightToken creates a 2/2 green and white Elf Knight creature token with vigilance.
func NewElfKnightToken() *Token {
	tok := NewToken("ElfKnightToken", "2/2 green and white Elf Knight creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELF")
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewElfWarriorToken
// ========================================

// NewElfWarriorToken creates a 1/1 green Elf Warrior creature token.
func NewElfWarriorToken() *Token {
	tok := NewToken("ElfWarriorToken", "1/1 green Elf Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELF")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewElkToken
// ========================================

// NewElkToken creates a 3/3 green Elk creature token.
func NewElkToken() *Token {
	tok := NewToken("ElkToken", "3/3 green Elk creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELK")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewEmptyToken
// ========================================

// NewEmptyToken creates a token.
func NewEmptyToken() *Token {
	tok := NewToken("EmptyToken", "")
	return tok
}

// ========================================
// NewEnchantmentBirdToken
// ========================================

// NewEnchantmentBirdToken creates a 2/2 blue Bird enchantment creature token with flying.
func NewEnchantmentBirdToken() *Token {
	tok := NewToken("EnchantmentBirdToken", "2/2 blue Bird enchantment creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewErrandOfDutyKnightToken
// ========================================

// NewErrandOfDutyKnightToken creates a 1/1 white Knight creature token with banding.
func NewErrandOfDutyKnightToken() *Token {
	tok := NewToken("ErrandOfDutyKnightToken", "1/1 white Knight creature token with banding")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewEtheriumCellToken
// ========================================

// NewEtheriumCellToken creates a token.
func NewEtheriumCellToken() *Token {
	tok := NewToken("EtheriumCellToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewEverywhereToken
// ========================================

// NewEverywhereToken creates a colorless land token named Everywhere that is every basic land type.
func NewEverywhereToken() *Token {
	tok := NewToken("EverywhereToken", "colorless land token named Everywhere that is every basic land type")
	tok.AddCardType(CardTypeLand)
	return tok
}

// ========================================
// NewEwokToken
// ========================================

// NewEwokToken creates a 1/1 green Ewok creature tokens.
func NewEwokToken() *Token {
	tok := NewToken("EwokToken", "1/1 green Ewok creature tokens")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("EWOK")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewExpansionSymbolToken
// ========================================

// NewExpansionSymbolToken creates a 1/1 colorless Expansion-Symbol creature token.
func NewExpansionSymbolToken() *Token {
	tok := NewToken("ExpansionSymbolToken", "1/1 colorless Expansion-Symbol creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("EXPANSION_SYMBOL")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFableOfTheMirrorBreakerToken
// ========================================

// NewFableOfTheMirrorBreakerToken creates a token.
func NewFableOfTheMirrorBreakerToken() *Token {
	tok := NewToken("FableOfTheMirrorBreakerToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("SHAMAN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewFaerieBlockFliersToken
// ========================================

// NewFaerieBlockFliersToken creates a token.
func NewFaerieBlockFliersToken() *Token {
	tok := NewToken("FaerieBlockFliersToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewFaerieBlueBlackToken
// ========================================

// NewFaerieBlueBlackToken creates a 1/1 blue and black Faerie Rogue creature token with flying.
func NewFaerieBlueBlackToken() *Token {
	tok := NewToken("FaerieBlueBlackToken", "1/1 blue and black Faerie Rogue creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Blue: true, Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewFaerieDragonToken
// ========================================

// NewFaerieDragonToken creates a 1/1 blue Faerie Dragon creature token with flying.
func NewFaerieDragonToken() *Token {
	tok := NewToken("FaerieDragonToken", "1/1 blue Faerie Dragon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewFaerieRogueToken
// ========================================

// NewFaerieRogueToken creates a 1/1 black Faerie Rogue creature token with flying.
func NewFaerieRogueToken() *Token {
	tok := NewToken("FaerieRogueToken", "1/1 black Faerie Rogue creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewFaerieToken
// ========================================

// NewFaerieToken creates a 1/1 blue Faerie creature token with flying.
func NewFaerieToken() *Token {
	tok := NewToken("FaerieToken", "1/1 blue Faerie creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewFeatherToken
// ========================================

// NewFeatherToken creates a token.
func NewFeatherToken() *Token {
	tok := NewToken("FeatherToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.SetColor(Color{Red: true})
	return tok
}

// ========================================
// NewFesteringGoblinToken
// ========================================

// NewFesteringGoblinToken creates a token.
func NewFesteringGoblinToken() *Token {
	tok := NewToken("FesteringGoblinToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("GOBLIN")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFirstMateRagavanToken
// ========================================

// NewFirstMateRagavanToken creates a First Mate Ragavan, a legendary 2/1 red Monkey Pirate creature token.
func NewFirstMateRagavanToken() *Token {
	tok := NewToken("FirstMateRagavanToken", "First Mate Ragavan, a legendary 2/1 red Monkey Pirate creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MONKEY")
	tok.AddSubtype("PIRATE")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewFishNoAbilityToken
// ========================================

// NewFishNoAbilityToken creates a 1/1 blue Fish creature token.
func NewFishNoAbilityToken() *Token {
	tok := NewToken("FishNoAbilityToken", "1/1 blue Fish creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FISH")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFishToken
// ========================================

// NewFishToken creates a token.
func NewFishToken() *Token {
	tok := NewToken("FishToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FISH")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFoodToken
// ========================================

// NewFoodToken creates a Food token.
func NewFoodToken() *Token {
	tok := NewToken("FoodToken", "Food token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("FOOD")
	return tok
}

// ========================================
// NewForestDryadToken
// ========================================

// NewForestDryadToken creates a 1/1 green Forest Dryad land creature token.
func NewForestDryadToken() *Token {
	tok := NewToken("ForestDryadToken", "1/1 green Forest Dryad land creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeLand)
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewForlornPseudammaZombieToken
// ========================================

// NewForlornPseudammaZombieToken creates a 2/2 black Zombie enchantment creature token.
func NewForlornPseudammaZombieToken() *Token {
	tok := NewToken("ForlornPseudammaZombieToken", "2/2 black Zombie enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewFox22VigilanceToken
// ========================================

// NewFox22VigilanceToken creates a 2/2 white Fox creature token with vigilance.
func NewFox22VigilanceToken() *Token {
	tok := NewToken("Fox22VigilanceToken", "2/2 white Fox creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FOX")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewFractalToken
// ========================================

// NewFractalToken creates a 0/0 green and blue Fractal creature token.
func NewFractalToken() *Token {
	tok := NewToken("FractalToken", "0/0 green and blue Fractal creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FRACTAL")
	tok.SetColor(Color{Blue: true, Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewFrogGreenToken
// ========================================

// NewFrogGreenToken creates a 1/1 green Frog creature token.
func NewFrogGreenToken() *Token {
	tok := NewToken("FrogGreenToken", "1/1 green Frog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FROG")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFrogLizardToken
// ========================================

// NewFrogLizardToken creates a 3/3 green Frog Lizard creature token.
func NewFrogLizardToken() *Token {
	tok := NewToken("FrogLizardToken", "3/3 green Frog Lizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FROG")
	tok.AddSubtype("LIZARD")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewFrogToken
// ========================================

// NewFrogToken creates a 1/1 blue Frog creature token.
func NewFrogToken() *Token {
	tok := NewToken("FrogToken", "1/1 blue Frog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FROG")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFungusBeastToken
// ========================================

// NewFungusBeastToken creates a 4/4 green Fungus Beast creature token with trample.
func NewFungusBeastToken() *Token {
	tok := NewToken("FungusBeastToken", "4/4 green Fungus Beast creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FUNGUS")
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewFungusCantBlockToken
// ========================================

// NewFungusCantBlockToken creates a token.
func NewFungusCantBlockToken() *Token {
	tok := NewToken("FungusCantBlockToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FUNGUS")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewFungusDinosaurToken
// ========================================

// NewFungusDinosaurToken creates a X/X green Fungus Dinosaur creature token.
func NewFungusDinosaurToken() *Token {
	tok := NewToken("FungusDinosaurToken", "X/X green Fungus Dinosaur creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FUNGUS")
	tok.AddSubtype("DINOSAUR")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewGargoyleToken
// ========================================

// NewGargoyleToken creates a 3/4 colorless Gargoyle artifact creature token with flying.
func NewGargoyleToken() *Token {
	tok := NewToken("GargoyleToken", "3/4 colorless Gargoyle artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GARGOYLE")
	tok.SetPowerToughness(3, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewGarrukApexPredatorBeastToken
// ========================================

// NewGarrukApexPredatorBeastToken creates a 3/3 black Beast creature token with deathtouch.
func NewGarrukApexPredatorBeastToken() *Token {
	tok := NewToken("GarrukApexPredatorBeastToken", "3/3 black Beast creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewGarrukCursedHuntsmanToken
// ========================================

// NewGarrukCursedHuntsmanToken creates a token.
func NewGarrukCursedHuntsmanToken() *Token {
	tok := NewToken("GarrukCursedHuntsmanToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewGeminiEngineTwinToken
// ========================================

// NewGeminiEngineTwinToken creates a colorless Construct artifact creature token named Twin that's attacking. Its power is equal to Gemini Engine's power and its toughness is equal to Gemini Engine's toughness..
func NewGeminiEngineTwinToken() *Token {
	tok := NewToken("GeminiEngineTwinToken", "colorless Construct artifact creature token named Twin that's attacking. Its power is equal to Gemini Engine's power and its toughness is equal to Gemini Engine's toughness.")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	return tok
}

// ========================================
// NewGiantBaitingGiantWarriorToken
// ========================================

// NewGiantBaitingGiantWarriorToken creates a 4/4 red and green Giant Warrior creature token with haste.
func NewGiantBaitingGiantWarriorToken() *Token {
	tok := NewToken("GiantBaitingGiantWarriorToken", "4/4 red and green Giant Warrior creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GIANT")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewGiantBirdToken
// ========================================

// NewGiantBirdToken creates a 4/4 red Giant Bird creature token.
func NewGiantBirdToken() *Token {
	tok := NewToken("GiantBirdToken", "4/4 red Giant Bird creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GIANT")
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewGiantOpportunityToken
// ========================================

// NewGiantOpportunityToken creates a 7/7 green Giant creature token.
func NewGiantOpportunityToken() *Token {
	tok := NewToken("GiantOpportunityToken", "7/7 green Giant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GIANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(7, 7)
	return tok
}

// ========================================
// NewGiantToken
// ========================================

// NewGiantToken creates a 4/4 red Giant creature token.
func NewGiantToken() *Token {
	tok := NewToken("GiantToken", "4/4 red Giant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GIANT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewGiantWarriorToken
// ========================================

// NewGiantWarriorToken creates a 5/5 white Giant Warrior creature token.
func NewGiantWarriorToken() *Token {
	tok := NewToken("GiantWarriorToken", "5/5 white Giant Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GIANT")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewGiantWizardToken
// ========================================

// NewGiantWizardToken creates a 4/4 blue Giant Wizard creature token.
func NewGiantWizardToken() *Token {
	tok := NewToken("GiantWizardToken", "4/4 blue Giant Wizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GIANT")
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewGlimmerToken
// ========================================

// NewGlimmerToken creates a 1/1 white Glimmer enchantment creature token.
func NewGlimmerToken() *Token {
	tok := NewToken("GlimmerToken", "1/1 white Glimmer enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("GLIMMER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGnomeSoldierStarStarToken
// ========================================

// NewGnomeSoldierStarStarToken creates a token.
func NewGnomeSoldierStarStarToken() *Token {
	tok := NewToken("GnomeSoldierStarStarToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GNOME")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	return tok
}

// ========================================
// NewGnomeToken
// ========================================

// NewGnomeToken creates a 1/1 colorless Gnome artifact creature token.
func NewGnomeToken() *Token {
	tok := NewToken("GnomeToken", "1/1 colorless Gnome artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GNOME")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGoatToken
// ========================================

// NewGoatToken creates a 0/1 white Goat creature token.
func NewGoatToken() *Token {
	tok := NewToken("GoatToken", "0/1 white Goat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOAT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewGoblinRogueToken
// ========================================

// NewGoblinRogueToken creates a 1/1 black Goblin Rogue creature token.
func NewGoblinRogueToken() *Token {
	tok := NewToken("GoblinRogueToken", "1/1 black Goblin Rogue creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGoblinScoutsToken
// ========================================

// NewGoblinScoutsToken creates a 1/1 red Goblin Scout creature tokens with mountainwalk.
func NewGoblinScoutsToken() *Token {
	tok := NewToken("GoblinScoutsToken", "1/1 red Goblin Scout creature tokens with mountainwalk")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("SCOUT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGoblinSoldierToken
// ========================================

// NewGoblinSoldierToken creates a 1/1 red and white Goblin Soldier creature token.
func NewGoblinSoldierToken() *Token {
	tok := NewToken("GoblinSoldierToken", "1/1 red and white Goblin Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGoblinToken
// ========================================

// NewGoblinToken creates a 1/1 red Goblin creature token.
func NewGoblinToken() *Token {
	tok := NewToken("GoblinToken", "1/1 red Goblin creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewGoblinWarriorToken
// ========================================

// NewGoblinWarriorToken creates a 1/1 red and green Goblin Warrior creature token.
func NewGoblinWarriorToken() *Token {
	tok := NewToken("GoblinWarriorToken", "1/1 red and green Goblin Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGoblinWizardToken
// ========================================

// NewGoblinWizardToken creates a 1/1 red Goblin Wizard creature token with prowess.
func NewGoblinWizardToken() *Token {
	tok := NewToken("GoblinWizardToken", "1/1 red Goblin Wizard creature token with prowess")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGodEternalOketraToken
// ========================================

// NewGodEternalOketraToken creates a 4/4 black Zombie Warrior creature token with vigilance.
func NewGodEternalOketraToken() *Token {
	tok := NewToken("GodEternalOketraToken", "4/4 black Zombie Warrior creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewGodFavoredGeneralSoldierToken
// ========================================

// NewGodFavoredGeneralSoldierToken creates a 1/1 white Soldier enchantment creature token.
func NewGodFavoredGeneralSoldierToken() *Token {
	tok := NewToken("GodFavoredGeneralSoldierToken", "1/1 white Soldier enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGodSireBeastToken
// ========================================

// NewGodSireBeastToken creates a 8/8 Beast creature token that's red, green, and white.
func NewGodSireBeastToken() *Token {
	tok := NewToken("GodSireBeastToken", "8/8 Beast creature token that's red, green, and white")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{White: true, Red: true, Green: true})
	tok.SetPowerToughness(8, 8)
	return tok
}

// ========================================
// NewGoldForgeGarrisonGolemToken
// ========================================

// NewGoldForgeGarrisonGolemToken creates a 4/4 colorless Golem artifact creature token.
func NewGoldForgeGarrisonGolemToken() *Token {
	tok := NewToken("GoldForgeGarrisonGolemToken", "4/4 colorless Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewGoldToken
// ========================================

// NewGoldToken creates a Gold token.
func NewGoldToken() *Token {
	tok := NewToken("GoldToken", "Gold token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("GOLD")
	return tok
}

// ========================================
// NewGoldmeadowHarrierToken
// ========================================

// NewGoldmeadowHarrierToken creates a token.
func NewGoldmeadowHarrierToken() *Token {
	tok := NewToken("GoldmeadowHarrierToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KITHKIN")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGolemFlyingToken
// ========================================

// NewGolemFlyingToken creates a 3/3 colorless Golem artifact creature token with flying.
func NewGolemFlyingToken() *Token {
	tok := NewToken("GolemFlyingToken", "3/3 colorless Golem artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewGolemToken
// ========================================

// NewGolemToken creates a 3/3 colorless Golem artifact creature token.
func NewGolemToken() *Token {
	tok := NewToken("GolemToken", "3/3 colorless Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewGolemTrampleToken
// ========================================

// NewGolemTrampleToken creates a 3/3 colorless Golem artifact creature token with trample.
func NewGolemTrampleToken() *Token {
	tok := NewToken("GolemTrampleToken", "3/3 colorless Golem artifact creature token with trample")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewGolemVigilanceToken
// ========================================

// NewGolemVigilanceToken creates a 3/3 colorless Golem artifact creature token with vigilance.
func NewGolemVigilanceToken() *Token {
	tok := NewToken("GolemVigilanceToken", "3/3 colorless Golem artifact creature token with vigilance")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewGolemWhiteBlueToken
// ========================================

// NewGolemWhiteBlueToken creates a 4/4 white and blue Golem artifact creature token.
func NewGolemWhiteBlueToken() *Token {
	tok := NewToken("GolemWhiteBlueToken", "4/4 white and blue Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetColor(Color{White: true, Blue: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewGolemXXToken
// ========================================

// NewGolemXXToken creates a X/X colorless Golem artifact creature token.
func NewGolemXXToken() *Token {
	tok := NewToken("GolemXXToken", "X/X colorless Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	return tok
}

// ========================================
// NewGrakmawSkyclaveRavagerHydraToken
// ========================================

// NewGrakmawSkyclaveRavagerHydraToken creates a X/X black and green Hydra creature token.
func NewGrakmawSkyclaveRavagerHydraToken() *Token {
	tok := NewToken("GrakmawSkyclaveRavagerHydraToken", "X/X black and green Hydra creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HYDRA")
	tok.SetColor(Color{Black: true, Green: true})
	return tok
}

// ========================================
// NewGravebornToken
// ========================================

// NewGravebornToken creates a 3/1 black and red Graveborn creature token with haste.
func NewGravebornToken() *Token {
	tok := NewToken("GravebornToken", "3/1 black and red Graveborn creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GRAVEBORN")
	tok.SetColor(Color{Black: true, Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewGreenAndWhiteElementalToken
// ========================================

// NewGreenAndWhiteElementalToken creates a 8/8 green and white Elemental creature token with vigilance.
func NewGreenAndWhiteElementalToken() *Token {
	tok := NewToken("GreenAndWhiteElementalToken", "8/8 green and white Elemental creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(8, 8)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewGreenCat2Token
// ========================================

// NewGreenCat2Token creates a 2/2 green Cat creature token.
func NewGreenCat2Token() *Token {
	tok := NewToken("GreenCat2Token", "2/2 green Cat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewGreenCatToken
// ========================================

// NewGreenCatToken creates a 1/1 green Cat creature token.
func NewGreenCatToken() *Token {
	tok := NewToken("GreenCatToken", "1/1 green Cat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGreenDogToken
// ========================================

// NewGreenDogToken creates a 1/1 green Dog creature token.
func NewGreenDogToken() *Token {
	tok := NewToken("GreenDogToken", "1/1 green Dog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGreenWhiteElfWarriorToken
// ========================================

// NewGreenWhiteElfWarriorToken creates a 1/1 green and white Elf Warrior creature token.
func NewGreenWhiteElfWarriorToken() *Token {
	tok := NewToken("GreenWhiteElfWarriorToken", "1/1 green and white Elf Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELF")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGremlin11Token
// ========================================

// NewGremlin11Token creates a 1/1 red Gremlin creature token.
func NewGremlin11Token() *Token {
	tok := NewToken("Gremlin11Token", "1/1 red Gremlin creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GREMLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewGremlinArtifactToken
// ========================================

// NewGremlinArtifactToken creates a 0/0 red Gremlin artifact creature token.
func NewGremlinArtifactToken() *Token {
	tok := NewToken("GremlinArtifactToken", "0/0 red Gremlin artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GREMLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewGremlinToken
// ========================================

// NewGremlinToken creates a 2/2 red Gremlin creature token.
func NewGremlinToken() *Token {
	tok := NewToken("GremlinToken", "2/2 red Gremlin creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GREMLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewGriffinToken
// ========================================

// NewGriffinToken creates a 2/2 white Griffin creature token with flying.
func NewGriffinToken() *Token {
	tok := NewToken("GriffinToken", "2/2 white Griffin creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GRIFFIN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewGuenhwyvarToken
// ========================================

// NewGuenhwyvarToken creates a Guenhwyvar, a legendary 4/1 green Cat creature token with trample.
func NewGuenhwyvarToken() *Token {
	tok := NewToken("GuenhwyvarToken", "Guenhwyvar, a legendary 4/1 green Cat creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 1)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewGutterGrimeToken
// ========================================

// NewGutterGrimeToken creates a token.
func NewGutterGrimeToken() *Token {
	tok := NewToken("GutterGrimeToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewGwaihirBirdToken
// ========================================

// NewGwaihirBirdToken creates a token.
func NewGwaihirBirdToken() *Token {
	tok := NewToken("GwaihirBirdToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewHalflingToken
// ========================================

// NewHalflingToken creates a 1/1 white Halfling creature token.
func NewHalflingToken() *Token {
	tok := NewToken("HalflingToken", "1/1 white Halfling creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HALFLING")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHammerOfPurphorosGolemToken
// ========================================

// NewHammerOfPurphorosGolemToken creates a 3/3 colorless Golem enchantment artifact creature token.
func NewHammerOfPurphorosGolemToken() *Token {
	tok := NewToken("HammerOfPurphorosGolemToken", "3/3 colorless Golem enchantment artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewHamsterToken
// ========================================

// NewHamsterToken creates a 1/1 red Hamster creature token.
func NewHamsterToken() *Token {
	tok := NewToken("HamsterToken", "1/1 red Hamster creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HAMSTER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHarpyToken
// ========================================

// NewHarpyToken creates a 1/1 black Harpy creature token with flying.
func NewHarpyToken() *Token {
	tok := NewToken("HarpyToken", "1/1 black Harpy creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HARPY")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewHasteGolemToken
// ========================================

// NewHasteGolemToken creates a X/X colorless Golem artifact creature token with haste.
func NewHasteGolemToken() *Token {
	tok := NewToken("HasteGolemToken", "X/X colorless Golem artifact creature token with haste")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewHauntedAngelToken
// ========================================

// NewHauntedAngelToken creates a 3/3 black Angel creature token with flying.
func NewHauntedAngelToken() *Token {
	tok := NewToken("HauntedAngelToken", "3/3 black Angel creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ANGEL")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewHeliodGodOfTheSunToken
// ========================================

// NewHeliodGodOfTheSunToken creates a 2/1 white Cleric enchantment creature token.
func NewHeliodGodOfTheSunToken() *Token {
	tok := NewToken("HeliodGodOfTheSunToken", "2/1 white Cleric enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("CLERIC")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewHellionToken
// ========================================

// NewHellionToken creates a 4/4 red Hellion creature token.
func NewHellionToken() *Token {
	tok := NewToken("HellionToken", "4/4 red Hellion creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HELLION")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewHeroToken
// ========================================

// NewHeroToken creates a 1/1 colorless Hero creature token.
func NewHeroToken() *Token {
	tok := NewToken("HeroToken", "1/1 colorless Hero creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HERO")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHippoToken
// ========================================

// NewHippoToken creates a 1/1 green Hippo creature token.
func NewHippoToken() *Token {
	tok := NewToken("HippoToken", "1/1 green Hippo creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HIPPO")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHomunculusToken
// ========================================

// NewHomunculusToken creates a 0/1 blue Homunculus artifact creature token.
func NewHomunculusToken() *Token {
	tok := NewToken("HomunculusToken", "0/1 blue Homunculus artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HOMUNCULUS")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewHornetToken
// ========================================

// NewHornetToken creates a 1/1 colorless Insect artifact creature token with flying and haste.
func NewHornetToken() *Token {
	tok := NewToken("HornetToken", "1/1 colorless Insect artifact creature token with flying and haste")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewHorror2Token
// ========================================

// NewHorror2Token creates a 1/1 black Horror creature token.
func NewHorror2Token() *Token {
	tok := NewToken("Horror2Token", "1/1 black Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHorror3Token
// ========================================

// NewHorror3Token creates a 2/2 black Horror creature token.
func NewHorror3Token() *Token {
	tok := NewToken("Horror3Token", "2/2 black Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewHorrorEnchantmentCreatureToken
// ========================================

// NewHorrorEnchantmentCreatureToken creates a 2/2 black Horror enchantment creature token.
func NewHorrorEnchantmentCreatureToken() *Token {
	tok := NewToken("HorrorEnchantmentCreatureToken", "2/2 black Horror enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewHorrorToken
// ========================================

// NewHorrorToken creates a 4/4 black Horror creature token.
func NewHorrorToken() *Token {
	tok := NewToken("HorrorToken", "4/4 black Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewHorrorXXBlackToken
// ========================================

// NewHorrorXXBlackToken creates a X/X black Horror creature token.
func NewHorrorXXBlackToken() *Token {
	tok := NewToken("HorrorXXBlackToken", "X/X black Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewHourOfNeedSphinxToken
// ========================================

// NewHourOfNeedSphinxToken creates a 4/4 blue Sphinx creature token with flying.
func NewHourOfNeedSphinxToken() *Token {
	tok := NewToken("HourOfNeedSphinxToken", "4/4 blue Sphinx creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPHINX")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewHuman11WithWard2Token
// ========================================

// NewHuman11WithWard2Token creates a 1/1 white Human creature token with ward {2}.
func NewHuman11WithWard2Token() *Token {
	tok := NewToken("Human11WithWard2Token", "1/1 white Human creature token with ward {2}")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanCitizenToken
// ========================================

// NewHumanCitizenToken creates a 1/1 green and white Human Citizen creature token.
func NewHumanCitizenToken() *Token {
	tok := NewToken("HumanCitizenToken", "1/1 green and white Human Citizen creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("CITIZEN")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanClericToken
// ========================================

// NewHumanClericToken creates a 1/1 white and black Human Cleric creature token.
func NewHumanClericToken() *Token {
	tok := NewToken("HumanClericToken", "1/1 white and black Human Cleric creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("CLERIC")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanKnightToken
// ========================================

// NewHumanKnightToken creates a 2/2 red Human Knight creature token with trample and haste.
func NewHumanKnightToken() *Token {
	tok := NewToken("HumanKnightToken", "2/2 red Human Knight creature token with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewHumanMonkToken
// ========================================

// NewHumanMonkToken creates a token.
func NewHumanMonkToken() *Token {
	tok := NewToken("HumanMonkToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("MONK")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanRogueToken
// ========================================

// NewHumanRogueToken creates a 1/1 white Human Rogue creature token.
func NewHumanRogueToken() *Token {
	tok := NewToken("HumanRogueToken", "1/1 white Human Rogue creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanSoldierToken
// ========================================

// NewHumanSoldierToken creates a 1/1 white Human Soldier creature token.
func NewHumanSoldierToken() *Token {
	tok := NewToken("HumanSoldierToken", "1/1 white Human Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanSoldierTrainingToken
// ========================================

// NewHumanSoldierTrainingToken creates a 1/1 green and white Human Soldier creature token with training.
func NewHumanSoldierTrainingToken() *Token {
	tok := NewToken("HumanSoldierTrainingToken", "1/1 green and white Human Soldier creature token with training")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanToken
// ========================================

// NewHumanToken creates a 1/1 white Human creature token.
func NewHumanToken() *Token {
	tok := NewToken("HumanToken", "1/1 white Human creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanWarriorToken
// ========================================

// NewHumanWarriorToken creates a 1/1 white Human Warrior creature token.
func NewHumanWarriorToken() *Token {
	tok := NewToken("HumanWarriorToken", "1/1 white Human Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHumanWizardToken
// ========================================

// NewHumanWizardToken creates a 1/1 blue Human Wizard creature token.
func NewHumanWizardToken() *Token {
	tok := NewToken("HumanWizardToken", "1/1 blue Human Wizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewHungryForMoreVampireToken
// ========================================

// NewHungryForMoreVampireToken creates a 3/1 black and red Vampire creature token with trample, lifelink, and haste.
func NewHungryForMoreVampireToken() *Token {
	tok := NewToken("HungryForMoreVampireToken", "3/1 black and red Vampire creature token with trample, lifelink, and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{Black: true, Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewHuntedCentaurToken
// ========================================

// NewHuntedCentaurToken creates a 3/3 green Centaur creature tokens with protection from black.
func NewHuntedCentaurToken() *Token {
	tok := NewToken("HuntedCentaurToken", "3/3 green Centaur creature tokens with protection from black")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CENTAUR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewHuntedDragonKnightToken
// ========================================

// NewHuntedDragonKnightToken creates a 2/2 white Knight creature tokens with first strike.
func NewHuntedDragonKnightToken() *Token {
	tok := NewToken("HuntedDragonKnightToken", "2/2 white Knight creature tokens with first strike")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("first strike")
	return tok
}

// ========================================
// NewHunterToken
// ========================================

// NewHunterToken creates a 4/4 red Hunter creature token.
func NewHunterToken() *Token {
	tok := NewToken("HunterToken", "4/4 red Hunter creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUNTER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewHydraBroodmasterToken
// ========================================

// NewHydraBroodmasterToken creates a green Hydra creature token.
func NewHydraBroodmasterToken() *Token {
	tok := NewToken("HydraBroodmasterToken", "green Hydra creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HYDRA")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewIcingdeathFrostTongueToken
// ========================================

// NewIcingdeathFrostTongueToken creates a token.
func NewIcingdeathFrostTongueToken() *Token {
	tok := NewToken("IcingdeathFrostTongueToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("EQUIPMENT")
	tok.SetColor(Color{White: true})
	return tok
}

// ========================================
// NewIcyManalithToken
// ========================================

// NewIcyManalithToken creates a token.
func NewIcyManalithToken() *Token {
	tok := NewToken("IcyManalithToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewIllusionToken
// ========================================

// NewIllusionToken creates a 2/2 blue Illusion creature token.
func NewIllusionToken() *Token {
	tok := NewToken("IllusionToken", "2/2 blue Illusion creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewIllusionVillainToken
// ========================================

// NewIllusionVillainToken creates a 3/3 blue Illusion Villain creature token.
func NewIllusionVillainToken() *Token {
	tok := NewToken("IllusionVillainToken", "3/3 blue Illusion Villain creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.AddSubtype("VILLAIN")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewImpToken
// ========================================

// NewImpToken creates a token.
func NewImpToken() *Token {
	tok := NewToken("ImpToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("IMP")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewIncubatorToken
// ========================================

// NewIncubatorToken creates a token.
func NewIncubatorToken() *Token {
	tok := NewToken("IncubatorToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("INCUBATOR")
	return tok
}

// ========================================
// NewInexorableBlobOozeToken
// ========================================

// NewInexorableBlobOozeToken creates a 3/3 green Ooze creature token.
func NewInexorableBlobOozeToken() *Token {
	tok := NewToken("InexorableBlobOozeToken", "3/3 green Ooze creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewInklingToken
// ========================================

// NewInklingToken creates a 2/1 white and black Inkling creature token with flying.
func NewInklingToken() *Token {
	tok := NewToken("InklingToken", "2/1 white and black Inkling creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INKLING")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewInsectBlackGreenFlyingToken
// ========================================

// NewInsectBlackGreenFlyingToken creates a 1/1 black and green Insect creature token with flying.
func NewInsectBlackGreenFlyingToken() *Token {
	tok := NewToken("InsectBlackGreenFlyingToken", "1/1 black and green Insect creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewInsectColorlessArtifactToken
// ========================================

// NewInsectColorlessArtifactToken creates a 1/1 colorless Insect artifact creature token with flying.
func NewInsectColorlessArtifactToken() *Token {
	tok := NewToken("InsectColorlessArtifactToken", "1/1 colorless Insect artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewInsectDeathToken
// ========================================

// NewInsectDeathToken creates a 1/1 green Insect creature token with flying and deathtouch.
func NewInsectDeathToken() *Token {
	tok := NewToken("InsectDeathToken", "1/1 green Insect creature token with flying and deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewInsectInfectToken
// ========================================

// NewInsectInfectToken creates a 1/1 green Phyrexian Insect creature token with infect.
func NewInsectInfectToken() *Token {
	tok := NewToken("InsectInfectToken", "1/1 green Phyrexian Insect creature token with infect")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewInsectToken
// ========================================

// NewInsectToken creates a 1/1 green Insect creature token.
func NewInsectToken() *Token {
	tok := NewToken("InsectToken", "1/1 green Insect creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewInsectWhiteToken
// ========================================

// NewInsectWhiteToken creates a 2/1 white Insect creature token with flying.
func NewInsectWhiteToken() *Token {
	tok := NewToken("InsectWhiteToken", "2/1 white Insect creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewIxalanVampireToken
// ========================================

// NewIxalanVampireToken creates a 1/1 white Vampire creature token with lifelink.
func NewIxalanVampireToken() *Token {
	tok := NewToken("IxalanVampireToken", "1/1 white Vampire creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewIzoniInsectToken
// ========================================

// NewIzoniInsectToken creates a 1/1 black and green Insect creature token.
func NewIzoniInsectToken() *Token {
	tok := NewToken("IzoniInsectToken", "1/1 black and green Insect creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewIzoniSpiderToken
// ========================================

// NewIzoniSpiderToken creates a 2/1 black and green Spider creature token with menace and reach.
func NewIzoniSpiderToken() *Token {
	tok := NewToken("IzoniSpiderToken", "2/1 black and green Spider creature token with menace and reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIDER")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewJaceCunningCastawayIllusionToken
// ========================================

// NewJaceCunningCastawayIllusionToken creates a token.
func NewJaceCunningCastawayIllusionToken() *Token {
	tok := NewToken("JaceCunningCastawayIllusionToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewJellyfishToken
// ========================================

// NewJellyfishToken creates a 1/1 blue Jellyfish creature token with flying.
func NewJellyfishToken() *Token {
	tok := NewToken("JellyfishToken", "1/1 blue Jellyfish creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("JELLYFISH")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewJoinTheRanksSoldierToken
// ========================================

// NewJoinTheRanksSoldierToken creates a 1/1 white Soldier Ally creature token.
func NewJoinTheRanksSoldierToken() *Token {
	tok := NewToken("JoinTheRanksSoldierToken", "1/1 white Soldier Ally creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.AddSubtype("ALLY")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewJumblebonesToken
// ========================================

// NewJumblebonesToken creates a token.
func NewJumblebonesToken() *Token {
	tok := NewToken("JumblebonesToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SKELETON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewJunkToken
// ========================================

// NewJunkToken creates a Junk token.
func NewJunkToken() *Token {
	tok := NewToken("JunkToken", "Junk token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("JUNK")
	return tok
}

// ========================================
// NewKaldraToken
// ========================================

// NewKaldraToken creates a Kaldra, a legendary 4/4 colorless Avatar creature token.
func NewKaldraToken() *Token {
	tok := NewToken("KaldraToken", "Kaldra, a legendary 4/4 colorless Avatar creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("AVATAR")
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewKalitasVampireToken
// ========================================

// NewKalitasVampireToken creates a token.
func NewKalitasVampireToken() *Token {
	tok := NewToken("KalitasVampireToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewKalonianTwingroveTreefolkWarriorToken
// ========================================

// NewKalonianTwingroveTreefolkWarriorToken creates a token.
func NewKalonianTwingroveTreefolkWarriorToken() *Token {
	tok := NewToken("KalonianTwingroveTreefolkWarriorToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TREEFOLK")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewKarnConstructToken
// ========================================

// NewKarnConstructToken creates a token.
func NewKarnConstructToken() *Token {
	tok := NewToken("KarnConstructToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewKaroxBladewingDragonToken
// ========================================

// NewKaroxBladewingDragonToken creates a Karox Bladewing, a legendary 4/4 red Dragon creature token with flying.
func NewKaroxBladewingDragonToken() *Token {
	tok := NewToken("KaroxBladewingDragonToken", "Karox Bladewing, a legendary 4/4 red Dragon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewKavuAllColorToken
// ========================================

// NewKavuAllColorToken creates a 3/3 Kavu creature token with trample that's all colors.
func NewKavuAllColorToken() *Token {
	tok := NewToken("KavuAllColorToken", "3/3 Kavu creature token with trample that's all colors")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KAVU")
	tok.SetColor(Color{White: true, Blue: true, Black: true, Red: true, Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewKeimiToken
// ========================================

// NewKeimiToken creates a token.
func NewKeimiToken() *Token {
	tok := NewToken("KeimiToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FROG")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewKelpToken
// ========================================

// NewKelpToken creates a 0/1 blue Plant Wall creature token with defender named Kelp.
func NewKelpToken() *Token {
	tok := NewToken("KelpToken", "0/1 blue Plant Wall creature token with defender named Kelp")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PLANT")
	tok.AddSubtype("WALL")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 1)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewKherKeepKoboldToken
// ========================================

// NewKherKeepKoboldToken creates a 0/1 red Kobold creature token named Kobolds of Kher Keep.
func NewKherKeepKoboldToken() *Token {
	tok := NewToken("KherKeepKoboldToken", "0/1 red Kobold creature token named Kobolds of Kher Keep")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KOBOLD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewKithkinSoldierToken
// ========================================

// NewKithkinSoldierToken creates a 1/1 white Kithkin Soldier creature token.
func NewKithkinSoldierToken() *Token {
	tok := NewToken("KithkinSoldierToken", "1/1 white Kithkin Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KITHKIN")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewKnight31RedToken
// ========================================

// NewKnight31RedToken creates a 3/1 red Knight creature token.
func NewKnight31RedToken() *Token {
	tok := NewToken("Knight31RedToken", "3/1 red Knight creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	return tok
}

// ========================================
// NewKnight33Token
// ========================================

// NewKnight33Token creates a 3/3 white Knight creature token.
func NewKnight33Token() *Token {
	tok := NewToken("Knight33Token", "3/3 white Knight creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewKnightAllyToken
// ========================================

// NewKnightAllyToken creates a 2/2 white Knight Ally creature token.
func NewKnightAllyToken() *Token {
	tok := NewToken("KnightAllyToken", "2/2 white Knight Ally creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.AddSubtype("ALLY")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewKnightToken
// ========================================

// NewKnightToken creates a 2/2 white Knight creature token with vigilance.
func NewKnightToken() *Token {
	tok := NewToken("KnightToken", "2/2 white Knight creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewKnightWhiteBlueToken
// ========================================

// NewKnightWhiteBlueToken creates a 2/2 white and blue Knight creature token with vigilance.
func NewKnightWhiteBlueToken() *Token {
	tok := NewToken("KnightWhiteBlueToken", "2/2 white and blue Knight creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true, Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewKomasCoilToken
// ========================================

// NewKomasCoilToken creates a 3/3 blue Serpent creature token named Koma's Coil.
func NewKomasCoilToken() *Token {
	tok := NewToken("KomasCoilToken", "3/3 blue Serpent creature token named Koma's Coil")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SERPENT")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewKorAllyToken
// ========================================

// NewKorAllyToken creates a 1/1 white Kor Ally creature token.
func NewKorAllyToken() *Token {
	tok := NewToken("KorAllyToken", "1/1 white Kor Ally creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KOR")
	tok.AddSubtype("ALLY")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewKorSoldierToken
// ========================================

// NewKorSoldierToken creates a 1/1 white Kor Soldier creature token.
func NewKorSoldierToken() *Token {
	tok := NewToken("KorSoldierToken", "1/1 white Kor Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KOR")
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewKorWarriorToken
// ========================================

// NewKorWarriorToken creates a 1/1 white Kor Warrior creature token.
func NewKorWarriorToken() *Token {
	tok := NewToken("KorWarriorToken", "1/1 white Kor Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KOR")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewKraken11Token
// ========================================

// NewKraken11Token creates a 1/1 blue Kraken creature token with trample.
func NewKraken11Token() *Token {
	tok := NewToken("Kraken11Token", "1/1 blue Kraken creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KRAKEN")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewKraken99Token
// ========================================

// NewKraken99Token creates a 9/9 blue Kraken creature token.
func NewKraken99Token() *Token {
	tok := NewToken("Kraken99Token", "9/9 blue Kraken creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KRAKEN")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(9, 9)
	return tok
}

// ========================================
// NewKrakenHexproofToken
// ========================================

// NewKrakenHexproofToken creates a 8/8 blue Kraken creature token with hexproof.
func NewKrakenHexproofToken() *Token {
	tok := NewToken("KrakenHexproofToken", "8/8 blue Kraken creature token with hexproof")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KRAKEN")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(8, 8)
	tok.AddAbility("hexproof")
	return tok
}

// ========================================
// NewKrakenToken
// ========================================

// NewKrakenToken creates a 8/8 blue Kraken creature token.
func NewKrakenToken() *Token {
	tok := NewToken("KrakenToken", "8/8 blue Kraken creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KRAKEN")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(8, 8)
	return tok
}

// ========================================
// NewLandMineToken
// ========================================

// NewLandMineToken creates a token.
func NewLandMineToken() *Token {
	tok := NewToken("LandMineToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewLanderToken
// ========================================

// NewLanderToken creates a Lander token.
func NewLanderToken() *Token {
	tok := NewToken("LanderToken", "Lander token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("LANDER")
	return tok
}

// ========================================
// NewLeafdrakeRoostDrakeToken
// ========================================

// NewLeafdrakeRoostDrakeToken creates a 2/2 green and blue Drake creature token with flying.
func NewLeafdrakeRoostDrakeToken() *Token {
	tok := NewToken("LeafdrakeRoostDrakeToken", "2/2 green and blue Drake creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAKE")
	tok.SetColor(Color{Blue: true, Green: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewLightningRagerToken
// ========================================

// NewLightningRagerToken creates a token.
func NewLightningRagerToken() *Token {
	tok := NewToken("LightningRagerToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(5, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewLizardToken
// ========================================

// NewLizardToken creates a 2/2 green Lizard creature token.
func NewLizardToken() *Token {
	tok := NewToken("LizardToken", "2/2 green Lizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("LIZARD")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewLlanowarElvesToken
// ========================================

// NewLlanowarElvesToken creates a token.
func NewLlanowarElvesToken() *Token {
	tok := NewToken("LlanowarElvesToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELF")
	tok.AddSubtype("DRUID")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewLolthSpiderToken
// ========================================

// NewLolthSpiderToken creates a 2/1 black Spider creature token with menace and reach.
func NewLolthSpiderToken() *Token {
	tok := NewToken("LolthSpiderToken", "2/1 black Spider creature token with menace and reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIDER")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewMagesAttendantToken
// ========================================

// NewMagesAttendantToken creates a token.
func NewMagesAttendantToken() *Token {
	tok := NewToken("MagesAttendantToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMapToken
// ========================================

// NewMapToken creates a Map token.
func NewMapToken() *Token {
	tok := NewToken("MapToken", "Map token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("MAP")
	return tok
}

// ========================================
// NewMarduStrikeLeaderWarriorToken
// ========================================

// NewMarduStrikeLeaderWarriorToken creates a 2/1 black Warrior creature token.
func NewMarduStrikeLeaderWarriorToken() *Token {
	tok := NewToken("MarduStrikeLeaderWarriorToken", "2/1 black Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewMaritLageToken
// ========================================

// NewMaritLageToken creates a Marit Lage, a legendary 20/20 black Avatar creature token with flying and indestructible.
func NewMaritLageToken() *Token {
	tok := NewToken("MaritLageToken", "Marit Lage, a legendary 20/20 black Avatar creature token with flying and indestructible")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("AVATAR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(20, 20)
	tok.AddAbility("flying")
	tok.AddAbility("indestructible")
	return tok
}

// ========================================
// NewMarkOfTheRaniToken
// ========================================

// NewMarkOfTheRaniToken creates a token.
func NewMarkOfTheRaniToken() *Token {
	tok := NewToken("MarkOfTheRaniToken", "")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.SetColor(Color{Red: true})
	return tok
}

// ========================================
// NewMaskToken
// ========================================

// NewMaskToken creates a token.
func NewMaskToken() *Token {
	tok := NewToken("MaskToken", "")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.SetColor(Color{White: true})
	return tok
}

// ========================================
// NewMasterOfWavesElementalToken
// ========================================

// NewMasterOfWavesElementalToken creates a 1/0 blue Elemental creature.
func NewMasterOfWavesElementalToken() *Token {
	tok := NewToken("MasterOfWavesElementalToken", "1/0 blue Elemental creature")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 0)
	return tok
}

// ========================================
// NewMechtitanToken
// ========================================

// NewMechtitanToken creates a Mechtitan, a legendary 10/10 Construct artifact creature token with flying, vigilance, trample, lifelink, and haste that's all colors.
func NewMechtitanToken() *Token {
	tok := NewToken("MechtitanToken", "Mechtitan, a legendary 10/10 Construct artifact creature token with flying, vigilance, trample, lifelink, and haste that's all colors")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetColor(Color{White: true, Blue: true, Black: true, Red: true, Green: true})
	tok.SetPowerToughness(10, 10)
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	tok.AddAbility("vigilance")
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewMelokuTheCloudedMirrorToken
// ========================================

// NewMelokuTheCloudedMirrorToken creates a 1/1 blue Illusion creature token with flying.
func NewMelokuTheCloudedMirrorToken() *Token {
	tok := NewToken("MelokuTheCloudedMirrorToken", "1/1 blue Illusion creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewMercenaryToken
// ========================================

// NewMercenaryToken creates a token.
func NewMercenaryToken() *Token {
	tok := NewToken("MercenaryToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MERCENARY")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMerfolkHexproofToken
// ========================================

// NewMerfolkHexproofToken creates a 1/1 blue Merfolk creature token with hexproof.
func NewMerfolkHexproofToken() *Token {
	tok := NewToken("MerfolkHexproofToken", "1/1 blue Merfolk creature token with hexproof")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MERFOLK")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("hexproof")
	return tok
}

// ========================================
// NewMerfolkToken
// ========================================

// NewMerfolkToken creates a 1/1 blue Merfolk creature token.
func NewMerfolkToken() *Token {
	tok := NewToken("MerfolkToken", "1/1 blue Merfolk creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MERFOLK")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMerfolkWizardToken
// ========================================

// NewMerfolkWizardToken creates a 1/1 blue Merfolk Wizard creature token.
func NewMerfolkWizardToken() *Token {
	tok := NewToken("MerfolkWizardToken", "1/1 blue Merfolk Wizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MERFOLK")
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMesmerizingBenthidToken
// ========================================

// NewMesmerizingBenthidToken creates a token.
func NewMesmerizingBenthidToken() *Token {
	tok := NewToken("MesmerizingBenthidToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 2)
	return tok
}

// ========================================
// NewMetallicSliverToken
// ========================================

// NewMetallicSliverToken creates a 1/1 colorless Sliver artifact creature token named Metallic Sliver.
func NewMetallicSliverToken() *Token {
	tok := NewToken("MetallicSliverToken", "1/1 colorless Sliver artifact creature token named Metallic Sliver")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SLIVER")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMetallurgicSummoningsConstructToken
// ========================================

// NewMetallurgicSummoningsConstructToken creates a X/X colorless Construct artifact creature token.
func NewMetallurgicSummoningsConstructToken() *Token {
	tok := NewToken("MetallurgicSummoningsConstructToken", "X/X colorless Construct artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	return tok
}

// ========================================
// NewMeteoriteToken
// ========================================

// NewMeteoriteToken creates a token.
func NewMeteoriteToken() *Token {
	tok := NewToken("MeteoriteToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewMinionToken
// ========================================

// NewMinionToken creates a 1/1 black Minion creature token.
func NewMinionToken() *Token {
	tok := NewToken("MinionToken", "1/1 black Minion creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MINION")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMinnWilyIllusionistToken
// ========================================

// NewMinnWilyIllusionistToken creates a token.
func NewMinnWilyIllusionistToken() *Token {
	tok := NewToken("MinnWilyIllusionistToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ILLUSION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMinorDemonToken
// ========================================

// NewMinorDemonToken creates a 1/1 black and red Demon creature token named Minor Demon.
func NewMinorDemonToken() *Token {
	tok := NewToken("MinorDemonToken", "1/1 black and red Demon creature token named Minor Demon")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true, Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMinotaurToken
// ========================================

// NewMinotaurToken creates a 2/3 red Minotaur creature token.
func NewMinotaurToken() *Token {
	tok := NewToken("MinotaurToken", "2/3 red Minotaur creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MINOTAUR")
	tok.SetPowerToughness(2, 3)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewMitoticSlimeOozeToken
// ========================================

// NewMitoticSlimeOozeToken creates a token.
func NewMitoticSlimeOozeToken() *Token {
	tok := NewToken("MitoticSlimeOozeToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewMonasteryMentorToken
// ========================================

// NewMonasteryMentorToken creates a 1/1 white Monk creature token with prowess.
func NewMonasteryMentorToken() *Token {
	tok := NewToken("MonasteryMentorToken", "1/1 white Monk creature token with prowess")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MONK")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMonkRedToken
// ========================================

// NewMonkRedToken creates a 1/1 red Monk creature token with prowess.
func NewMonkRedToken() *Token {
	tok := NewToken("MonkRedToken", "1/1 red Monk creature token with prowess")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MONK")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMonkeyToken
// ========================================

// NewMonkeyToken creates a 2/2 green Monkey creature token.
func NewMonkeyToken() *Token {
	tok := NewToken("MonkeyToken", "2/2 green Monkey creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MONKEY")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewMonsterRoleToken
// ========================================

// NewMonsterRoleToken creates a Monster Role token.
func NewMonsterRoleToken() *Token {
	tok := NewToken("MonsterRoleToken", "Monster Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewMoogleToken
// ========================================

// NewMoogleToken creates a 1/2 white Moogle creature token with lifelink.
func NewMoogleToken() *Token {
	tok := NewToken("MoogleToken", "1/2 white Moogle creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MOOGLE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewMoonfolk12FlyingToken
// ========================================

// NewMoonfolk12FlyingToken creates a 1/2 blue Moonfolk creature token with flying.
func NewMoonfolk12FlyingToken() *Token {
	tok := NewToken("Moonfolk12FlyingToken", "1/2 blue Moonfolk creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MOONFOLK")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewMouseToken
// ========================================

// NewMouseToken creates a 1/1 white Mouse creature token.
func NewMouseToken() *Token {
	tok := NewToken("MouseToken", "1/1 white Mouse creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MOUSE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewMowuToken
// ========================================

// NewMowuToken creates a Mowu, a legendary 3/3 green Dog creature token.
func NewMowuToken() *Token {
	tok := NewToken("MowuToken", "Mowu, a legendary 3/3 green Dog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewMuYanlingSkyDancerToken
// ========================================

// NewMuYanlingSkyDancerToken creates a 4/4 blue Elemental Bird creature token with flying.
func NewMuYanlingSkyDancerToken() *Token {
	tok := NewToken("MuYanlingSkyDancerToken", "4/4 blue Elemental Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewMunitionsToken
// ========================================

// NewMunitionsToken creates a token.
func NewMunitionsToken() *Token {
	tok := NewToken("MunitionsToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewMutagenToken
// ========================================

// NewMutagenToken creates a Mutagen token.
func NewMutagenToken() *Token {
	tok := NewToken("MutagenToken", "Mutagen token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("MUTAGEN")
	return tok
}

// ========================================
// NewMutant33DeathtouchToken
// ========================================

// NewMutant33DeathtouchToken creates a 3/3 green mutant creature token with deathtouch.
func NewMutant33DeathtouchToken() *Token {
	tok := NewToken("Mutant33DeathtouchToken", "3/3 green mutant creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MUTANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewMutavaultToken
// ========================================

// NewMutavaultToken creates a Mutavault token.
func NewMutavaultToken() *Token {
	tok := NewToken("MutavaultToken", "Mutavault token")
	tok.AddCardType(CardTypeLand)
	return tok
}

// ========================================
// NewMyrToken
// ========================================

// NewMyrToken creates a 1/1 colorless Myr artifact creature token.
func NewMyrToken() *Token {
	tok := NewToken("MyrToken", "1/1 colorless Myr artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MYR")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewNahiriTheLithomancerEquipmentToken
// ========================================

// NewNahiriTheLithomancerEquipmentToken creates a token.
func NewNahiriTheLithomancerEquipmentToken() *Token {
	tok := NewToken("NahiriTheLithomancerEquipmentToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("EQUIPMENT")
	tok.AddAbility("double strike")
	tok.AddAbility("indestructible")
	return tok
}

// ========================================
// NewNalaarAetherjetToken
// ========================================

// NewNalaarAetherjetToken creates a X/X colorless Vehicle artifact token named Nalaar Aetherjet with flying and crew 2.
func NewNalaarAetherjetToken() *Token {
	tok := NewToken("NalaarAetherjetToken", "X/X colorless Vehicle artifact token named Nalaar Aetherjet with flying and crew 2")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("VEHICLE")
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewNecronWarriorToken
// ========================================

// NewNecronWarriorToken creates a 2/2 black Necron Warrior artifact creature token.
func NewNecronWarriorToken() *Token {
	tok := NewToken("NecronWarriorToken", "2/2 black Necron Warrior artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("NECRON")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewNestOfScarabsBlackInsectToken
// ========================================

// NewNestOfScarabsBlackInsectToken creates a 1/1 black Insect creature token.
func NewNestOfScarabsBlackInsectToken() *Token {
	tok := NewToken("NestOfScarabsBlackInsectToken", "1/1 black Insect creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewNestingDragonToken
// ========================================

// NewNestingDragonToken creates a token.
func NewNestingDragonToken() *Token {
	tok := NewToken("NestingDragonToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.AddSubtype("EGG")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(0, 2)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewNettlingNuisancePirateToken
// ========================================

// NewNettlingNuisancePirateToken creates a token.
func NewNettlingNuisancePirateToken() *Token {
	tok := NewToken("NettlingNuisancePirateToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PIRATE")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 2)
	return tok
}

// ========================================
// NewNightwingHorrorToken
// ========================================

// NewNightwingHorrorToken creates a 1/1 blue and black Horror creature token with flying.
func NewNightwingHorrorToken() *Token {
	tok := NewToken("NightwingHorrorToken", "1/1 blue and black Horror creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Blue: true, Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewNinjaToken
// ========================================

// NewNinjaToken creates a token.
func NewNinjaToken() *Token {
	tok := NewToken("NinjaToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("NINJA")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewNissaSageAnimistToken
// ========================================

// NewNissaSageAnimistToken creates a token.
func NewNissaSageAnimistToken() *Token {
	tok := NewToken("NissaSageAnimistToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewNoFlyingSpiritWhiteToken
// ========================================

// NewNoFlyingSpiritWhiteToken creates a 1/1 white Spirit creature token.
func NewNoFlyingSpiritWhiteToken() *Token {
	tok := NewToken("NoFlyingSpiritWhiteToken", "1/1 white Spirit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewNymphToken
// ========================================

// NewNymphToken creates a 2/2 white Nymph enchantment creature token.
func NewNymphToken() *Token {
	tok := NewToken("NymphToken", "2/2 white Nymph enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("NYMPH")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewOctopusToken
// ========================================

// NewOctopusToken creates a 8/8 blue Octopus creature token.
func NewOctopusToken() *Token {
	tok := NewToken("OctopusToken", "8/8 blue Octopus creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OCTOPUS")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(8, 8)
	return tok
}

// ========================================
// NewOgreToken
// ========================================

// NewOgreToken creates a 3/3 red Ogre creature token.
func NewOgreToken() *Token {
	tok := NewToken("OgreToken", "3/3 red Ogre creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OGRE")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewOgreWarriorToken
// ========================================

// NewOgreWarriorToken creates a 4/3 black Ogre Warrior creature token.
func NewOgreWarriorToken() *Token {
	tok := NewToken("OgreWarriorToken", "4/3 black Ogre Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OGRE")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(4, 3)
	return tok
}

// ========================================
// NewOminousRoostBirdToken
// ========================================

// NewOminousRoostBirdToken creates a token.
func NewOminousRoostBirdToken() *Token {
	tok := NewToken("OminousRoostBirdToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewOmnathElementalToken
// ========================================

// NewOmnathElementalToken creates a 5/5 red and green Elemental creature token.
func NewOmnathElementalToken() *Token {
	tok := NewToken("OmnathElementalToken", "5/5 red and green Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewOneDozenEyesBeastToken
// ========================================

// NewOneDozenEyesBeastToken creates a 5/5 green Beast creature token.
func NewOneDozenEyesBeastToken() *Token {
	tok := NewToken("OneDozenEyesBeastToken", "5/5 green Beast creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewOonaQueenFaerieRogueToken
// ========================================

// NewOonaQueenFaerieRogueToken creates a 1/1 blue and black Faerie Rogue creature token with flying.
func NewOonaQueenFaerieRogueToken() *Token {
	tok := NewToken("OonaQueenFaerieRogueToken", "1/1 blue and black Faerie Rogue creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FAERIE")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Blue: true, Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewOozeToken
// ========================================

// NewOozeToken creates a X/X green Ooze creature token.
func NewOozeToken() *Token {
	tok := NewToken("OozeToken", "X/X green Ooze creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewOozeTrampleToken
// ========================================

// NewOozeTrampleToken creates a 0/0 green Ooze creature token with trample.
func NewOozeTrampleToken() *Token {
	tok := NewToken("OozeTrampleToken", "0/0 green Ooze creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OOZE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewOphiomancerSnakeToken
// ========================================

// NewOphiomancerSnakeToken creates a 1/1 black Snake creature token with deathtouch.
func NewOphiomancerSnakeToken() *Token {
	tok := NewToken("OphiomancerSnakeToken", "1/1 black Snake creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAKE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewOrcArmyToken
// ========================================

// NewOrcArmyToken creates a 0/0 black Orc Army creature token.
func NewOrcArmyToken() *Token {
	tok := NewToken("OrcArmyToken", "0/0 black Orc Army creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ORC")
	tok.AddSubtype("ARMY")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewOrnithopterToken
// ========================================

// NewOrnithopterToken creates a 0/2 colorless Thopter artifact creature token with flying named Ornithopter.
func NewOrnithopterToken() *Token {
	tok := NewToken("OrnithopterToken", "0/2 colorless Thopter artifact creature token with flying named Ornithopter")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("THOPTER")
	tok.SetPowerToughness(0, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewOtterProwessToken
// ========================================

// NewOtterProwessToken creates a 1/1 blue and red Otter creature token with prowess.
func NewOtterProwessToken() *Token {
	tok := NewToken("OtterProwessToken", "1/1 blue and red Otter creature token with prowess")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OTTER")
	tok.SetColor(Color{Blue: true, Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewOutlawsMerrimentClericToken
// ========================================

// NewOutlawsMerrimentClericToken creates a 2/1 Human Cleric with lifelink and haste.
func NewOutlawsMerrimentClericToken() *Token {
	tok := NewToken("OutlawsMerrimentClericToken", "2/1 Human Cleric with lifelink and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("CLERIC")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("haste")
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewOutlawsMerrimentRogueToken
// ========================================

// NewOutlawsMerrimentRogueToken creates a token.
func NewOutlawsMerrimentRogueToken() *Token {
	tok := NewToken("OutlawsMerrimentRogueToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewOutlawsMerrimentWarriorToken
// ========================================

// NewOutlawsMerrimentWarriorToken creates a 3/1 Human Warrior with trample and haste.
func NewOutlawsMerrimentWarriorToken() *Token {
	tok := NewToken("OutlawsMerrimentWarriorToken", "3/1 Human Warrior with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewOviyaPashiriSageLifecrafterToken
// ========================================

// NewOviyaPashiriSageLifecrafterToken creates a X/X colorless Construct artifact creature token, where X is the number of creatures you control.
func NewOviyaPashiriSageLifecrafterToken() *Token {
	tok := NewToken("OviyaPashiriSageLifecrafterToken", "X/X colorless Construct artifact creature token, where X is the number of creatures you control")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	return tok
}

// ========================================
// NewOx22Token
// ========================================

// NewOx22Token creates a 2/2 white Ox creature token.
func NewOx22Token() *Token {
	tok := NewToken("Ox22Token", "2/2 white Ox creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OX")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewOxGreenToken
// ========================================

// NewOxGreenToken creates a 4/4 green Ox creature token.
func NewOxGreenToken() *Token {
	tok := NewToken("OxGreenToken", "4/4 green Ox creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OX")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewOxToken
// ========================================

// NewOxToken creates a 2/4 white Ox creature token.
func NewOxToken() *Token {
	tok := NewToken("OxToken", "2/4 white Ox creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OX")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 4)
	return tok
}

// ========================================
// NewPatagiaViperSnakeToken
// ========================================

// NewPatagiaViperSnakeToken creates a 1/1 green and blue Snake creature token.
func NewPatagiaViperSnakeToken() *Token {
	tok := NewToken("PatagiaViperSnakeToken", "1/1 green and blue Snake creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAKE")
	tok.SetColor(Color{Blue: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPegasusToken
// ========================================

// NewPegasusToken creates a 1/1 white Pegasus creature token with flying.
func NewPegasusToken() *Token {
	tok := NewToken("PegasusToken", "1/1 white Pegasus creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PEGASUS")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewPentaviteToken
// ========================================

// NewPentaviteToken creates a 1/1 colorless Pentavite artifact creature token with flying.
func NewPentaviteToken() *Token {
	tok := NewToken("PentaviteToken", "1/1 colorless Pentavite artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PENTAVITE")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewPenumbraBobcatToken
// ========================================

// NewPenumbraBobcatToken creates a 2/1 black Cat creature token.
func NewPenumbraBobcatToken() *Token {
	tok := NewToken("PenumbraBobcatToken", "2/1 black Cat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewPenumbraKavuToken
// ========================================

// NewPenumbraKavuToken creates a 3/3 black Kavu creature token.
func NewPenumbraKavuToken() *Token {
	tok := NewToken("PenumbraKavuToken", "3/3 black Kavu creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KAVU")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewPenumbraSpiderToken
// ========================================

// NewPenumbraSpiderToken creates a 2/4 black Spider creature token with reach.
func NewPenumbraSpiderToken() *Token {
	tok := NewToken("PenumbraSpiderToken", "2/4 black Spider creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIDER")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 4)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewPenumbraWurmToken
// ========================================

// NewPenumbraWurmToken creates a 6/6 black Wurm creature token with trample.
func NewPenumbraWurmToken() *Token {
	tok := NewToken("PenumbraWurmToken", "6/6 black Wurm creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(6, 6)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewPest11GainLifeToken
// ========================================

// NewPest11GainLifeToken creates a token.
func NewPest11GainLifeToken() *Token {
	tok := NewToken("Pest11GainLifeToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PEST")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPestToken
// ========================================

// NewPestToken creates a 0/1 colorless Pest artifact creature token.
func NewPestToken() *Token {
	tok := NewToken("PestToken", "0/1 colorless Pest artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PEST")
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewPharikaSnakeToken
// ========================================

// NewPharikaSnakeToken creates a 1/1 black and green Snake enchantment creature token with deathtouch.
func NewPharikaSnakeToken() *Token {
	tok := NewToken("PharikaSnakeToken", "1/1 black and green Snake enchantment creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("SNAKE")
	tok.SetColor(Color{Black: true, Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewPhobosToken
// ========================================

// NewPhobosToken creates a Phobos, a legendary 3/2 red Horse creature token.
func NewPhobosToken() *Token {
	tok := NewToken("PhobosToken", "Phobos, a legendary 3/2 red Horse creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORSE")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewPhyrexian00Token
// ========================================

// NewPhyrexian00Token creates a 0/0 Phyrexian artifact creature token.
func NewPhyrexian00Token() *Token {
	tok := NewToken("Phyrexian00Token", "0/0 Phyrexian artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewPhyrexianBeastToken
// ========================================

// NewPhyrexianBeastToken creates a 4/4 green Phyrexian Beast creature token.
func NewPhyrexianBeastToken() *Token {
	tok := NewToken("PhyrexianBeastToken", "4/4 green Phyrexian Beast creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewPhyrexianBeastToxicToken
// ========================================

// NewPhyrexianBeastToxicToken creates a 3/3 green Phyrexian Beast creature token with toxic 1.
func NewPhyrexianBeastToxicToken() *Token {
	tok := NewToken("PhyrexianBeastToxicToken", "3/3 green Phyrexian Beast creature token with toxic 1")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewPhyrexianGermToken
// ========================================

// NewPhyrexianGermToken creates a 0/0 black Phyrexian Germ creature token.
func NewPhyrexianGermToken() *Token {
	tok := NewToken("PhyrexianGermToken", "0/0 black Phyrexian Germ creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("GERM")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewPhyrexianGoblinHasteToken
// ========================================

// NewPhyrexianGoblinHasteToken creates a 1/1 red Phyrexian Goblin creature token with haste.
func NewPhyrexianGoblinHasteToken() *Token {
	tok := NewToken("PhyrexianGoblinHasteToken", "1/1 red Phyrexian Goblin creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("GOBLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewPhyrexianGoblinToken
// ========================================

// NewPhyrexianGoblinToken creates a 1/1 red Phyrexian Goblin creature token.
func NewPhyrexianGoblinToken() *Token {
	tok := NewToken("PhyrexianGoblinToken", "1/1 red Phyrexian Goblin creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("GOBLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPhyrexianGolemToken
// ========================================

// NewPhyrexianGolemToken creates a 3/3 colorless Phyrexian Golem artifact creature token.
func NewPhyrexianGolemToken() *Token {
	tok := NewToken("PhyrexianGolemToken", "3/3 colorless Phyrexian Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewPhyrexianHorrorGreenToken
// ========================================

// NewPhyrexianHorrorGreenToken creates a X/X green Phyrexian Horror creature token.
func NewPhyrexianHorrorGreenToken() *Token {
	tok := NewToken("PhyrexianHorrorGreenToken", "X/X green Phyrexian Horror creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewPhyrexianHorrorRedToken
// ========================================

// NewPhyrexianHorrorRedToken creates a X/1 red Phyrexian Horror creature token with trample and haste.
func NewPhyrexianHorrorRedToken() *Token {
	tok := NewToken("PhyrexianHorrorRedToken", "X/1 red Phyrexian Horror creature token with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Red: true})
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewPhyrexianHydraWithLifelinkToken
// ========================================

// NewPhyrexianHydraWithLifelinkToken creates a 3/3 green and white Phyrexian Hydra creature token with lifelink.
func NewPhyrexianHydraWithLifelinkToken() *Token {
	tok := NewToken("PhyrexianHydraWithLifelinkToken", "3/3 green and white Phyrexian Hydra creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("HYDRA")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewPhyrexianHydraWithReachToken
// ========================================

// NewPhyrexianHydraWithReachToken creates a 3/3 green and white Phyrexian Hydra creature token with reach.
func NewPhyrexianHydraWithReachToken() *Token {
	tok := NewToken("PhyrexianHydraWithReachToken", "3/3 green and white Phyrexian Hydra creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("HYDRA")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewPhyrexianMinionToken
// ========================================

// NewPhyrexianMinionToken creates a X/X black Phyrexian Minion creature token.
func NewPhyrexianMinionToken() *Token {
	tok := NewToken("PhyrexianMinionToken", "X/X black Phyrexian Minion creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("MINION")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewPhyrexianMiteToken
// ========================================

// NewPhyrexianMiteToken creates a token.
func NewPhyrexianMiteToken() *Token {
	tok := NewToken("PhyrexianMiteToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("MITE")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPhyrexianMyrToken
// ========================================

// NewPhyrexianMyrToken creates a 1/1 colorless Phyrexian Myr artifact creature token.
func NewPhyrexianMyrToken() *Token {
	tok := NewToken("PhyrexianMyrToken", "1/1 colorless Phyrexian Myr artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("MYR")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPhyrexianRebirthHorrorToken
// ========================================

// NewPhyrexianRebirthHorrorToken creates a X/X colorless Phyrexian Horror artifact creature token.
func NewPhyrexianRebirthHorrorToken() *Token {
	tok := NewToken("PhyrexianRebirthHorrorToken", "X/X colorless Phyrexian Horror artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("HORROR")
	return tok
}

// ========================================
// NewPhyrexianSaprolingToken
// ========================================

// NewPhyrexianSaprolingToken creates a 1/1 green Phyrexian Saproling creature token.
func NewPhyrexianSaprolingToken() *Token {
	tok := NewToken("PhyrexianSaprolingToken", "1/1 green Phyrexian Saproling creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("SAPROLING")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPhyrexianToken
// ========================================

// NewPhyrexianToken creates a 2/2 black Phyrexian creature token.
func NewPhyrexianToken() *Token {
	tok := NewToken("PhyrexianToken", "2/2 black Phyrexian creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewPhyrexianWurm12DeathtouchToken
// ========================================

// NewPhyrexianWurm12DeathtouchToken creates a 1/2 black Phyrexian Wurm artifact creature token with deathtouch.
func NewPhyrexianWurm12DeathtouchToken() *Token {
	tok := NewToken("PhyrexianWurm12DeathtouchToken", "1/2 black Phyrexian Wurm artifact creature token with deathtouch")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewPhyrexianWurm21LifelinkToken
// ========================================

// NewPhyrexianWurm21LifelinkToken creates a 2/1 black Phyrexian Wurm artifact creature token with lifelink.
func NewPhyrexianWurm21LifelinkToken() *Token {
	tok := NewToken("PhyrexianWurm21LifelinkToken", "2/1 black Phyrexian Wurm artifact creature token with lifelink")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewPhyrexianWurmToken
// ========================================

// NewPhyrexianWurmToken creates a X/X green Phyrexian Wurm creature token with trample and toxic 1.
func NewPhyrexianWurmToken() *Token {
	tok := NewToken("PhyrexianWurmToken", "X/X green Phyrexian Wurm creature token with trample and toxic 1")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewPhyrexianZombieToken
// ========================================

// NewPhyrexianZombieToken creates a 2/2 black Phyrexian Zombie creature token.
func NewPhyrexianZombieToken() *Token {
	tok := NewToken("PhyrexianZombieToken", "2/2 black Phyrexian Zombie creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewPilotCrewToken
// ========================================

// NewPilotCrewToken creates a token.
func NewPilotCrewToken() *Token {
	tok := NewToken("PilotCrewToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PILOT")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPilotSaddleCrewToken
// ========================================

// NewPilotSaddleCrewToken creates a token.
func NewPilotSaddleCrewToken() *Token {
	tok := NewToken("PilotSaddleCrewToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PILOT")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPincherToken
// ========================================

// NewPincherToken creates a 2/2 colorless Pincher creature token.
func NewPincherToken() *Token {
	tok := NewToken("PincherToken", "2/2 colorless Pincher creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PINCHER")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewPirateRedToken
// ========================================

// NewPirateRedToken creates a 1/1 red Pirate creature token with menace and haste.
func NewPirateRedToken() *Token {
	tok := NewToken("PirateRedToken", "1/1 red Pirate creature token with menace and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PIRATE")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewPirateToken
// ========================================

// NewPirateToken creates a 2/2 black Pirate creature token with menace.
func NewPirateToken() *Token {
	tok := NewToken("PirateToken", "2/2 black Pirate creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PIRATE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewPlaguebearerOfNurgleToken
// ========================================

// NewPlaguebearerOfNurgleToken creates a 1/3 black Demon creature token named Plaguebearer of Nurgle.
func NewPlaguebearerOfNurgleToken() *Token {
	tok := NewToken("PlaguebearerOfNurgleToken", "1/3 black Demon creature token named Plaguebearer of Nurgle")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 3)
	return tok
}

// ========================================
// NewPlanewideCelebrationToken
// ========================================

// NewPlanewideCelebrationToken creates a 2/2 Citizen creature token that's all colors.
func NewPlanewideCelebrationToken() *Token {
	tok := NewToken("PlanewideCelebrationToken", "2/2 Citizen creature token that's all colors")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CITIZEN")
	tok.SetColor(Color{White: true, Blue: true, Black: true, Red: true, Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewPlant11Token
// ========================================

// NewPlant11Token creates a 1/1 green Plant creature token.
func NewPlant11Token() *Token {
	tok := NewToken("Plant11Token", "1/1 green Plant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PLANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewPlantToken
// ========================================

// NewPlantToken creates a 0/1 green Plant creature token.
func NewPlantToken() *Token {
	tok := NewToken("PlantToken", "0/1 green Plant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PLANT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewPlantWarriorToken
// ========================================

// NewPlantWarriorToken creates a 4/2 green Plant Warrior creature token with reach.
func NewPlantWarriorToken() *Token {
	tok := NewToken("PlantWarriorToken", "4/2 green Plant Warrior creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PLANT")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 2)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewPorgToken
// ========================================

// NewPorgToken creates a token.
func NewPorgToken() *Token {
	tok := NewToken("PorgToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewPowerstoneToken
// ========================================

// NewPowerstoneToken creates a Powerstone token.
func NewPowerstoneToken() *Token {
	tok := NewToken("PowerstoneToken", "Powerstone token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("POWERSTONE")
	return tok
}

// ========================================
// NewPrimoTheIndivisibleToken
// ========================================

// NewPrimoTheIndivisibleToken creates a token.
func NewPrimoTheIndivisibleToken() *Token {
	tok := NewToken("PrimoTheIndivisibleToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FRACTAL")
	tok.SetColor(Color{Blue: true, Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewPrismToken
// ========================================

// NewPrismToken creates a 0/1 colorless Prism artifact creature token.
func NewPrismToken() *Token {
	tok := NewToken("PrismToken", "0/1 colorless Prism artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PRISM")
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewPurphorossInterventionToken
// ========================================

// NewPurphorossInterventionToken creates a X/1 red Elemental creature token with trample and haste.
func NewPurphorossInterventionToken() *Token {
	tok := NewToken("PurphorossInterventionToken", "X/1 red Elemental creature token with trample and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewPursuedWhaleToken
// ========================================

// NewPursuedWhaleToken creates a token.
func NewPursuedWhaleToken() *Token {
	tok := NewToken("PursuedWhaleToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PIRATE")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewQueenMarchesaAssassinToken
// ========================================

// NewQueenMarchesaAssassinToken creates a 1/1 black Assassin creature token with deathtouch and haste.
func NewQueenMarchesaAssassinToken() *Token {
	tok := NewToken("QueenMarchesaAssassinToken", "1/1 black Assassin creature token with deathtouch and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ASSASSIN")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewQuestForTheGravelordZombieToken
// ========================================

// NewQuestForTheGravelordZombieToken creates a 5/5 black Zombie Giant creature token.
func NewQuestForTheGravelordZombieToken() *Token {
	tok := NewToken("QuestForTheGravelordZombieToken", "5/5 black Zombie Giant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("GIANT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewRabbitToken
// ========================================

// NewRabbitToken creates a 1/1 white Rabbit creature token.
func NewRabbitToken() *Token {
	tok := NewToken("RabbitToken", "1/1 white Rabbit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RABBIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRabidSheepToken
// ========================================

// NewRabidSheepToken creates a 2/2 green Sheep creature token named Rabid Sheep.
func NewRabidSheepToken() *Token {
	tok := NewToken("RabidSheepToken", "2/2 green Sheep creature token named Rabid Sheep")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHEEP")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewRaccoonToken
// ========================================

// NewRaccoonToken creates a 3/3 green Raccoon creature token.
func NewRaccoonToken() *Token {
	tok := NewToken("RaccoonToken", "3/3 green Raccoon creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RACCOON")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewRagavanToken
// ========================================

// NewRagavanToken creates a Ragavan, a legendary 2/1 red Monkey creature token.
func NewRagavanToken() *Token {
	tok := NewToken("RagavanToken", "Ragavan, a legendary 2/1 red Monkey creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("MONKEY")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewRakdosGuildmageGoblinToken
// ========================================

// NewRakdosGuildmageGoblinToken creates a 2/1 red Goblin creature token with haste.
func NewRakdosGuildmageGoblinToken() *Token {
	tok := NewToken("RakdosGuildmageGoblinToken", "2/1 red Goblin creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewRasputinKnightToken
// ========================================

// NewRasputinKnightToken creates a 2/2 white Knight creature token with protection from red.
func NewRasputinKnightToken() *Token {
	tok := NewToken("RasputinKnightToken", "2/2 white Knight creature token with protection from red")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewRat11LifelinkToken
// ========================================

// NewRat11LifelinkToken creates a 1/1 black Rat creature token with lifelink.
func NewRat11LifelinkToken() *Token {
	tok := NewToken("Rat11LifelinkToken", "1/1 black Rat creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewRatCantBlockToken
// ========================================

// NewRatCantBlockToken creates a token.
func NewRatCantBlockToken() *Token {
	tok := NewToken("RatCantBlockToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRatRogueToken
// ========================================

// NewRatRogueToken creates a 1/1 black Rat Rogue creature token.
func NewRatRogueToken() *Token {
	tok := NewToken("RatRogueToken", "1/1 black Rat Rogue creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RAT")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRatToken
// ========================================

// NewRatToken creates a 1/1 black Rat creature token.
func NewRatToken() *Token {
	tok := NewToken("RatToken", "1/1 black Rat creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRebelRedToken
// ========================================

// NewRebelRedToken creates a 2/2 red Rebel creature token.
func NewRebelRedToken() *Token {
	tok := NewToken("RebelRedToken", "2/2 red Rebel creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("REBEL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewRebelStarshipToken
// ========================================

// NewRebelStarshipToken creates a 2/3 blue Rebel Starship artifact creature tokens with spaceflight name B-Wing.
func NewRebelStarshipToken() *Token {
	tok := NewToken("RebelStarshipToken", "2/3 blue Rebel Starship artifact creature tokens with spaceflight name B-Wing")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("REBEL")
	tok.AddSubtype("STARSHIP")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 3)
	return tok
}

// ========================================
// NewRebelToken
// ========================================

// NewRebelToken creates a 1/1 white Rebel creature token.
func NewRebelToken() *Token {
	tok := NewToken("RebelToken", "1/1 white Rebel creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("REBEL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRedElementalToken
// ========================================

// NewRedElementalToken creates a 1/1 red Elemental creature token.
func NewRedElementalToken() *Token {
	tok := NewToken("RedElementalToken", "1/1 red Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRedGreenBeastToken
// ========================================

// NewRedGreenBeastToken creates a 4/4 red and green Beast creature token with trample.
func NewRedGreenBeastToken() *Token {
	tok := NewToken("RedGreenBeastToken", "4/4 red and green Beast creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewRedHumanToken
// ========================================

// NewRedHumanToken creates a 1/1 red Human creature token.
func NewRedHumanToken() *Token {
	tok := NewToken("RedHumanToken", "1/1 red Human creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRedWarriorToken
// ========================================

// NewRedWarriorToken creates a 1/1 red Warrior creature token.
func NewRedWarriorToken() *Token {
	tok := NewToken("RedWarriorToken", "1/1 red Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewRedWhiteGolemToken
// ========================================

// NewRedWhiteGolemToken creates a 4/4 red and white Golem artifact creature token.
func NewRedWhiteGolemToken() *Token {
	tok := NewToken("RedWhiteGolemToken", "4/4 red and white Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewRedWolfToken
// ========================================

// NewRedWolfToken creates a 3/2 red Wolf creature token.
func NewRedWolfToken() *Token {
	tok := NewToken("RedWolfToken", "3/2 red Wolf creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewReefWormFishToken
// ========================================

// NewReefWormFishToken creates a token.
func NewReefWormFishToken() *Token {
	tok := NewToken("ReefWormFishToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("FISH")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewReefWormWhaleToken
// ========================================

// NewReefWormWhaleToken creates a token.
func NewReefWormWhaleToken() *Token {
	tok := NewToken("ReefWormWhaleToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WHALE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(6, 6)
	return tok
}

// ========================================
// NewReflectionBlueToken
// ========================================

// NewReflectionBlueToken creates a 3/2 blue Reflection creature token.
func NewReflectionBlueToken() *Token {
	tok := NewToken("ReflectionBlueToken", "3/2 blue Reflection creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("REFLECTION")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewReflectionPureToken
// ========================================

// NewReflectionPureToken creates a X/X white Reflection creature token, where X is the mana value of that spell.
func NewReflectionPureToken() *Token {
	tok := NewToken("ReflectionPureToken", "X/X white Reflection creature token, where X is the mana value of that spell")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("REFLECTION")
	tok.SetColor(Color{White: true})
	return tok
}

// ========================================
// NewReflectionToken
// ========================================

// NewReflectionToken creates a 2/2 white Reflection creature token.
func NewReflectionToken() *Token {
	tok := NewToken("ReflectionToken", "2/2 white Reflection creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("REFLECTION")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewRekindlingPhoenixToken
// ========================================

// NewRekindlingPhoenixToken creates a token.
func NewRekindlingPhoenixToken() *Token {
	tok := NewToken("RekindlingPhoenixToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(0, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewRelicRobberToken
// ========================================

// NewRelicRobberToken creates a token.
func NewRelicRobberToken() *Token {
	tok := NewToken("RelicRobberToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewReliquaryDragonToken
// ========================================

// NewReliquaryDragonToken creates a token.
func NewReliquaryDragonToken() *Token {
	tok := NewToken("ReliquaryDragonToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{White: true, Blue: true, Black: true, Red: true, Green: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewRenownedWeaverSpiderToken
// ========================================

// NewRenownedWeaverSpiderToken creates a 1/3 green Spider enchantment creature token with reach.
func NewRenownedWeaverSpiderToken() *Token {
	tok := NewToken("RenownedWeaverSpiderToken", "1/3 green Spider enchantment creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("SPIDER")
	tok.SetPowerToughness(1, 3)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewReplicatedRingToken
// ========================================

// NewReplicatedRingToken creates a token.
func NewReplicatedRingToken() *Token {
	tok := NewToken("ReplicatedRingToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewResearchDevelopmentToken
// ========================================

// NewResearchDevelopmentToken creates a 3/1 red Elemental creature token.
func NewResearchDevelopmentToken() *Token {
	tok := NewToken("ResearchDevelopmentToken", "3/1 red Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	return tok
}

// ========================================
// NewRhinoToken
// ========================================

// NewRhinoToken creates a 4/4 green Rhino creature token with trample.
func NewRhinoToken() *Token {
	tok := NewToken("RhinoToken", "4/4 green Rhino creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RHINO")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewRhinoWarriorToken
// ========================================

// NewRhinoWarriorToken creates a 4/4 green Rhino Warrior creature token.
func NewRhinoWarriorToken() *Token {
	tok := NewToken("RhinoWarriorToken", "4/4 green Rhino Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RHINO")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewRhonassLastStandToken
// ========================================

// NewRhonassLastStandToken creates a 5/4 green Snake creature token.
func NewRhonassLastStandToken() *Token {
	tok := NewToken("RhonassLastStandToken", "5/4 green Snake creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAKE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 4)
	return tok
}

// ========================================
// NewRiftmarkedKnightToken
// ========================================

// NewRiftmarkedKnightToken creates a 2/2 black Knight creature token with flanking, protection from white, and haste.
func NewRiftmarkedKnightToken() *Token {
	tok := NewToken("RiftmarkedKnightToken", "2/2 black Knight creature token with flanking, protection from white, and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewRiptideReplicatorToken
// ========================================

// NewRiptideReplicatorToken creates a X/X creature token of the chosen color and type.
func NewRiptideReplicatorToken() *Token {
	tok := NewToken("RiptideReplicatorToken", "X/X creature token of the chosen color and type")
	tok.AddCardType(CardTypeCreature)
	return tok
}

// ========================================
// NewRiseOfTheAntsInsectToken
// ========================================

// NewRiseOfTheAntsInsectToken creates a 3/3 green Insect creature token.
func NewRiseOfTheAntsInsectToken() *Token {
	tok := NewToken("RiseOfTheAntsInsectToken", "3/3 green Insect creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewRitualOfTheReturnedZombieToken
// ========================================

// NewRitualOfTheReturnedZombieToken creates a black Zombie creature token with power equal to the exiled card's power and toughness equal to the exiled card's toughness.
func NewRitualOfTheReturnedZombieToken() *Token {
	tok := NewToken("RitualOfTheReturnedZombieToken", "black Zombie creature token with power equal to the exiled card's power and toughness equal to the exiled card's toughness")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	return tok
}

// ========================================
// NewRobot33Token
// ========================================

// NewRobot33Token creates a 3/3 colorless Robot artifact creature token.
func NewRobot33Token() *Token {
	tok := NewToken("Robot33Token", "3/3 colorless Robot artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ROBOT")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewRobotBlueToken
// ========================================

// NewRobotBlueToken creates a 3/3 blue Robot Warrior artifact creature token.
func NewRobotBlueToken() *Token {
	tok := NewToken("RobotBlueToken", "3/3 blue Robot Warrior artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ROBOT")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewRobotCantBlockToken
// ========================================

// NewRobotCantBlockToken creates a token.
func NewRobotCantBlockToken() *Token {
	tok := NewToken("RobotCantBlockToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ROBOT")
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewRobotFlyingToken
// ========================================

// NewRobotFlyingToken creates a 1/1 colorless Robot artifact creature tokens with flying.
func NewRobotFlyingToken() *Token {
	tok := NewToken("RobotFlyingToken", "1/1 colorless Robot artifact creature tokens with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ROBOT")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewRobotToken
// ========================================

// NewRobotToken creates a 2/2 colorless Robot artifact creature token.
func NewRobotToken() *Token {
	tok := NewToken("RobotToken", "2/2 colorless Robot artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ROBOT")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewRocEggToken
// ========================================

// NewRocEggToken creates a 3/3 white Bird creature token with flying.
func NewRocEggToken() *Token {
	tok := NewToken("RocEggToken", "3/3 white Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewRockToken
// ========================================

// NewRockToken creates a token.
func NewRockToken() *Token {
	tok := NewToken("RockToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("EQUIPMENT")
	return tok
}

// ========================================
// NewRogueToken
// ========================================

// NewRogueToken creates a 2/2 black Rogue creature token.
func NewRogueToken() *Token {
	tok := NewToken("RogueToken", "2/2 black Rogue creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewRoyalGuardToken
// ========================================

// NewRoyalGuardToken creates a 2/2 red Soldier creature token with first strike named Royal Guard.
func NewRoyalGuardToken() *Token {
	tok := NewToken("RoyalGuardToken", "2/2 red Soldier creature token with first strike named Royal Guard")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("first strike")
	return tok
}

// ========================================
// NewRoyalRoleToken
// ========================================

// NewRoyalRoleToken creates a Royal Role token.
func NewRoyalRoleToken() *Token {
	tok := NewToken("RoyalRoleToken", "Royal Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	return tok
}

// ========================================
// NewRukhEggBirdToken
// ========================================

// NewRukhEggBirdToken creates a 4/4 red Bird creature token with flying.
func NewRukhEggBirdToken() *Token {
	tok := NewToken("RukhEggBirdToken", "4/4 red Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSalamanderWarriorToken
// ========================================

// NewSalamanderWarriorToken creates a 4/3 blue Salamander Warrior creature token.
func NewSalamanderWarriorToken() *Token {
	tok := NewToken("SalamanderWarriorToken", "4/3 blue Salamander Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SALAMANDER")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(4, 3)
	return tok
}

// ========================================
// NewSamuraiToken
// ========================================

// NewSamuraiToken creates a 2/2 white Samurai creature token with vigilance..
func NewSamuraiToken() *Token {
	tok := NewToken("SamuraiToken", "2/2 white Samurai creature token with vigilance.")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SAMURAI")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewSandWarriorToken
// ========================================

// NewSandWarriorToken creates a 1/1 red, green, and white Sand Warrior creature token.
func NewSandWarriorToken() *Token {
	tok := NewToken("SandWarriorToken", "1/1 red, green, and white Sand Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SAND")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true, Red: true, Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSaprolingBurstToken
// ========================================

// NewSaprolingBurstToken creates a token.
func NewSaprolingBurstToken() *Token {
	tok := NewToken("SaprolingBurstToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SAPROLING")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSaprolingToken
// ========================================

// NewSaprolingToken creates a 1/1 green Saproling creature token.
func NewSaprolingToken() *Token {
	tok := NewToken("SaprolingToken", "1/1 green Saproling creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SAPROLING")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSatyrCantBlockToken
// ========================================

// NewSatyrCantBlockToken creates a token.
func NewSatyrCantBlockToken() *Token {
	tok := NewToken("SatyrCantBlockToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SATYR")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSatyrNyxSmithElementalToken
// ========================================

// NewSatyrNyxSmithElementalToken creates a 3/1 red Elemental enchantment creature token with haste.
func NewSatyrNyxSmithElementalToken() *Token {
	tok := NewToken("SatyrNyxSmithElementalToken", "3/1 red Elemental enchantment creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewScionOfTheDeepToken
// ========================================

// NewScionOfTheDeepToken creates a Scion of the Deep, a legendary 8/8 blue Octopus creature token.
func NewScionOfTheDeepToken() *Token {
	tok := NewToken("ScionOfTheDeepToken", "Scion of the Deep, a legendary 8/8 blue Octopus creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("OCTOPUS")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(8, 8)
	return tok
}

// ========================================
// NewScorpionDragonToken
// ========================================

// NewScorpionDragonToken creates a 4/4 red Scorpion Dragon creature token with flying and haste.
func NewScorpionDragonToken() *Token {
	tok := NewToken("ScorpionDragonToken", "4/4 red Scorpion Dragon creature token with flying and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SCORPION")
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewScrapToken
// ========================================

// NewScrapToken creates a colorless artifact token named Scrap.
func NewScrapToken() *Token {
	tok := NewToken("ScrapToken", "colorless artifact token named Scrap")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewSeizeTheStormElementalToken
// ========================================

// NewSeizeTheStormElementalToken creates a token.
func NewSeizeTheStormElementalToken() *Token {
	tok := NewToken("SeizeTheStormElementalToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(0, 0)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewSengirNosferatuBatToken
// ========================================

// NewSengirNosferatuBatToken creates a 1/2 black Bat creature token with flying.
func NewSengirNosferatuBatToken() *Token {
	tok := NewToken("SengirNosferatuBatToken", "1/2 black Bat creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSerfToken
// ========================================

// NewSerfToken creates a 0/1 black Serf creature token.
func NewSerfToken() *Token {
	tok := NewToken("SerfToken", "0/1 black Serf creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SERF")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewSerpentGeneratorSnakeToken
// ========================================

// NewSerpentGeneratorSnakeToken creates a token.
func NewSerpentGeneratorSnakeToken() *Token {
	tok := NewToken("SerpentGeneratorSnakeToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAKE")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewServoToken
// ========================================

// NewServoToken creates a 1/1 colorless Servo artifact creature token.
func NewServoToken() *Token {
	tok := NewToken("ServoToken", "1/1 colorless Servo artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SERVO")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSettlementToken
// ========================================

// NewSettlementToken creates a Settlement token.
func NewSettlementToken() *Token {
	tok := NewToken("SettlementToken", "Settlement token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewShapeshifter32Token
// ========================================

// NewShapeshifter32Token creates a 3/2 colorless Shapeshifter creature token with changeling.
func NewShapeshifter32Token() *Token {
	tok := NewToken("Shapeshifter32Token", "3/2 colorless Shapeshifter creature token with changeling")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHAPESHIFTER")
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewShapeshifterBlueToken
// ========================================

// NewShapeshifterBlueToken creates a 2/2 blue Shapeshifter creature token with changeling.
func NewShapeshifterBlueToken() *Token {
	tok := NewToken("ShapeshifterBlueToken", "2/2 blue Shapeshifter creature token with changeling")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHAPESHIFTER")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewShapeshifterDeathtouchToken
// ========================================

// NewShapeshifterDeathtouchToken creates a X/X colorless Shapeshifter creature token with changeling and deathtouch.
func NewShapeshifterDeathtouchToken() *Token {
	tok := NewToken("ShapeshifterDeathtouchToken", "X/X colorless Shapeshifter creature token with changeling and deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHAPESHIFTER")
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewShapeshifterToken
// ========================================

// NewShapeshifterToken creates a 2/2 colorless Shapeshifter creature token with changeling.
func NewShapeshifterToken() *Token {
	tok := NewToken("ShapeshifterToken", "2/2 colorless Shapeshifter creature token with changeling")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHAPESHIFTER")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewShardToken
// ========================================

// NewShardToken creates a Shard token.
func NewShardToken() *Token {
	tok := NewToken("ShardToken", "Shard token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("SHARD")
	return tok
}

// ========================================
// NewShark33Token
// ========================================

// NewShark33Token creates a 3/3 blue Shark creature token.
func NewShark33Token() *Token {
	tok := NewToken("Shark33Token", "3/3 blue Shark creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHARK")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewSharkToken
// ========================================

// NewSharkToken creates a X/X blue Shark creature token with flying.
func NewSharkToken() *Token {
	tok := NewToken("SharkToken", "X/X blue Shark creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHARK")
	tok.SetColor(Color{Blue: true})
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSheepToken
// ========================================

// NewSheepToken creates a 0/1 green Sheep creature token.
func NewSheepToken() *Token {
	tok := NewToken("SheepToken", "0/1 green Sheep creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHEEP")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewSheepWhiteToken
// ========================================

// NewSheepWhiteToken creates a 1/1 white Sheep creature token.
func NewSheepWhiteToken() *Token {
	tok := NewToken("SheepWhiteToken", "1/1 white Sheep creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SHEEP")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewShrineToken
// ========================================

// NewShrineToken creates a 1/1 colorless Shrine enchantment creature token.
func NewShrineToken() *Token {
	tok := NewToken("ShrineToken", "1/1 colorless Shrine enchantment creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("SHRINE")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSkeletonMenaceToken
// ========================================

// NewSkeletonMenaceToken creates a 4/1 black Skeleton creature token with menace.
func NewSkeletonMenaceToken() *Token {
	tok := NewToken("SkeletonMenaceToken", "4/1 black Skeleton creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SKELETON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(4, 1)
	return tok
}

// ========================================
// NewSkeletonPirateToken
// ========================================

// NewSkeletonPirateToken creates a 2/2 black Skeleton Pirate creature token.
func NewSkeletonPirateToken() *Token {
	tok := NewToken("SkeletonPirateToken", "2/2 black Skeleton Pirate creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SKELETON")
	tok.AddSubtype("PIRATE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSkeletonRegenerateToken
// ========================================

// NewSkeletonRegenerateToken creates a token.
func NewSkeletonRegenerateToken() *Token {
	tok := NewToken("SkeletonRegenerateToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SKELETON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSkeletonToken
// ========================================

// NewSkeletonToken creates a 1/1 black Skeleton creature token.
func NewSkeletonToken() *Token {
	tok := NewToken("SkeletonToken", "1/1 black Skeleton creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SKELETON")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSliverArmyToken
// ========================================

// NewSliverArmyToken creates a 0/0 black Sliver Army creature token.
func NewSliverArmyToken() *Token {
	tok := NewToken("SliverArmyToken", "0/0 black Sliver Army creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SLIVER")
	tok.AddSubtype("ARMY")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewSliverToken
// ========================================

// NewSliverToken creates a 1/1 colorless Sliver creature token.
func NewSliverToken() *Token {
	tok := NewToken("SliverToken", "1/1 colorless Sliver creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SLIVER")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSlugToken
// ========================================

// NewSlugToken creates a 1/1 black Slug creature token.
func NewSlugToken() *Token {
	tok := NewToken("SlugToken", "1/1 black Slug creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SLUG")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSmaugToken
// ========================================

// NewSmaugToken creates a token.
func NewSmaugToken() *Token {
	tok := NewToken("SmaugToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(6, 6)
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewSmokeBlessingToken
// ========================================

// NewSmokeBlessingToken creates a token.
func NewSmokeBlessingToken() *Token {
	tok := NewToken("SmokeBlessingToken", "")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.SetColor(Color{Red: true})
	return tok
}

// ========================================
// NewSnailToken
// ========================================

// NewSnailToken creates a 1/1 black Snail creature token.
func NewSnailToken() *Token {
	tok := NewToken("SnailToken", "1/1 black Snail creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAIL")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSnakeToken
// ========================================

// NewSnakeToken creates a 1/1 green Snake creature token.
func NewSnakeToken() *Token {
	tok := NewToken("SnakeToken", "1/1 green Snake creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SNAKE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSoldier22Token
// ========================================

// NewSoldier22Token creates a 2/2 white Soldier creature token.
func NewSoldier22Token() *Token {
	tok := NewToken("Soldier22Token", "2/2 white Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSoldierArtifactToken
// ========================================

// NewSoldierArtifactToken creates a 1/1 colorless Soldier artifact creature token.
func NewSoldierArtifactToken() *Token {
	tok := NewToken("SoldierArtifactToken", "1/1 colorless Soldier artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSoldierFirebendingToken
// ========================================

// NewSoldierFirebendingToken creates a 2/2 red Soldier creature token with firebending 1.
func NewSoldierFirebendingToken() *Token {
	tok := NewToken("SoldierFirebendingToken", "2/2 red Soldier creature token with firebending 1")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSoldierLifelinkToken
// ========================================

// NewSoldierLifelinkToken creates a 1/1 white Soldier creature token with lifelink.
func NewSoldierLifelinkToken() *Token {
	tok := NewToken("SoldierLifelinkToken", "1/1 white Soldier creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewSoldierRedToken
// ========================================

// NewSoldierRedToken creates a 2/2 red Soldier creature token.
func NewSoldierRedToken() *Token {
	tok := NewToken("SoldierRedToken", "2/2 red Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSoldierToken
// ========================================

// NewSoldierToken creates a 1/1 white Soldier creature token.
func NewSoldierToken() *Token {
	tok := NewToken("SoldierToken", "1/1 white Soldier creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSoldierVigilanceToken
// ========================================

// NewSoldierVigilanceToken creates a 2/2 white Soldier creature token with vigilance.
func NewSoldierVigilanceToken() *Token {
	tok := NewToken("SoldierVigilanceToken", "2/2 white Soldier creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SOLDIER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewSorcererRoleToken
// ========================================

// NewSorcererRoleToken creates a Sorcerer Role token.
func NewSorcererRoleToken() *Token {
	tok := NewToken("SorcererRoleToken", "Sorcerer Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	return tok
}

// ========================================
// NewSorinLordOfInnistradVampireToken
// ========================================

// NewSorinLordOfInnistradVampireToken creates a 1/1 black Vampire creature token with lifelink.
func NewSorinLordOfInnistradVampireToken() *Token {
	tok := NewToken("SorinLordOfInnistradVampireToken", "1/1 black Vampire creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewSoundTheCallToken
// ========================================

// NewSoundTheCallToken creates a token.
func NewSoundTheCallToken() *Token {
	tok := NewToken("SoundTheCallToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSparkElementalToken
// ========================================

// NewSparkElementalToken creates a token.
func NewSparkElementalToken() *Token {
	tok := NewToken("SparkElementalToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("haste")
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewSpawnToken
// ========================================

// NewSpawnToken creates a 3/3 red Spawn creature token.
func NewSpawnToken() *Token {
	tok := NewToken("SpawnToken", "3/3 red Spawn creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPAWN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewSpawningGroundsBeastToken
// ========================================

// NewSpawningGroundsBeastToken creates a 5/5 green Beast creature token with trample.
func NewSpawningGroundsBeastToken() *Token {
	tok := NewToken("SpawningGroundsBeastToken", "5/5 green Beast creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BEAST")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewSpawningPitToken
// ========================================

// NewSpawningPitToken creates a 2/2 colorless Spawn artifact creature token.
func NewSpawningPitToken() *Token {
	tok := NewToken("SpawningPitToken", "2/2 colorless Spawn artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPAWN")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSpellgorgerWeirdToken
// ========================================

// NewSpellgorgerWeirdToken creates a Spellgorger Weird token.
func NewSpellgorgerWeirdToken() *Token {
	tok := NewToken("SpellgorgerWeirdToken", "Spellgorger Weird token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WEIRD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSpider21Token
// ========================================

// NewSpider21Token creates a 2/1 green Spider creature token with reach.
func NewSpider21Token() *Token {
	tok := NewToken("Spider21Token", "2/1 green Spider creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIDER")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 1)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewSpider22Token
// ========================================

// NewSpider22Token creates a 2/2 green Spider creature token with reach.
func NewSpider22Token() *Token {
	tok := NewToken("Spider22Token", "2/2 green Spider creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIDER")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewSpiderToken
// ========================================

// NewSpiderToken creates a 1/2 green Spider creature token with reach.
func NewSpiderToken() *Token {
	tok := NewToken("SpiderToken", "1/2 green Spider creature token with reach")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIDER")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewSpikeToken
// ========================================

// NewSpikeToken creates a 1/1 green Spike creature token.
func NewSpikeToken() *Token {
	tok := NewToken("SpikeToken", "1/1 green Spike creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIKE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSpirit22Token
// ========================================

// NewSpirit22Token creates a 2/2 white Spirit creature token with flying.
func NewSpirit22Token() *Token {
	tok := NewToken("Spirit22Token", "2/2 white Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSpirit31Token
// ========================================

// NewSpirit31Token creates a 3/1 white Spirit creature token with flying.
func NewSpirit31Token() *Token {
	tok := NewToken("Spirit31Token", "3/1 white Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSpirit32Token
// ========================================

// NewSpirit32Token creates a 3/2 red and white Spirit creature token.
func NewSpirit32Token() *Token {
	tok := NewToken("Spirit32Token", "3/2 red and white Spirit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewSpiritBlueToken
// ========================================

// NewSpiritBlueToken creates a 1/1 blue Spirit creature token with flying.
func NewSpiritBlueToken() *Token {
	tok := NewToken("SpiritBlueToken", "1/1 blue Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSpiritClericToken
// ========================================

// NewSpiritClericToken creates a token.
func NewSpiritClericToken() *Token {
	tok := NewToken("SpiritClericToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.AddSubtype("CLERIC")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewSpiritEvilBorosCharmToken
// ========================================

// NewSpiritEvilBorosCharmToken creates a 1/1 colorless Spirit creature token with lifelink and haste.
func NewSpiritEvilBorosCharmToken() *Token {
	tok := NewToken("SpiritEvilBorosCharmToken", "1/1 colorless Spirit creature token with lifelink and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewSpiritGreenToken
// ========================================

// NewSpiritGreenToken creates a 4/5 green Spirit creature token.
func NewSpiritGreenToken() *Token {
	tok := NewToken("SpiritGreenToken", "4/5 green Spirit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 5)
	return tok
}

// ========================================
// NewSpiritGreenXToken
// ========================================

// NewSpiritGreenXToken creates a X/X green Spirit creature token.
func NewSpiritGreenXToken() *Token {
	tok := NewToken("SpiritGreenXToken", "X/X green Spirit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewSpiritRedToken
// ========================================

// NewSpiritRedToken creates a 2/2 red Spirit creature token with menace.
func NewSpiritRedToken() *Token {
	tok := NewToken("SpiritRedToken", "2/2 red Spirit creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewSpiritTeferiToken
// ========================================

// NewSpiritTeferiToken creates a token.
func NewSpiritTeferiToken() *Token {
	tok := NewToken("SpiritTeferiToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewSpiritToken
// ========================================

// NewSpiritToken creates a 1/1 colorless Spirit creature token.
func NewSpiritToken() *Token {
	tok := NewToken("SpiritToken", "1/1 colorless Spirit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSpiritWarriorToken
// ========================================

// NewSpiritWarriorToken creates a X/X black and green Spirit Warrior creature token, where X is the greatest toughness among creatures you control.
func NewSpiritWarriorToken() *Token {
	tok := NewToken("SpiritWarriorToken", "X/X black and green Spirit Warrior creature token, where X is the greatest toughness among creatures you control")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Black: true, Green: true})
	return tok
}

// ========================================
// NewSpiritWhiteToken
// ========================================

// NewSpiritWhiteToken creates a 1/1 white Spirit creature token with flying.
func NewSpiritWhiteToken() *Token {
	tok := NewToken("SpiritWhiteToken", "1/1 white Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSpiritWorldToken
// ========================================

// NewSpiritWorldToken creates a token.
func NewSpiritWorldToken() *Token {
	tok := NewToken("SpiritWorldToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSpiritXXToken
// ========================================

// NewSpiritXXToken creates a token.
func NewSpiritXXToken() *Token {
	tok := NewToken("SpiritXXToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true})
	return tok
}

// ========================================
// NewSplinterToken
// ========================================

// NewSplinterToken creates a 1/1 green Splinter creature token.
func NewSplinterToken() *Token {
	tok := NewToken("SplinterToken", "1/1 green Splinter creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPLINTER")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSpyMasterGoblinToken
// ========================================

// NewSpyMasterGoblinToken creates a token.
func NewSpyMasterGoblinToken() *Token {
	tok := NewToken("SpyMasterGoblinToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSquidToken
// ========================================

// NewSquidToken creates a 1/1 blue Squid creature token with islandwalk.
func NewSquidToken() *Token {
	tok := NewToken("SquidToken", "1/1 blue Squid creature token with islandwalk")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SQUID")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSquirrelToken
// ========================================

// NewSquirrelToken creates a 1/1 green Squirrel creature token.
func NewSquirrelToken() *Token {
	tok := NewToken("SquirrelToken", "1/1 green Squirrel creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SQUIRREL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewStanggTwinToken
// ========================================

// NewStanggTwinToken creates a Stangg Twin, a legendary 3/4 red and green Human Warrior creature token.
func NewStanggTwinToken() *Token {
	tok := NewToken("StanggTwinToken", "Stangg Twin, a legendary 3/4 red and green Human Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(3, 4)
	return tok
}

// ========================================
// NewStarfishToken
// ========================================

// NewStarfishToken creates a 0/1 blue Starfish creature token.
func NewStarfishToken() *Token {
	tok := NewToken("StarfishToken", "0/1 blue Starfish creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("STARFISH")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewStitcherGeralfZombieToken
// ========================================

// NewStitcherGeralfZombieToken creates a X/X blue Zombie creature token.
func NewStitcherGeralfZombieToken() *Token {
	tok := NewToken("StitcherGeralfZombieToken", "X/X blue Zombie creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Blue: true})
	return tok
}

// ========================================
// NewStitchersApprenticeHomunculusToken
// ========================================

// NewStitchersApprenticeHomunculusToken creates a 2/2 blue Homunculus creature token.
func NewStitchersApprenticeHomunculusToken() *Token {
	tok := NewToken("StitchersApprenticeHomunculusToken", "2/2 blue Homunculus creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HOMUNCULUS")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewStoneIdolToken
// ========================================

// NewStoneIdolToken creates a 6/12 colorless Construct artifact creature token with trample.
func NewStoneIdolToken() *Token {
	tok := NewToken("StoneIdolToken", "6/12 colorless Construct artifact creature token with trample")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CONSTRUCT")
	tok.SetPowerToughness(6, 12)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewStormCrowToken
// ========================================

// NewStormCrowToken creates a 1/2 blue Bird creature token with flying named Storm Crow.
func NewStormCrowToken() *Token {
	tok := NewToken("StormCrowToken", "1/2 blue Bird creature token with flying named Storm Crow")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSubterraneanTremorsLizardToken
// ========================================

// NewSubterraneanTremorsLizardToken creates a 8/8 red Lizard creature token.
func NewSubterraneanTremorsLizardToken() *Token {
	tok := NewToken("SubterraneanTremorsLizardToken", "8/8 red Lizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("LIZARD")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(8, 8)
	return tok
}

// ========================================
// NewSurvivorToken
// ========================================

// NewSurvivorToken creates a 1/1 red Survivor creature token.
func NewSurvivorToken() *Token {
	tok := NewToken("SurvivorToken", "1/1 red Survivor creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SURVIVOR")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewSwanSongBirdToken
// ========================================

// NewSwanSongBirdToken creates a 2/2 blue Bird creature token with flying.
func NewSwanSongBirdToken() *Token {
	tok := NewToken("SwanSongBirdToken", "2/2 blue Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewSwordToken
// ========================================

// NewSwordToken creates a token.
func NewSwordToken() *Token {
	tok := NewToken("SwordToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("EQUIPMENT")
	return tok
}

// ========================================
// NewSylvanOfferingTreefolkToken
// ========================================

// NewSylvanOfferingTreefolkToken creates a X/X green Treefolk creature token.
func NewSylvanOfferingTreefolkToken() *Token {
	tok := NewToken("SylvanOfferingTreefolkToken", "X/X green Treefolk creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TREEFOLK")
	tok.SetColor(Color{Green: true})
	return tok
}

// ========================================
// NewTIEFighterToken
// ========================================

// NewTIEFighterToken creates a 1/1 black Starship artifact creature tokens with Spaceflight named TIE Fighter.
func NewTIEFighterToken() *Token {
	tok := NewToken("TIEFighterToken", "1/1 black Starship artifact creature tokens with Spaceflight named TIE Fighter")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("STARSHIP")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTamiyosNotebookToken
// ========================================

// NewTamiyosNotebookToken creates a token.
func NewTamiyosNotebookToken() *Token {
	tok := NewToken("TamiyosNotebookToken", "")
	tok.AddCardType(CardTypeArtifact)
	return tok
}

// ========================================
// NewTarmogoyfToken
// ========================================

// NewTarmogoyfToken creates a Tarmogoyf token.
func NewTarmogoyfToken() *Token {
	tok := NewToken("TarmogoyfToken", "Tarmogoyf token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("LHURGOYF")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	return tok
}

// ========================================
// NewTatsumaDragonToken
// ========================================

// NewTatsumaDragonToken creates a 5/5 blue Dragon Spirit creature token with flying.
func NewTatsumaDragonToken() *Token {
	tok := NewToken("TatsumaDragonToken", "5/5 blue Dragon Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewTentacleToken
// ========================================

// NewTentacleToken creates a 1/1 blue Tentacle creature token.
func NewTentacleToken() *Token {
	tok := NewToken("TentacleToken", "1/1 blue Tentacle creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TENTACLE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTetraviteToken
// ========================================

// NewTetraviteToken creates a token.
func NewTetraviteToken() *Token {
	tok := NewToken("TetraviteToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TETRAVITE")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewTeyoToken
// ========================================

// NewTeyoToken creates a 0/3 white Wall creature token with defender.
func NewTeyoToken() *Token {
	tok := NewToken("TeyoToken", "0/3 white Wall creature token with defender")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 3)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewThatcherHumanToken
// ========================================

// NewThatcherHumanToken creates a 1/1 red Human creature token with haste.
func NewThatcherHumanToken() *Token {
	tok := NewToken("ThatcherHumanToken", "1/1 red Human creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewTheAtropalToken
// ========================================

// NewTheAtropalToken creates a The Atropal, a legendary 4/4 black God Horror creature token with deathtouch.
func NewTheAtropalToken() *Token {
	tok := NewToken("TheAtropalToken", "The Atropal, a legendary 4/4 black God Horror creature token with deathtouch")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOD")
	tok.AddSubtype("HORROR")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewTheBlackjackToken
// ========================================

// NewTheBlackjackToken creates a The Blackjack, a legendary 3/3 colorless Vehicle artifact token with flying and crew 2.
func NewTheBlackjackToken() *Token {
	tok := NewToken("TheBlackjackToken", "The Blackjack, a legendary 3/3 colorless Vehicle artifact token with flying and crew 2")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("VEHICLE")
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewTheEleventhHourToken
// ========================================

// NewTheEleventhHourToken creates a token.
func NewTheEleventhHourToken() *Token {
	tok := NewToken("TheEleventhHourToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HUMAN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTheGirlInTheFireplaceHorseToken
// ========================================

// NewTheGirlInTheFireplaceHorseToken creates a token.
func NewTheGirlInTheFireplaceHorseToken() *Token {
	tok := NewToken("TheGirlInTheFireplaceHorseToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HORSE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewTheGirlInTheFireplaceHumanNobleToken
// ========================================

// NewTheGirlInTheFireplaceHumanNobleToken creates a token.
func NewTheGirlInTheFireplaceHumanNobleToken() *Token {
	tok := NewToken("TheGirlInTheFireplaceHumanNobleToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTheHollowSentinelToken
// ========================================

// NewTheHollowSentinelToken creates a The Hollow Sentinel, a legendary 3/3 colorless Phyrexian Golem artifact creature token.
func NewTheHollowSentinelToken() *Token {
	tok := NewToken("TheHollowSentinelToken", "The Hollow Sentinel, a legendary 3/3 colorless Phyrexian Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewTheLocustGodInsectToken
// ========================================

// NewTheLocustGodInsectToken creates a 1/1 blue and red Insect creature token with flying and haste.
func NewTheLocustGodInsectToken() *Token {
	tok := NewToken("TheLocustGodInsectToken", "1/1 blue and red Insect creature token with flying and haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Blue: true, Red: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewThePrydwenSteelFlagshipHumanKnightToken
// ========================================

// NewThePrydwenSteelFlagshipHumanKnightToken creates a token.
func NewThePrydwenSteelFlagshipHumanKnightToken() *Token {
	tok := NewToken("ThePrydwenSteelFlagshipHumanKnightToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewThopter00ColorlessToken
// ========================================

// NewThopter00ColorlessToken creates a 0/0 colorless Thopter artifact creature token with flying.
func NewThopter00ColorlessToken() *Token {
	tok := NewToken("Thopter00ColorlessToken", "0/0 colorless Thopter artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("THOPTER")
	tok.SetPowerToughness(0, 0)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewThopterColorlessToken
// ========================================

// NewThopterColorlessToken creates a 1/1 colorless Thopter artifact creature token with flying.
func NewThopterColorlessToken() *Token {
	tok := NewToken("ThopterColorlessToken", "1/1 colorless Thopter artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("THOPTER")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewThopterToken
// ========================================

// NewThopterToken creates a 1/1 blue Thopter artifact creature token with flying.
func NewThopterToken() *Token {
	tok := NewToken("ThopterToken", "1/1 blue Thopter artifact creature token with flying")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("THOPTER")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewThrullToken
// ========================================

// NewThrullToken creates a 1/1 black Thrull creature token.
func NewThrullToken() *Token {
	tok := NewToken("ThrullToken", "1/1 black Thrull creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("THRULL")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTidalWaveWallToken
// ========================================

// NewTidalWaveWallToken creates a 5/5 blue Wall creature token with defender.
func NewTidalWaveWallToken() *Token {
	tok := NewToken("TidalWaveWallToken", "5/5 blue Wall creature token with defender")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewTinyToken
// ========================================

// NewTinyToken creates a Tiny, a legendary 2/2 green Dog Detective creature token with trample.
func NewTinyToken() *Token {
	tok := NewToken("TinyToken", "Tiny, a legendary 2/2 green Dog Detective creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewTitanForgeGolemToken
// ========================================

// NewTitanForgeGolemToken creates a 9/9 colorless Golem artifact creature token.
func NewTitanForgeGolemToken() *Token {
	tok := NewToken("TitanForgeGolemToken", "9/9 colorless Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(9, 9)
	return tok
}

// ========================================
// NewTitaniaProtectorOfArgothElementalToken
// ========================================

// NewTitaniaProtectorOfArgothElementalToken creates a 5/3 green Elemental creature token.
func NewTitaniaProtectorOfArgothElementalToken() *Token {
	tok := NewToken("TitaniaProtectorOfArgothElementalToken", "5/3 green Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 3)
	return tok
}

// ========================================
// NewTolsimirMidnightsLightToken
// ========================================

// NewTolsimirMidnightsLightToken creates a Voja Fenstalker, a legendary 5/5 green and white Wolf creature token with trample.
func NewTolsimirMidnightsLightToken() *Token {
	tok := NewToken("TolsimirMidnightsLightToken", "Voja Fenstalker, a legendary 5/5 green and white Wolf creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewTombspawnZombieToken
// ========================================

// NewTombspawnZombieToken creates a 2/2 black Zombie creature token with haste named Tombspawn.
func NewTombspawnZombieToken() *Token {
	tok := NewToken("TombspawnZombieToken", "2/2 black Zombie creature token with haste named Tombspawn")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewToyToken
// ========================================

// NewToyToken creates a 1/1 white Toy artifact creature token.
func NewToyToken() *Token {
	tok := NewToken("ToyToken", "1/1 white Toy artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TOY")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTreasureToken
// ========================================

// NewTreasureToken creates a Treasure token.
func NewTreasureToken() *Token {
	tok := NewToken("TreasureToken", "Treasure token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("TREASURE")
	return tok
}

// ========================================
// NewTreefolkShamanToken
// ========================================

// NewTreefolkShamanToken creates a 2/5 green Treefolk Shaman creature token.
func NewTreefolkShamanToken() *Token {
	tok := NewToken("TreefolkShamanToken", "2/5 green Treefolk Shaman creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TREEFOLK")
	tok.AddSubtype("SHAMAN")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 5)
	return tok
}

// ========================================
// NewTriskelaviteToken
// ========================================

// NewTriskelaviteToken creates a token.
func NewTriskelaviteToken() *Token {
	tok := NewToken("TriskelaviteToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TRISKELAVITE")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewTrollWarriorToken
// ========================================

// NewTrollWarriorToken creates a 4/4 green Troll Warrior creature token with trample.
func NewTrollWarriorToken() *Token {
	tok := NewToken("TrollWarriorToken", "4/4 green Troll Warrior creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TROLL")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewTrooperToken
// ========================================

// NewTrooperToken creates a 1/1 white Trooper creature token.
func NewTrooperToken() *Token {
	tok := NewToken("TrooperToken", "1/1 white Trooper creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TROOPER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTuktukTheReturnedToken
// ========================================

// NewTuktukTheReturnedToken creates a Tuktuk the Returned, a legendary 5/5 colorless Goblin Golem artifact creature token.
func NewTuktukTheReturnedToken() *Token {
	tok := NewToken("TuktukTheReturnedToken", "Tuktuk the Returned, a legendary 5/5 colorless Goblin Golem artifact creature token")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("GOBLIN")
	tok.AddSubtype("GOLEM")
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewTuskenRaiderToken
// ========================================

// NewTuskenRaiderToken creates a 1/1 white Tusken Raider creature token.
func NewTuskenRaiderToken() *Token {
	tok := NewToken("TuskenRaiderToken", "1/1 white Tusken Raider creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TUSKEN")
	tok.AddSubtype("RAIDER")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTyranid55Token
// ========================================

// NewTyranid55Token creates a 5/5 green Tyranid creature token.
func NewTyranid55Token() *Token {
	tok := NewToken("Tyranid55Token", "5/5 green Tyranid creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TYRANID")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewTyranidGargoyleToken
// ========================================

// NewTyranidGargoyleToken creates a 1/1 blue Tyranid Gargoyle creature token with flying.
func NewTyranidGargoyleToken() *Token {
	tok := NewToken("TyranidGargoyleToken", "1/1 blue Tyranid Gargoyle creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TYRANID")
	tok.AddSubtype("GARGOYLE")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewTyranidToken
// ========================================

// NewTyranidToken creates a 1/1 green Tyranid creature token.
func NewTyranidToken() *Token {
	tok := NewToken("TyranidToken", "1/1 green Tyranid creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TYRANID")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewTyranidWarriorToken
// ========================================

// NewTyranidWarriorToken creates a 3/3 green Tyranid Warrior creature token with trample.
func NewTyranidWarriorToken() *Token {
	tok := NewToken("TyranidWarriorToken", "3/3 green Tyranid Warrior creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TYRANID")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewUginTheIneffableToken
// ========================================

// NewUginTheIneffableToken creates a 2/2 colorless Spirit creature token.
func NewUginTheIneffableToken() *Token {
	tok := NewToken("UginTheIneffableToken", "2/2 colorless Spirit creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewUktabiKongApeToken
// ========================================

// NewUktabiKongApeToken creates a 1/1 green Ape creature token.
func NewUktabiKongApeToken() *Token {
	tok := NewToken("UktabiKongApeToken", "1/1 green Ape creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("APE")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewUnicornToken
// ========================================

// NewUnicornToken creates a 2/2 white Unicorn creature token.
func NewUnicornToken() *Token {
	tok := NewToken("UnicornToken", "2/2 white Unicorn creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("UNICORN")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewUramiToken
// ========================================

// NewUramiToken creates a Urami, a legendary 5/5 black Demon Spirit creature token with flying.
func NewUramiToken() *Token {
	tok := NewToken("UramiToken", "Urami, a legendary 5/5 black Demon Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DEMON")
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewUtvaraHellkiteDragonToken
// ========================================

// NewUtvaraHellkiteDragonToken creates a 6/6 red Dragon creature token with flying.
func NewUtvaraHellkiteDragonToken() *Token {
	tok := NewToken("UtvaraHellkiteDragonToken", "6/6 red Dragon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(6, 6)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewValorsReachTagTeamToken
// ========================================

// NewValorsReachTagTeamToken creates a token.
func NewValorsReachTagTeamToken() *Token {
	tok := NewToken("ValorsReachTagTeamToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true, Red: true})
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewVampireDemonToken
// ========================================

// NewVampireDemonToken creates a 4/3 white and black Vampire Demon creature token with flying.
func NewVampireDemonToken() *Token {
	tok := NewToken("VampireDemonToken", "4/3 white and black Vampire Demon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.AddSubtype("DEMON")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(4, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewVampireKnightToken
// ========================================

// NewVampireKnightToken creates a 1/1 black Vampire Knight creature token with lifelink.
func NewVampireKnightToken() *Token {
	tok := NewToken("VampireKnightToken", "1/1 black Vampire Knight creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewVampireLifelinkToken
// ========================================

// NewVampireLifelinkToken creates a 2/3 black Vampire creature token with flying and lifelink.
func NewVampireLifelinkToken() *Token {
	tok := NewToken("VampireLifelinkToken", "2/3 black Vampire creature token with flying and lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 3)
	tok.AddAbility("flying")
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewVampireRogueToken
// ========================================

// NewVampireRogueToken creates a 1/1 black Vampire Rogue creature token with lifelink.
func NewVampireRogueToken() *Token {
	tok := NewToken("VampireRogueToken", "1/1 black Vampire Rogue creature token with lifelink")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewVampireToken
// ========================================

// NewVampireToken creates a 2/2 black Vampire creature token with flying.
func NewVampireToken() *Token {
	tok := NewToken("VampireToken", "2/2 black Vampire creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VAMPIRE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewVarmintToken
// ========================================

// NewVarmintToken creates a 2/1 green Varmint creature token.
func NewVarmintToken() *Token {
	tok := NewToken("VarmintToken", "2/1 green Varmint creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("VARMINT")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 1)
	return tok
}

// ========================================
// NewVecnaToken
// ========================================

// NewVecnaToken creates a Vecna, a legendary 8/8 black Zombie God creature token with indestructible and all triggered abilities of the exiled cards.
func NewVecnaToken() *Token {
	tok := NewToken("VecnaToken", "Vecna, a legendary 8/8 black Zombie God creature token with indestructible and all triggered abilities of the exiled cards")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("GOD")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(8, 8)
	tok.AddAbility("indestructible")
	return tok
}

// ========================================
// NewVehicleToken
// ========================================

// NewVehicleToken creates a 3/2 colorless Vehicle artifact token with crew 1.
func NewVehicleToken() *Token {
	tok := NewToken("VehicleToken", "3/2 colorless Vehicle artifact token with crew 1")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("VEHICLE")
	tok.SetPowerToughness(3, 2)
	return tok
}

// ========================================
// NewVirtuousRoleToken
// ========================================

// NewVirtuousRoleToken creates a Virtuous Role token.
func NewVirtuousRoleToken() *Token {
	tok := NewToken("VirtuousRoleToken", "Virtuous Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	return tok
}

// ========================================
// NewVoiceOfResurgenceToken
// ========================================

// NewVoiceOfResurgenceToken creates a token.
func NewVoiceOfResurgenceToken() *Token {
	tok := NewToken("VoiceOfResurgenceToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewVoiceOfTheWoodsElementalToken
// ========================================

// NewVoiceOfTheWoodsElementalToken creates a 7/7 green Elemental creature token with trample.
func NewVoiceOfTheWoodsElementalToken() *Token {
	tok := NewToken("VoiceOfTheWoodsElementalToken", "7/7 green Elemental creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(7, 7)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewVojaFriendToElvesToken
// ========================================

// NewVojaFriendToElvesToken creates a token.
func NewVojaFriendToElvesToken() *Token {
	tok := NewToken("VojaFriendToElvesToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewVojaToken
// ========================================

// NewVojaToken creates a Voja, a legendary 2/2 green and white Wolf creature token.
func NewVojaToken() *Token {
	tok := NewToken("VojaToken", "Voja, a legendary 2/2 green and white Wolf creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{White: true, Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewVolosJournalToken
// ========================================

// NewVolosJournalToken creates a token.
func NewVolosJournalToken() *Token {
	tok := NewToken("VolosJournalToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddAbility("hexproof")
	return tok
}

// ========================================
// NewVolrathsLaboratoryToken
// ========================================

// NewVolrathsLaboratoryToken creates a 2/2 creature token of the chosen color and type.
func NewVolrathsLaboratoryToken() *Token {
	tok := NewToken("VolrathsLaboratoryToken", "2/2 creature token of the chosen color and type")
	tok.AddCardType(CardTypeCreature)
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewVrenRatToken
// ========================================

// NewVrenRatToken creates a token.
func NewVrenRatToken() *Token {
	tok := NewToken("VrenRatToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("RAT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewVrondissRageOfAncientsToken
// ========================================

// NewVrondissRageOfAncientsToken creates a token.
func NewVrondissRageOfAncientsToken() *Token {
	tok := NewToken("VrondissRageOfAncientsToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DRAGON")
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(5, 4)
	return tok
}

// ========================================
// NewWalkerToken
// ========================================

// NewWalkerToken creates a Walker token.
func NewWalkerToken() *Token {
	tok := NewToken("WalkerToken", "Walker token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewWall13Token
// ========================================

// NewWall13Token creates a 1/3 white Wall creature token with defender.
func NewWall13Token() *Token {
	tok := NewToken("Wall13Token", "1/3 white Wall creature token with defender")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 3)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewWallFlyingToken
// ========================================

// NewWallFlyingToken creates a 0/4 white Wall creature token with defender and flying.
func NewWallFlyingToken() *Token {
	tok := NewToken("WallFlyingToken", "0/4 white Wall creature token with defender and flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 4)
	tok.AddAbility("flying")
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewWallOfResurgenceToken
// ========================================

// NewWallOfResurgenceToken creates a 0/0 Elemental creature with haste.
func NewWallOfResurgenceToken() *Token {
	tok := NewToken("WallOfResurgenceToken", "0/0 Elemental creature with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetPowerToughness(0, 0)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewWallToken
// ========================================

// NewWallToken creates a 0/2 colorless Wall artifact creature token with defender.
func NewWallToken() *Token {
	tok := NewToken("WallToken", "0/2 colorless Wall artifact creature token with defender")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetPowerToughness(0, 2)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewWallWhiteToken
// ========================================

// NewWallWhiteToken creates a 0/4 white Wall creature token with defender.
func NewWallWhiteToken() *Token {
	tok := NewToken("WallWhiteToken", "0/4 white Wall creature token with defender")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(0, 4)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewWandOfTheElementsFirstToken
// ========================================

// NewWandOfTheElementsFirstToken creates a 2/2 blue Elemental creature token with flying.
func NewWandOfTheElementsFirstToken() *Token {
	tok := NewToken("WandOfTheElementsFirstToken", "2/2 blue Elemental creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWandOfTheElementsSecondToken
// ========================================

// NewWandOfTheElementsSecondToken creates a 3/3 red Elemental creature token.
func NewWandOfTheElementsSecondToken() *Token {
	tok := NewToken("WandOfTheElementsSecondToken", "3/3 red Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewWardenSphinxToken
// ========================================

// NewWardenSphinxToken creates a 4/4 white and blue Sphinx creature token with flying and vigilance.
func NewWardenSphinxToken() *Token {
	tok := NewToken("WardenSphinxToken", "4/4 white and blue Sphinx creature token with flying and vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPHINX")
	tok.SetColor(Color{White: true, Blue: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewWarriorToken
// ========================================

// NewWarriorToken creates a 1/1 white Warrior creature token.
func NewWarriorToken() *Token {
	tok := NewToken("WarriorToken", "1/1 white Warrior creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewWarriorVigilantToken
// ========================================

// NewWarriorVigilantToken creates a 1/1 white Warrior creature token with vigilance.
func NewWarriorVigilantToken() *Token {
	tok := NewToken("WarriorVigilantToken", "1/1 white Warrior creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewWasitoraCatDragonToken
// ========================================

// NewWasitoraCatDragonToken creates a 3/3 black, red, and green Cat Dragon creature token with flying.
func NewWasitoraCatDragonToken() *Token {
	tok := NewToken("WasitoraCatDragonToken", "3/3 black, red, and green Cat Dragon creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("CAT")
	tok.AddSubtype("DRAGON")
	tok.SetColor(Color{Black: true, Red: true, Green: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWaspToken
// ========================================

// NewWaspToken creates a 1/1 colorless Insect artifact creature token with flying named Wasp.
func NewWaspToken() *Token {
	tok := NewToken("WaspToken", "1/1 colorless Insect artifact creature token with flying named Wasp")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWastelandSurvivalGuideToken
// ========================================

// NewWastelandSurvivalGuideToken creates a token.
func NewWastelandSurvivalGuideToken() *Token {
	tok := NewToken("WastelandSurvivalGuideToken", "")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("EQUIPMENT")
	return tok
}

// ========================================
// NewWaylayToken
// ========================================

// NewWaylayToken creates a 2/2 white Knight creature token.
func NewWaylayToken() *Token {
	tok := NewToken("WaylayToken", "2/2 white Knight creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("KNIGHT")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewWeirdToken
// ========================================

// NewWeirdToken creates a 3/3 blue Weird create token with defender and flying.
func NewWeirdToken() *Token {
	tok := NewToken("WeirdToken", "3/3 blue Weird create token with defender and flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WEIRD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("flying")
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewWhiteAstartesWarriorToken
// ========================================

// NewWhiteAstartesWarriorToken creates a 2/2 white Astartes Warrior creature token with vigilance.
func NewWhiteAstartesWarriorToken() *Token {
	tok := NewToken("WhiteAstartesWarriorToken", "2/2 white Astartes Warrior creature token with vigilance")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ASTARTES")
	tok.AddSubtype("WARRIOR")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("vigilance")
	return tok
}

// ========================================
// NewWhiteBlackSpiritToken
// ========================================

// NewWhiteBlackSpiritToken creates a 1/1 white and black Spirit creature token with flying.
func NewWhiteBlackSpiritToken() *Token {
	tok := NewToken("WhiteBlackSpiritToken", "1/1 white and black Spirit creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SPIRIT")
	tok.SetColor(Color{White: true, Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWhiteBlueBirdToken
// ========================================

// NewWhiteBlueBirdToken creates a 1/1 white and blue Bird creature token with flying.
func NewWhiteBlueBirdToken() *Token {
	tok := NewToken("WhiteBlueBirdToken", "1/1 white and blue Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true, Blue: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWhiteDogToken
// ========================================

// NewWhiteDogToken creates a 1/1 white Dog creature token.
func NewWhiteDogToken() *Token {
	tok := NewToken("WhiteDogToken", "1/1 white Dog creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("DOG")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewWhiteElementalToken
// ========================================

// NewWhiteElementalToken creates a 4/4 white Elemental creature token with flying.
func NewWhiteElementalToken() *Token {
	tok := NewToken("WhiteElementalToken", "4/4 white Elemental creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(4, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWickedRoleToken
// ========================================

// NewWickedRoleToken creates a Wicked Role token.
func NewWickedRoleToken() *Token {
	tok := NewToken("WickedRoleToken", "Wicked Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	return tok
}

// ========================================
// NewWildfireAwakenerToken
// ========================================

// NewWildfireAwakenerToken creates a token.
func NewWildfireAwakenerToken() *Token {
	tok := NewToken("WildfireAwakenerToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Red: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewWingmateRocToken
// ========================================

// NewWingmateRocToken creates a 3/4 white Bird creature token with flying.
func NewWingmateRocToken() *Token {
	tok := NewToken("WingmateRocToken", "3/4 white Bird creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BIRD")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(3, 4)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWireflyToken
// ========================================

// NewWireflyToken creates a 2/2 colorless Insect artifact creature token with flying named Wirefly.
func NewWireflyToken() *Token {
	tok := NewToken("WireflyToken", "2/2 colorless Insect artifact creature token with flying named Wirefly")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewWizardToken
// ========================================

// NewWizardToken creates a 2/2 blue Wizard creature token.
func NewWizardToken() *Token {
	tok := NewToken("WizardToken", "2/2 blue Wizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Blue: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewWolfToken
// ========================================

// NewWolfToken creates a 2/2 green Wolf creature token.
func NewWolfToken() *Token {
	tok := NewToken("WolfToken", "2/2 green Wolf creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewWolfsQuarryToken
// ========================================

// NewWolfsQuarryToken creates a token.
func NewWolfsQuarryToken() *Token {
	tok := NewToken("WolfsQuarryToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("BOAR")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewWolvesOfTheHuntToken
// ========================================

// NewWolvesOfTheHuntToken creates a token.
func NewWolvesOfTheHuntToken() *Token {
	tok := NewToken("WolvesOfTheHuntToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WOLF")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewWoodToken
// ========================================

// NewWoodToken creates a 0/1 green Wall creature token with defender named Wood.
func NewWoodToken() *Token {
	tok := NewToken("WoodToken", "0/1 green Wall creature token with defender named Wood")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WALL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 1)
	tok.AddAbility("defender")
	return tok
}

// ========================================
// NewWraithToken
// ========================================

// NewWraithToken creates a 3/3 black Wraith creature token with menace.
func NewWraithToken() *Token {
	tok := NewToken("WraithToken", "3/3 black Wraith creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WRAITH")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(3, 3)
	return tok
}

// ========================================
// NewWrennAndSevenTreefolkToken
// ========================================

// NewWrennAndSevenTreefolkToken creates a token.
func NewWrennAndSevenTreefolkToken() *Token {
	tok := NewToken("WrennAndSevenTreefolkToken", "")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("TREEFOLK")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	tok.AddAbility("reach")
	return tok
}

// ========================================
// NewWurm44Token
// ========================================

// NewWurm44Token creates a 4/4 green Wurm creature token.
func NewWurm44Token() *Token {
	tok := NewToken("Wurm44Token", "4/4 green Wurm creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(4, 4)
	return tok
}

// ========================================
// NewWurm55Token
// ========================================

// NewWurm55Token creates a 5/5 green Wurm creature token.
func NewWurm55Token() *Token {
	tok := NewToken("Wurm55Token", "5/5 green Wurm creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 5)
	return tok
}

// ========================================
// NewWurmCallingWurmToken
// ========================================

// NewWurmCallingWurmToken creates a X/X green Wurm creature token.
func NewWurmCallingWurmToken() *Token {
	tok := NewToken("WurmCallingWurmToken", "X/X green Wurm creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewWurmToken
// ========================================

// NewWurmToken creates a 6/6 green Wurm creature token.
func NewWurmToken() *Token {
	tok := NewToken("WurmToken", "6/6 green Wurm creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(6, 6)
	return tok
}

// ========================================
// NewWurmWithDeathtouchToken
// ========================================

// NewWurmWithDeathtouchToken creates a 3/3 colorless Phyrexian Wurm artifact creature token with deathtouch.
func NewWurmWithDeathtouchToken() *Token {
	tok := NewToken("WurmWithDeathtouchToken", "3/3 colorless Phyrexian Wurm artifact creature token with deathtouch")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("WURM")
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("deathtouch")
	return tok
}

// ========================================
// NewWurmWithLifelinkToken
// ========================================

// NewWurmWithLifelinkToken creates a 3/3 colorless Phyrexian Wurm artifact creature token with lifelink.
func NewWurmWithLifelinkToken() *Token {
	tok := NewToken("WurmWithLifelinkToken", "3/3 colorless Phyrexian Wurm artifact creature token with lifelink")
	tok.AddCardType(CardTypeArtifact)
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("PHYREXIAN")
	tok.AddSubtype("WURM")
	tok.SetPowerToughness(3, 3)
	tok.AddAbility("lifelink")
	return tok
}

// ========================================
// NewWurmWithTrampleToken
// ========================================

// NewWurmWithTrampleToken creates a 5/5 green Wurm creature token with trample.
func NewWurmWithTrampleToken() *Token {
	tok := NewToken("WurmWithTrampleToken", "5/5 green Wurm creature token with trample")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("WURM")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("trample")
	return tok
}

// ========================================
// NewXenagosSatyrToken
// ========================================

// NewXenagosSatyrToken creates a 2/2 red and green Satyr creature token with haste.
func NewXenagosSatyrToken() *Token {
	tok := NewToken("XenagosSatyrToken", "2/2 red and green Satyr creature token with haste")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("SATYR")
	tok.SetColor(Color{Red: true, Green: true})
	tok.SetPowerToughness(2, 2)
	tok.AddAbility("haste")
	return tok
}

// ========================================
// NewXiraBlackInsectToken
// ========================================

// NewXiraBlackInsectToken creates a 1/1 black Insect creature token with flying.
func NewXiraBlackInsectToken() *Token {
	tok := NewToken("XiraBlackInsectToken", "1/1 black Insect creature token with flying")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("INSECT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(1, 1)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewXmageToken
// ========================================

// NewXmageToken creates a token.
func NewXmageToken() *Token {
	tok := NewToken("XmageToken", "")
	return tok
}

// ========================================
// NewYoungHeroRoleToken
// ========================================

// NewYoungHeroRoleToken creates a Young Hero Role token.
func NewYoungHeroRoleToken() *Token {
	tok := NewToken("YoungHeroRoleToken", "Young Hero Role token")
	tok.AddCardType(CardTypeEnchantment)
	tok.AddSubtype("AURA")
	tok.AddSubtype("ROLE")
	return tok
}

// ========================================
// NewZaxaraTheExemplaryHydraToken
// ========================================

// NewZaxaraTheExemplaryHydraToken creates a 0/0 green Hydra creature token.
func NewZaxaraTheExemplaryHydraToken() *Token {
	tok := NewToken("ZaxaraTheExemplaryHydraToken", "0/0 green Hydra creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("HYDRA")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewZendikarsRoilElementalToken
// ========================================

// NewZendikarsRoilElementalToken creates a 2/2 green Elemental creature token.
func NewZendikarsRoilElementalToken() *Token {
	tok := NewToken("ZendikarsRoilElementalToken", "2/2 green Elemental creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ELEMENTAL")
	tok.SetColor(Color{Green: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZeppelinToken
// ========================================

// NewZeppelinToken creates a 5/5 colorless Vehicle artifact token named Zeppelin with flying and crew 3.
func NewZeppelinToken() *Token {
	tok := NewToken("ZeppelinToken", "5/5 colorless Vehicle artifact token named Zeppelin with flying and crew 3")
	tok.AddCardType(CardTypeArtifact)
	tok.AddSubtype("VEHICLE")
	tok.SetPowerToughness(5, 5)
	tok.AddAbility("flying")
	return tok
}

// ========================================
// NewZombieArmyToken
// ========================================

// NewZombieArmyToken creates a 0/0 black Zombie Army creature token.
func NewZombieArmyToken() *Token {
	tok := NewToken("ZombieArmyToken", "0/0 black Zombie Army creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("ARMY")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(0, 0)
	return tok
}

// ========================================
// NewZombieBerserkerToken
// ========================================

// NewZombieBerserkerToken creates a 2/2 black Zombie Berserker creature token.
func NewZombieBerserkerToken() *Token {
	tok := NewToken("ZombieBerserkerToken", "2/2 black Zombie Berserker creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("BERSERKER")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieDecayedToken
// ========================================

// NewZombieDecayedToken creates a 2/2 black Zombie creature token with decayed.
func NewZombieDecayedToken() *Token {
	tok := NewToken("ZombieDecayedToken", "2/2 black Zombie creature token with decayed")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieDruidToken
// ========================================

// NewZombieDruidToken creates a 2/2 black Zombie Druid creature token.
func NewZombieDruidToken() *Token {
	tok := NewToken("ZombieDruidToken", "2/2 black Zombie Druid creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("DRUID")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieKnightToken
// ========================================

// NewZombieKnightToken creates a 2/2 black Zombie Knight creature token with menace.
func NewZombieKnightToken() *Token {
	tok := NewToken("ZombieKnightToken", "2/2 black Zombie Knight creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieMenaceToken
// ========================================

// NewZombieMenaceToken creates a X/X blue and black Zombie creature token with menace.
func NewZombieMenaceToken() *Token {
	tok := NewToken("ZombieMenaceToken", "X/X blue and black Zombie creature token with menace")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Blue: true, Black: true})
	return tok
}

// ========================================
// NewZombieMutantToken
// ========================================

// NewZombieMutantToken creates a 2/2 black Zombie Mutant creature token.
func NewZombieMutantToken() *Token {
	tok := NewToken("ZombieMutantToken", "2/2 black Zombie Mutant creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("MUTANT")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieRogueToken
// ========================================

// NewZombieRogueToken creates a 2/2 blue and black Zombie Rogue creature token.
func NewZombieRogueToken() *Token {
	tok := NewToken("ZombieRogueToken", "2/2 blue and black Zombie Rogue creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("ROGUE")
	tok.SetColor(Color{Blue: true, Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieToken
// ========================================

// NewZombieToken creates a 2/2 black Zombie creature token.
func NewZombieToken() *Token {
	tok := NewToken("ZombieToken", "2/2 black Zombie creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{Black: true})
	tok.SetPowerToughness(2, 2)
	return tok
}

// ========================================
// NewZombieWhiteToken
// ========================================

// NewZombieWhiteToken creates a 1/1 white Zombie creature token.
func NewZombieWhiteToken() *Token {
	tok := NewToken("ZombieWhiteToken", "1/1 white Zombie creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.SetColor(Color{White: true})
	tok.SetPowerToughness(1, 1)
	return tok
}

// ========================================
// NewZombieWizardToken
// ========================================

// NewZombieWizardToken creates a 1/1 blue and black Zombie Wizard creature token.
func NewZombieWizardToken() *Token {
	tok := NewToken("ZombieWizardToken", "1/1 blue and black Zombie Wizard creature token")
	tok.AddCardType(CardTypeCreature)
	tok.AddSubtype("ZOMBIE")
	tok.AddSubtype("WIZARD")
	tok.SetColor(Color{Blue: true, Black: true})
	tok.SetPowerToughness(1, 1)
	return tok
}
