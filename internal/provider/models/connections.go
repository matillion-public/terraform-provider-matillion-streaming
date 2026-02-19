package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// DatabaseConnection provides a common interface for all database connection models
// This allows generic mapping of connection fields from API responses
type DatabaseConnection interface {
	SetHost(types.String)
	SetPort(types.Int32)
	SetDatabase(types.String) // no-op for MySQL/DB2 which don't have a database field
	SetUsername(types.String)
	GetPasswordRef() *SecretReferenceModel
	GetJdbcPropertiesRef() *map[string]types.String
}

// PostgresConnectionModel represents PostgreSQL connection parameters in the provider schema
type PostgresConnectionModel struct {
	Host           types.String            `tfsdk:"host"`
	Port           types.Int32             `tfsdk:"port"`
	Database       types.String            `tfsdk:"database"`
	Username       types.String            `tfsdk:"username"`
	Password       SecretReferenceModel    `tfsdk:"password"`
	JdbcProperties map[string]types.String `tfsdk:"jdbc_properties"`
}

// DatabaseConnection interface implementation for PostgresConnectionModel
func (c *PostgresConnectionModel) SetHost(v types.String)     { c.Host = v }
func (c *PostgresConnectionModel) SetPort(v types.Int32)      { c.Port = v }
func (c *PostgresConnectionModel) SetDatabase(v types.String) { c.Database = v }
func (c *PostgresConnectionModel) SetUsername(v types.String) { c.Username = v }
func (c *PostgresConnectionModel) GetPasswordRef() *SecretReferenceModel {
	return &c.Password
}
func (c *PostgresConnectionModel) GetJdbcPropertiesRef() *map[string]types.String {
	return &c.JdbcProperties
}

// SqlServerConnectionModel represents SQL Server connection parameters in the provider schema
type SqlServerConnectionModel struct {
	Host           types.String            `tfsdk:"host"`
	Port           types.Int32             `tfsdk:"port"`
	Database       types.String            `tfsdk:"database"`
	Username       types.String            `tfsdk:"username"`
	Password       SecretReferenceModel    `tfsdk:"password"`
	JdbcProperties map[string]types.String `tfsdk:"jdbc_properties"`
}

// DatabaseConnection interface implementation for SqlServerConnectionModel
func (c *SqlServerConnectionModel) SetHost(v types.String)     { c.Host = v }
func (c *SqlServerConnectionModel) SetPort(v types.Int32)      { c.Port = v }
func (c *SqlServerConnectionModel) SetDatabase(v types.String) { c.Database = v }
func (c *SqlServerConnectionModel) SetUsername(v types.String) { c.Username = v }
func (c *SqlServerConnectionModel) GetPasswordRef() *SecretReferenceModel {
	return &c.Password
}
func (c *SqlServerConnectionModel) GetJdbcPropertiesRef() *map[string]types.String {
	return &c.JdbcProperties
}

// MysqlConnectionModel represents MySQL connection parameters in the provider schema
type MysqlConnectionModel struct {
	Host           types.String            `tfsdk:"host"`
	Port           types.Int32             `tfsdk:"port"`
	Username       types.String            `tfsdk:"username"`
	Password       SecretReferenceModel    `tfsdk:"password"`
	JdbcProperties map[string]types.String `tfsdk:"jdbc_properties"`
}

// DatabaseConnection interface implementation for MysqlConnectionModel
func (c *MysqlConnectionModel) SetHost(v types.String)     { c.Host = v }
func (c *MysqlConnectionModel) SetPort(v types.Int32)      { c.Port = v }
func (c *MysqlConnectionModel) SetDatabase(v types.String) {} // no-op - MySQL doesn't have a database field
func (c *MysqlConnectionModel) SetUsername(v types.String) { c.Username = v }
func (c *MysqlConnectionModel) GetPasswordRef() *SecretReferenceModel {
	return &c.Password
}
func (c *MysqlConnectionModel) GetJdbcPropertiesRef() *map[string]types.String {
	return &c.JdbcProperties
}

// OracleConnectionModel represents Oracle connection parameters in the provider schema
type OracleConnectionModel struct {
	Host           types.String            `tfsdk:"host"`
	Port           types.Int32             `tfsdk:"port"`
	Database       types.String            `tfsdk:"database"`
	Pdb            types.String            `tfsdk:"pdb"`
	Username       types.String            `tfsdk:"username"`
	Password       SecretReferenceModel    `tfsdk:"password"`
	JdbcProperties map[string]types.String `tfsdk:"jdbc_properties"`
}

// DatabaseConnection interface implementation for OracleConnectionModel
func (c *OracleConnectionModel) SetHost(v types.String)     { c.Host = v }
func (c *OracleConnectionModel) SetPort(v types.Int32)      { c.Port = v }
func (c *OracleConnectionModel) SetDatabase(v types.String) { c.Database = v }
func (c *OracleConnectionModel) SetUsername(v types.String) { c.Username = v }
func (c *OracleConnectionModel) GetPasswordRef() *SecretReferenceModel {
	return &c.Password
}
func (c *OracleConnectionModel) GetJdbcPropertiesRef() *map[string]types.String {
	return &c.JdbcProperties
}

// Db2IbmIConnectionModel represents DB2 for IBM i connection parameters in the provider schema
type Db2IbmIConnectionModel struct {
	Host           types.String            `tfsdk:"host"`
	Port           types.Int32             `tfsdk:"port"`
	Username       types.String            `tfsdk:"username"`
	Password       SecretReferenceModel    `tfsdk:"password"`
	JdbcProperties map[string]types.String `tfsdk:"jdbc_properties"`
}

// DatabaseConnection interface implementation for Db2IbmIConnectionModel
func (c *Db2IbmIConnectionModel) SetHost(v types.String)     { c.Host = v }
func (c *Db2IbmIConnectionModel) SetPort(v types.Int32)      { c.Port = v }
func (c *Db2IbmIConnectionModel) SetDatabase(v types.String) {} // no-op - DB2 doesn't have a database field
func (c *Db2IbmIConnectionModel) SetUsername(v types.String) { c.Username = v }
func (c *Db2IbmIConnectionModel) GetPasswordRef() *SecretReferenceModel {
	return &c.Password
}
func (c *Db2IbmIConnectionModel) GetJdbcPropertiesRef() *map[string]types.String {
	return &c.JdbcProperties
}

// SnowflakeAuthenticationModel represents Snowflake authentication configuration
type SnowflakeAuthenticationModel struct {
	PrivateKey SnowflakeSecretReferenceModel  `tfsdk:"private_key"`
	Passphrase *SnowflakeSecretReferenceModel `tfsdk:"passphrase"`
}

// SnowflakeConnectionModel represents Snowflake connection configuration
type SnowflakeConnectionModel struct {
	AccountName    types.String                 `tfsdk:"account_name"`
	Username       types.String                 `tfsdk:"username"`
	Authentication SnowflakeAuthenticationModel `tfsdk:"authentication"`
	JdbcProperties map[string]types.String      `tfsdk:"jdbc_properties"`
}
