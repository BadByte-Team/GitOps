package models
type Credentials struct { Username string `json:"username"`; Password string `json:"password"` }
type Episode struct { ID int `json:"id"`; ModuleID int `json:"module_id"`; Title string `json:"title"`; VideoURL string `json:"video_url"`; IsHidden bool `json:"is_hidden"` }
type Module struct { ID int `json:"id"`; Title string `json:"title"`; IsHidden bool `json:"is_hidden"`; Episodes []Episode `json:"episodes"` }
