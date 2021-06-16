package gui

type Grid interface {
	GUIObject
	Add(...GUIObject)
}

type grid struct {
	objects  []GUIObject
	sections int32
	fix      int32
	inset    int32
}

func NewHGrid(cols, rowH, inset int32) Grid {
	g := &grid{
		sections: cols,
		fix:      rowH,
		inset:    inset,
	}
	return g
}

func (g *grid) Add(obj ...GUIObject) {
	g.objects = append(g.objects, obj...)
}

func (g *grid) GetChildrens() []GUIObject {
	return g.objects
}

func (g *grid) Render() {
	for _, o := range g.objects {
		o.Render()
	}
}

func (g *grid) Resize(r Rect) {
	if r.W == 0 || r.H == 0 {
		return
	}

	w := (r.W / g.sections)
	for idx, obj := range g.objects {
		row := int32(idx) / g.sections
		col := int32(idx) % g.sections
		rec := r.Relative(Rect{w*col + g.inset, r.H - g.fix - g.fix*row + g.inset, w - g.inset*2, g.fix - g.inset*2})
		obj.Resize(rec)
	}
}
