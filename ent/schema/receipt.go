package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Receipt holds the schema definition for the Receipt entity.
type Receipt struct {
	ent.Schema
}

// Fields of the Receipt.
func (Receipt) Fields() []ent.Field {
	return []ent.Field{
		field.String("retailer").
			NotEmpty(),
		field.String("purchase_date").
			NotEmpty(),
		field.String("purchase_time").
			NotEmpty(),
		field.String("total").
			NotEmpty(),
	}
}

// Edges of the Receipt.
func (Receipt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Item.Type),
	}
}
