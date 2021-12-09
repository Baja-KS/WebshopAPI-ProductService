package database

import (
	"gorm.io/gorm"
	"os"
)

type Group struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name string `gorm:"not null;unique" json:"name"`
	Description string `gorm:"" json:"description,omitempty"`
}

func (g *Group) Out() GroupOut {
	return GroupOut{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
	}
}

func (g *Group) OutWithCategories() GroupOutWithCategories {
	var categories []CategoryOut
	categories,err:=GetCategories(g.ID,os.Getenv("CATEGORY_SERVICE"))
	if err != nil {
		return GroupOutWithCategories{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Categories:  nil,
		}
	}
	return GroupOutWithCategories{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Categories:  categories,
	}
}

func (g *GroupIn) In() Group {
	return Group{
		Name:        g.Name,
		Description: g.Description,
	}
}

func GroupArrayOut(groupModels []Group) []GroupOut {
	outArr:=make([]GroupOut,len(groupModels))
	for i,group := range groupModels {
		outArr[i]=group.Out()
	}
	return outArr
}

func GroupArrayOutWithCategories(groupModels []Group) []GroupOutWithCategories {
	outArr:=make([]GroupOutWithCategories,len(groupModels))
	for i, group := range groupModels {
		outArr[i]=group.OutWithCategories()
	}
	return outArr
}

type GroupIn struct {
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
}

type GroupOut struct {
	ID uint `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
}

type GroupOutWithCategories struct {
	ID uint `json:"id,omitempty"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Categories []CategoryOut `json:"Categories"`
}



