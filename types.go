package main

type message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type lookupResult struct {
	Ok   bool `json:"ok"`
	User struct {
		ID string `json:"id"`
	} `json:"user"`
}

type sendResult struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type usersList struct {
	Ok      bool   `json:"ok"`
	Error   string `json:"error"`
	Members []struct {
		Name    string `json:"name"`
		Deleted bool   `json:"deleted"`
		IsBot   bool   `json:"is_bot"`
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	} `json:"members"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}
