package types

type Message struct {
	Body     string `json:"body"`
	RoomID   string `json:"roomID"`
	Username string `json:"username"`
	UserID   string `json:"userID"`
}

type CreateRoomReq struct {
	ID string `json:"id" validate:"required"`
}

type RoomRes struct {
	ID string `json:"id"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
