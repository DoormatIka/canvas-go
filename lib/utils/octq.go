package utils

import (
	"image/color"
)

type OctreeNode struct {
    Color        color.RGBA
    PixelCount   int
    PaletteIndex int
    Children     [8]*OctreeNode
}

const MaxDepth = 8

type OctreeQuantizer struct {
    Levels map[int][]*OctreeNode
    Root   *OctreeNode
}

func NewColor(red, green, blue, alpha uint8) color.RGBA {
    return color.RGBA{R: red, G: green, B: blue, A: alpha}
}

func NewOctreeNode(level int, parent *OctreeQuantizer) *OctreeNode {
    node := &OctreeNode{
        Color: color.RGBA{0, 0, 0, 0},
    }
    if level < MaxDepth-1 {
        parent.AddLevelNode(level, node)
    }
    return node
}

func (node *OctreeNode) IsLeaf() bool {
    return node.PixelCount > 0
}

func (node *OctreeNode) GetLeafNodes() []*OctreeNode {
    var leafNodes []*OctreeNode
    for _, child := range node.Children {
        if child != nil {
            if child.IsLeaf() {
                leafNodes = append(leafNodes, child)
            } else {
                leafNodes = append(leafNodes, child.GetLeafNodes()...)
            }
        }
    }
    return leafNodes
}

func (node *OctreeNode) GetNodesPixelCount() int {
    sumCount := node.PixelCount
    for _, child := range node.Children {
        if child != nil {
            sumCount += child.PixelCount
        }
    }
    return sumCount
}

func (node *OctreeNode) AddColor(color color.RGBA, level int, parent *OctreeQuantizer) {
    if level >= MaxDepth {
        node.Color.R += color.R
        node.Color.G += color.G
        node.Color.B += color.B
		node.Color.A += color.A
        node.PixelCount++
        return
    }
    index := node.GetColorIndexForLevel(color, level)
    if node.Children[index] == nil {
        node.Children[index] = NewOctreeNode(level, parent)
    }
    node.Children[index].AddColor(color, level+1, parent)
}

func (node *OctreeNode) GetPaletteIndex(color color.RGBA, level int) int {
    if node.IsLeaf() {
        return node.PaletteIndex
    }
    index := node.GetColorIndexForLevel(color, level)
    if node.Children[index] != nil {
        return node.Children[index].GetPaletteIndex(color, level+1)
    }
    for _, child := range node.Children {
        if child != nil {
            return child.GetPaletteIndex(color, level+1)
        }
    }
    return 0
}

func (node *OctreeNode) RemoveLeaves() int {
    result := 0
    for i := range node.Children {
        if node.Children[i] != nil {
            node.Color.R += node.Children[i].Color.R
            node.Color.G += node.Children[i].Color.G
            node.Color.B += node.Children[i].Color.B
			node.Color.A += node.Children[i].Color.A
            node.PixelCount += node.Children[i].PixelCount
            result++
        }
    }
    return result - 1
}

func (node *OctreeNode) GetColorIndexForLevel(color color.RGBA, level int) int {
    index := 0
    mask := 0x80 >> level
    if int(color.R)&mask != 0 {
        index |= 4
    }
    if int(color.G)&mask != 0 {
        index |= 2
    }
    if int(color.B)&mask != 0 {
        index |= 1
    }
    return index
}

func (node *OctreeNode) GetColor() color.RGBA {
    if node.PixelCount == 0 {
        return color.RGBA{0, 0, 0, 0}
    }
    return color.RGBA{
        R: uint8(int(node.Color.R) / node.PixelCount),
        G: uint8(int(node.Color.G) / node.PixelCount),
        B: uint8(int(node.Color.B) / node.PixelCount),
		A: uint8(int(node.Color.A) / node.PixelCount),
    }
}

func NewOctreeQuantizer() *OctreeQuantizer {
    quantizer := &OctreeQuantizer{
        Levels: make(map[int][]*OctreeNode),
    }
    quantizer.Root = NewOctreeNode(0, quantizer)
    return quantizer
}

func (quantizer *OctreeQuantizer) GetLeaves() []*OctreeNode {
    return quantizer.Root.GetLeafNodes()
}

func (quantizer *OctreeQuantizer) AddLevelNode(level int, node *OctreeNode) {
    quantizer.Levels[level] = append(quantizer.Levels[level], node)
}

func (quantizer *OctreeQuantizer) AddColor(color color.RGBA) {
    quantizer.Root.AddColor(color, 0, quantizer)
}

func (quantizer *OctreeQuantizer) MakePalette(colorCount int) []color.RGBA {
    var palette []color.RGBA
    paletteIndex := 0
    leafCount := len(quantizer.GetLeaves())
    for level := MaxDepth - 1; level >= 0; level-- {
        if nodes, exists := quantizer.Levels[level]; exists {
            for _, node := range nodes {
                leafCount -= node.RemoveLeaves()
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
    for _, node := range quantizer.GetLeaves() {
        if paletteIndex >= colorCount {
            break
        }
        if node.IsLeaf() {
            palette = append(palette, node.GetColor())
            node.PaletteIndex = paletteIndex
            paletteIndex++
        }
    }
    return palette
}

func (quantizer *OctreeQuantizer) GetPaletteIndex(color color.RGBA) int {
    return quantizer.Root.GetPaletteIndex(color, 0)
}

