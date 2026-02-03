package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DatabaseSourceSchema returns a reusable schema for database sources
func DatabaseSourceSchema(dbType string, connectionSchema schema.SingleNestedAttribute) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: dbType + " source configuration",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"connection": connectionSchema,
			"tables":     TablesSchema(),
		},
	}
}

// TablesSchema returns a reusable schema for database tables
func TablesSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Required:   true,
		Validators: []validator.List{listvalidator.SizeAtLeast(1)},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"schema": schema.StringAttribute{
					Required: true,
				},
				"table": schema.StringAttribute{
					Required: true,
				},
			},
		},
	}
}

// PostgresSourceSchema returns the complete Postgres source schema
func PostgresSourceSchema() schema.SingleNestedAttribute {
	return DatabaseSourceSchema("Postgres", PostgresConnectionSchema())
}

// PostgresConnectionSchema returns schema for PostgreSQL connections
func PostgresConnectionSchema() schema.SingleNestedAttribute {
	base := BasicDatabaseConnectionSchema("Postgres")
	base.Attributes["database"] = schema.StringAttribute{
		MarkdownDescription: "Postgres database name",
		Required:            true,
	}
	return base
}

// SQLServerSourceSchema returns the complete SQLServer source schema
func SQLServerSourceSchema() schema.SingleNestedAttribute {
	return DatabaseSourceSchema("SQL Server", SQLServerConnectionSchema())
}

// SQLServerConnectionSchema returns schema for SQL Server connections
func SQLServerConnectionSchema() schema.SingleNestedAttribute {
	base := BasicDatabaseConnectionSchema("SQL Server")
	base.Attributes["database"] = schema.StringAttribute{
		MarkdownDescription: "SQL Server database name",
		Required:            true,
	}
	return base
}

// MySQLSourceSchema returns schema for MySQL connections (basic connection without database field)
func MySQLSourceSchema() schema.SingleNestedAttribute {
	return DatabaseSourceSchema("MySQL", BasicDatabaseConnectionSchema("MySQL"))
}

func OracleSourceSchema() schema.SingleNestedAttribute {
	return DatabaseSourceSchema("Oracle", OracleConnectionSchema())
}

// OracleConnectionSchema returns schema for Oracle connections (includes database and pdb fields)
func OracleConnectionSchema() schema.SingleNestedAttribute {
	base := BasicDatabaseConnectionSchema("Oracle")
	base.Attributes["database"] = schema.StringAttribute{
		MarkdownDescription: "Oracle database name",
		Required:            true,
	}
	base.Attributes["pdb"] = schema.StringAttribute{
		MarkdownDescription: "Oracle pluggable database (PDB)",
		Optional:            true,
	}
	return base
}

func DB2IbmISourceSchema() schema.SingleNestedAttribute {
	return DatabaseSourceSchema("DB2 for IBM i", BasicDatabaseConnectionSchema("DB2 for IBM i"))
}
