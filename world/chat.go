package world

func (self *storage) chat(data map[string]interface{}) {
	message := map[string]interface{}{
		"action": "chat",
		"details": map[string]interface{}{
			"sender":  data["sender"],
			"message": data["message"],
		},
	}
	self.Lock()
	defer self.Unlock()
	for _, v := range self.userList {
		v.Conn.WriteJSON(message)
	}
}
