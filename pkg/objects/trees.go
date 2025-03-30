package objects

import "math/rand"

type Tree struct {
	X, Y   int
	Sprite []string
}

var TreeSprite = []string{
	"  ^  ",
	" /|\\ ",
	"  |  ",
}

var Trees []Tree

// InitTrees генерирует n деревьев в зоне около земли.
func InitTrees(n int) {
	Trees = make([]Tree, n)
	for i := 0; i < n; i++ {
		Trees[i] = Tree{
			X:      rand.Intn(WorldWidth),
			Y:      GroundLevel - len(TreeSprite), // чтобы дерево "стоило" на земле
			Sprite: TreeSprite,
		}
	}
}
