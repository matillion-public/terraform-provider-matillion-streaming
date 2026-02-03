package models

// PostgresSourceModel represents a PostgreSQL source configuration in the provider schema
type PostgresSourceModel struct {
	Connection PostgresConnectionModel `tfsdk:"connection"`
	Tables     []TableResourceModel    `tfsdk:"tables"`
}

// SqlServerSourceModel represents a SQL Server source configuration in the provider schema
type SqlServerSourceModel struct {
	Connection SqlServerConnectionModel `tfsdk:"connection"`
	Tables     []TableResourceModel     `tfsdk:"tables"`
}

// MysqlSourceModel represents a MySQL source configuration in the provider schema
type MysqlSourceModel struct {
	Connection MysqlConnectionModel `tfsdk:"connection"`
	Tables     []TableResourceModel `tfsdk:"tables"`
}

// OracleSourceModel represents an Oracle source configuration in the provider schema
type OracleSourceModel struct {
	Connection OracleConnectionModel `tfsdk:"connection"`
	Tables     []TableResourceModel  `tfsdk:"tables"`
}

// Db2IbmISourceModel represents a DB2 for IBM i source configuration in the provider schema
type Db2IbmISourceModel struct {
	Connection Db2IbmIConnectionModel `tfsdk:"connection"`
	Tables     []TableResourceModel   `tfsdk:"tables"`
}
