package gui

type Grid interface {
	GUIObject
	Add(...GUIObject)
}

type grid struct {
	objects  []GUIObject
	sections int32
	fix      int32
}

func NewHGrid(cols, rowH uint32) Grid {
	g := &grid{
		sections: int32(cols),
		fix:      int32(rowH),
	}
	return g
}

func (g *grid) Add(obj ...GUIObject) {
	g.objects = append(g.objects, obj...)
}

func (g *grid) Render() {
	for _, o := range g.objects {
		o.Render()
	}
}

func (g *grid) Resize(r Rect) {
	w := r.W / g.sections
	for idx, obj := range g.objects {
		row := int32(idx) / g.sections
		col := int32(idx) % g.sections
		obj.Resize(Rect{w * col, g.fix * row, w, g.fix})
	}
}
