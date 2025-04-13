package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
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
		field.Int("total").
			NonNegative(),
		field.Int("points").
			Default(0),
	}
}

// Edges of the Receipt.
func (Receipt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Item.Type),
	}
}

// Annotations of the Receipt.
func (Receipt) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "receipts"},
	}
}
