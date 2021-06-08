package gui

type HCT interface {
	GUIObject
	SetHead(ui GUIObject, size int32)
	SetCenter(ui GUIObject)
	SetTail(ui GUIObject, size int32)
}

type hct struct {
	vertical bool

	uiHead   GUIObject
	sizeHead int32
	uiCenter GUIObject
	uiTail   GUIObject
	sizeTail int32
}

func NewVerticalHCT() HCT   { return &hct{vertical: true} }
func NewHorizontalHCT() HCT { return &hct{vertical: false} }

func (hct *hct) Render() {
	if hct.uiHead != nil {
		hct.uiHead.Render()
	}
	if hct.uiCenter != nil {
		hct.uiCenter.Render()
	}
	if hct.uiTail != nil {
		hct.uiTail.Render()
	}
}

func (hct *hct) GetMouseTargets() []MouseTarget {
	var res []MouseTarget
	for _, obj := range append([]GUIObject{}, hct.uiHead, hct.uiCenter, hct.uiTail) {
		if obj != nil {
			res = append(res, obj.GetMouseTargets()...)
		}
	}
	return res
}

func (hct *hct) Resize(r Rect) {
	if r.H > 0 && r.W > 0 {
		if hct.vertical {
			sizeCenter := r.H - hct.sizeHead - hct.sizeTail
			if hct.uiHead != nil {
				hct.uiHead.Resize(r.Relative(Rect{0, r.H - hct.sizeHead, r.W, hct.sizeHead}))
			}
			if hct.uiCenter != nil {
				hct.uiCenter.Resize(r.Relative(Rect{0, hct.sizeTail, r.W, sizeCenter}))
			}
			if hct.uiTail != nil {
				hct.uiTail.Resize(r.Relative(Rect{0, 0, r.W, hct.sizeTail}))
			}
		} else {
			sizeCenter := r.W - hct.sizeHead - hct.sizeTail
			if hct.uiHead != nil {
				hct.uiHead.Resize(r.Relative(Rect{0, 0, hct.sizeHead, r.H}))
			}
			if hct.uiCenter != nil {
				hct.uiCenter.Resize(r.Relative(Rect{hct.sizeHead, 0, sizeCenter, r.H}))
			}
			if hct.uiTail != nil {
				hct.uiTail.Resize(r.Relative(Rect{r.W - hct.sizeTail, 0, hct.sizeTail, r.H}))
			}
		}
	}
}

func (hct *hct) SetHead(ui GUIObject, size int32) {
	hct.uiHead = ui
	hct.sizeHead = size
}

func (hct *hct) SetCenter(ui GUIObject) {
	hct.uiCenter = ui
}

func (hct *hct) SetTail(ui GUIObject, size int32) {
	hct.uiTail = ui
	hct.sizeTail = size
}
