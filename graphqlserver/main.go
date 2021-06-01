package main

import (
	"encoding/json"
	"fmt"
	"github.com/asura/go-graphql-tutorial/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	tutorials := models.Populate()

	var commentType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"body": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	var authorType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Tutorials": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})
	var tutorilType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Tutorial",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: authorType,
			},
			"comments": &graphql.Field{
				Type: graphql.NewList(commentType),
			},
		},
	})
	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type:        tutorilType,
				Description: "Create a new Tutorial",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					tutorial := models.Tutorial{
						Title: p.Args["title"].(string),
					}
					tutorials = append(tutorials, tutorial)
					return tutorial, nil
				},
			},
		},
	})

	//schema
	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorilType,
			Description: "Get Tutorial By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.Id) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorilType),
			Description: "Get tutorial list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				dsn := "root:123456@tcp(127.0.0.1:3306)/gozore?charset=utf8mb4&parseTime=True&loc=Local"
				db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
				if err != nil {
					log.Fatalf("gorm.open mysql error: %v", err)
				}
				var tutors []models.Tutor
				err = db.Model(&models.Tutor{}).Find(&tutors).Error
				if err != nil {
					log.Fatalf("db query tutorial error: %v", err)
				}
				return tutors, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error : %v", err)
	}

	/*query := `
	{
	list {
		id
		title
		comments {
			body
		}
		author {
			Name
			Tutorials
		}
	}
	}
	`*/
	/*query := `
		{
			tutorial(id:1) {
				title
				author {
					Name
					Tutorials
				}
			}
		}
	`*/
	query := `
	mutation {
	create(title: "Hello World") {
		title
		}
	}
`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJson, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJson)

	query = `
		{
			list {
				id
				title
			}
		}
		`
	params = graphql.Params{Schema: schema, RequestString: query}
	r = graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors; %+v", r.Errors)
	}
	rJson, _ = json.Marshal(r)
	fmt.Printf("%s \n", rJson)
}
