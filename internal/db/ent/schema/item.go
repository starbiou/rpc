package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("short_description").
			NotEmpty(),
		field.Int("price").
			NonNegative(),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("receipts", Receipt.Type).
			Ref("items").
			Unique(),
	}
}

// Annotations of the Item.
func (Item) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "items"},
	}
}
