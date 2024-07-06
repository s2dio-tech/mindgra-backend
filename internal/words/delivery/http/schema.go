package http

type WordCreateRequestSchema struct {
	Content     string    `json:"content" validate:"required,max=50"`
	Description *string   `json:"description" validate:"omitempty,max=512"`
	Refs        *[]string `json:"refs" validate:"omitempty,dive,required"`
	GraphId     string    `json:"graphId" validate:"required,max=50"`
	SourceId    *string   `json:"sourceId" validate:"omitempty,max=50"`
}

type WordUpdateRequestSchema struct {
	Content     string    `json:"content" validate:"required,max=50"`
	Description *string   `json:"description" validate:"omitempty,max=512"`
	Refs        *[]string `json:"refs" validate:"omitempty,dive,required"`
}

type Link2WordsRequestSchema struct {
	SourceId string `json:"sourceId" validate:"required,max=50"`
	TargetId string `json:"targetId" validate:"required,max=50"`
}

type LinkCreateRequestSchema struct {
	Word1Id     string    `json:"word1Id" validate:"required,max=50"`
	Word2Id     string    `json:"word2Id" validate:"required,max=50"`
	Content     string    `json:"content" validate:"required,max=50"`
	Description *string   `json:"description" validate:"omitempty,max=512"`
	Refs        *[]string `json:"refs" validate:"omitempty,dive,required"`
}

type LinkUpdateRequestSchema struct {
	Content     string    `json:"content" validate:"required,max=50"`
	Description *string   `json:"description" validate:"omitempty,max=512"`
	Refs        *[]string `json:"refs" validate:"omitempty,dive,required"`
}

type LinkRemoveRequestSchema struct {
	Word1Id string `json:"word1Id" validate:"required,max=50"`
	Word2Id string `json:"word2Id" validate:"required,max=50"`
}

type GraphCreateRequestSchema struct {
	Name string `json:"name" validate:"required,max=128"`
}

type GraphUpdateRequestSchema struct {
	Name string `json:"name" validate:"max=128"`
	Type string `json:"type" validate:"max=20"`
}
