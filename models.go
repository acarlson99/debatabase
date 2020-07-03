package main

type DocArticle struct {
	Name        string   `json:"name" maximum:"512" example:"google"`
	URL         string   `json:"url" maximum:"512" example:"google.com"`
	Description string   `json:"description" maximum:"1024" example:"a popular search engine"`
	Tags        []string `json:"tags" example:"engine,search,browser"`
}

type DocTag struct {
	Name        string `json:"name" maximum:"16" example:"engine"`
	Description string `json:"description" maximum:"256" example:"a machine designed to convert one form of energy into mechanical energy"`
}

type DocUser struct {
	Name   string `json:"name" example:"John"`
	Passwd string `json:"password" example:"password"`
}
