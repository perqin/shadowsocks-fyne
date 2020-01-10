package main

type Profile struct {
	Name       string `json:"name"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Acl        string `json:"acl"`
}
