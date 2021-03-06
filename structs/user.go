package structs

import (
	"CaribbeanWarServer/rtree"
	"github.com/gorilla/websocket"
	"math"
	"sync"
)

const (
	left               = -1
	right              = 1
	none               = 0
	angleSpeed float64 = 0.075
	velocity   float64 = 0.01
)

type NearestUser struct {
	ID            *uint           `json:"id"`
	Conn          *websocket.Conn `json:""`
	Ship          *Ship           `json:"ship"`
	Nick          *string         `json:"nick"`
	Location      *Point          `json:"location"`
	RotationAngle float64         `json:"alpha"`
}

type User struct {
	ID                uint          `json:"id"`
	Email             string        `json:"email"`
	Nick              string        `json:"nick"`
	Cash              uint          `json:"cash"`
	Location          *Point        `json:"location"`
	Ships             []Ship        `json:"ships"`
	SelectedShip      *Ship         `json:"selectedShip"`
	NearestUsers      []NearestUser `json:"nearestUsers"`
	RotationAngle     float64       `json:"alpha"`
	conn              *websocket.Conn
	inWorld           bool
	sailsMode         int16
	speedRatio        float64
	rotationDirection int16
	sync.Mutex
}

func (self *User) Bounds() *rtree.Rect {
	bound, _ := rtree.NewRect(rtree.Point{self.Location.X, self.Location.Y}, []float64{1, 1})
	return bound
}

func (self *User) SetMove(moveType string) {
	self.Lock()
	switch moveType {
	case "upward":
		self.sailsMode = int16(math.Min(float64(self.sailsMode+1), 3))
	case "backward":
		self.sailsMode = int16(math.Max(float64(self.sailsMode-1), 0))
	case "left":
		self.rotationDirection = left
	case "right":
		self.rotationDirection = right
	case "none":
		self.rotationDirection = none
	default:
		self.GetConn().WriteJSON(map[string]string{
			"action":  "error",
			"details": "ERRORS_UNKNOWN_ACTION",
		})
	}
	self.Unlock()
}

func (self *User) UpdatePosition(delta float64) {
	ship := self.SelectedShip
	if ship != nil {
		self.speedRatio = lerp(self.speedRatio, float64(self.sailsMode)*ship.Speed*delta/4.0, velocity)
		self.Location.X += self.speedRatio * math.Cos(self.RotationAngle)
		self.Location.Y += self.speedRatio * math.Sin(-self.RotationAngle)
		self.RotationAngle = math.Mod(self.RotationAngle+(float64(self.rotationDirection)*angleSpeed*self.speedRatio)/(float64(self.sailsMode)+1.0), 2*math.Pi)
	}
}

func lerp(start, end, delta float64) float64 {
	return start + delta*(end-start)
}

func (self *User) GetConn() *websocket.Conn {
	return self.conn
}

func (self *User) SetConn(conn *websocket.Conn) {
	self.conn = conn
}

func (self *User) IsInWorld() bool {
	return self.inWorld
}

func (self *User) SetIsInWorld(is bool) {
	self.inWorld = is
}
