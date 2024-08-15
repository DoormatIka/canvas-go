package utils

import (
	"image/gif"
)

type FlatOctreeNode struct {
    Color        Color
    PixelCount   int
    PaletteIndex int
}
type FlatOctreeNodeIndex struct {
	Index 		 int
	Children     [8]*FlatOctreeNodeIndex
}
type FlatOctree struct {
	Root  FlatOctreeNodeIndex // a tree of ints pointing to FlatOctree.Nodes
	Leaves []int // ints pointing to FlatOctree.Nodes
	Nodes []FlatOctreeNode
    Levels map[int][]*FlatOctreeNodeIndex
}


func NewNodeIndex(index int) FlatOctreeNodeIndex {
	return FlatOctreeNodeIndex{Index: index, Children: [8]*FlatOctreeNodeIndex{}};
}

func NewFlatOctree() *FlatOctree {
    quantizer := &FlatOctree {
        Levels: make(map[int][]*FlatOctreeNodeIndex),
    }
	quantizer.Nodes = append(quantizer.Nodes, *NewFlatNode(0, quantizer));
	quantizer.Root = NewNodeIndex(len(quantizer.Nodes)-1);
    return quantizer
}

func NewFlatNode(level int, parent *FlatOctree) *FlatOctreeNode {
    node := FlatOctreeNode{
        Color: Color{0, 0, 0, 0},
    }
	parent.Nodes = append(parent.Nodes, node);
	n := NewNodeIndex(len(parent.Nodes)-1);
    if level < MaxDepth-1 {
        parent.AddLevelNode(level, &n)
    }
    return &node;
}

func (octree *FlatOctree) AddLevelNode(level int, node *FlatOctreeNodeIndex) {
    octree.Levels[level] = append(octree.Levels[level], node);
}

func (quantizer *FlatOctree) AddColor(color Color) {
    quantizer.Root.AddColor(color, 0, quantizer)
}
func (ind *FlatOctreeNodeIndex) AddColor(color Color, level int, parent *FlatOctree) {
	if level >= MaxDepth {
		node := parent.Nodes[ind.Index];

		node.Color.Red += color.Red
		node.Color.Green += color.Green
		node.Color.Blue += color.Blue
		node.Color.Alpha += color.Alpha
		node.PixelCount++
		return
	}
    index := GetColorIndexForLevel(color, level)

	// FlatOctreeNodeIndex is cheaper to deref than the FlatOctreeNode.
	// which gets called every time, instead of FlatOctreeNode being called every traversal.
	// i could probably skip the traversal entirely if i knew how the octree worked.
	// i could convert this into an O(1) operation.. :D
    if ind.Children[index] == nil {
		parent.Nodes = append(parent.Nodes, *NewFlatNode(level, parent));
		n := NewNodeIndex(len(parent.Nodes) - 1);
        ind.Children[index] = &n;
    }
	ind.Children[index].AddColor(color, level+1, parent);
}


func (quantizer *FlatOctree) GetPaletteIndex(color Color) int {
    return quantizer.Root.GetPaletteIndex(color, 0, quantizer);
}
func (ind *FlatOctreeNodeIndex) GetPaletteIndex(color Color, level int, parent *FlatOctree) int {
	node := parent.Nodes[ind.Index]; // since nodes is being stored as a value now, no ptr deref impact.
    if node.IsLeaf() {
        return node.PaletteIndex
    }
    index := GetColorIndexForLevel(color, level)
    if ind.Children[index] != nil {
        return ind.Children[index].GetPaletteIndex(color, level+1, parent);
    }
    for _, child := range ind.Children {
        if child != nil {
            return child.GetPaletteIndex(color, level+1, parent);
        }
    }
    return 0
}

func (node *FlatOctreeNode) IsLeaf() bool {
    return node.PixelCount > 0
}
func (ind *FlatOctreeNodeIndex) GetLeafNodes(parent *FlatOctree) []int {
    var leafNodes []int
    for _, child := range ind.Children {
        if child != nil {
			node := parent.Nodes[ind.Index];
            if node.IsLeaf() {
                leafNodes = append(leafNodes, child.Index)
            } else {
                leafNodes = append(leafNodes, child.GetLeafNodes(parent)...)
            }
        }
    }
    return leafNodes
}
func (quantizer *FlatOctree) GetLeaves() []int {
    return quantizer.Root.GetLeafNodes(quantizer);
}
func (ind *FlatOctreeNodeIndex) RemoveLeaves(parent *FlatOctree) int {
    result := 0
    for i := range ind.Children {
        if ind.Children[i] != nil {
			node := parent.Nodes[ind.Index];
			child := parent.Nodes[ind.Children[i].Index];
            node.Color.Red += child.Color.Red
            node.Color.Green += child.Color.Green
            node.Color.Blue += child.Color.Blue
            node.Color.Alpha += child.Color.Alpha
            node.PixelCount += child.PixelCount
            result++
        }
    }
    return result - 1
}

func (quantizer *FlatOctree) MakePalette(colorCount int) []Color {
    var palette []Color
    paletteIndex := 0
    leafCount := len(quantizer.GetLeaves())
    for level := MaxDepth - 1; level >= 0; level-- {
        if nodes, exists := quantizer.Levels[level]; exists {
            for _, node := range nodes {
                leafCount -= node.RemoveLeaves(quantizer);
                if leafCount <= colorCount {
                    break
                }
            }
            if leafCount <= colorCount {
                break
            }
            quantizer.Levels[level] = nil
        }
    }
    for _, ind := range quantizer.GetLeaves() {
		node := quantizer.Nodes[ind];
        if paletteIndex >= colorCount {
            break
        }
        if node.IsLeaf() {
            palette = append(palette, node.GetColor());
            node.PaletteIndex = paletteIndex
            paletteIndex++
        }
    }
    return palette
}

func AddColorsToFlatOctree(q *FlatOctree, g *gif.GIF) {
    // Add colors from each frame to the quantizer
    for _, frame := range g.Image {
        bounds := frame.Bounds()
        for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                r, g, b, a := frame.At(x, y).RGBA()
                color := NewColor(int(r>>8), int(g>>8), int(b>>8), int(a>>8))
                q.AddColor(color) // called every pixel in every frame!
            }
        }
    }
}

// heres to hoping this function doesn't get compiled to the FlatOctreeNode.
func GetColorIndexForLevel(color Color, level int) int {
    index := 0
    mask := 0x80 >> level
    if color.Red&mask != 0 {
        index |= 4
    }
    if color.Green&mask != 0 {
        index |= 2
    }
    if color.Blue&mask != 0 {
        index |= 1
    }
    return index
}

func (node *FlatOctreeNode) GetColor() Color {
    if node.PixelCount == 0 {
        return Color{0, 0, 0, 0}
    }
    return Color{
        Red:   node.Color.Red / node.PixelCount,
        Green: node.Color.Green / node.PixelCount,
        Blue:  node.Color.Blue / node.PixelCount,
        Alpha: node.Color.Alpha / node.PixelCount,
    }
}

