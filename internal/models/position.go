package models

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

type Position = data.Vec2[float64]
type ConnectionIndexes = data.Vec2[int]

type Positions []Position

func (this *Positions) Add(position Position) { *this = append(*this, position) }
