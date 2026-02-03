package client

// Agent represents an agent in the Matillion API
type Agent struct {
	AgentId       string `json:"agentId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Deployment    string `json:"deployment"`
	CloudProvider string `json:"cloudProvider"`
	AgentType     string `json:"agentType"`
}

// Pipeline represents a pipeline in the streaming platform
type Pipeline struct {
	StreamingPipelineId string                 `json:"streamingPipelineId"`
	Name                string                 `json:"name"`
	AgentId             string                 `json:"agentId"`
	StreamingSource     map[string]interface{} `json:"streamingSource"`
	StreamingTarget     map[string]interface{} `json:"streamingTarget"`
	AdvancedProperties  map[string]string      `json:"advancedProperties"`
}

// Source is the base type for all source types
type Source struct {
	Type string `json:"type"`
}

// Target is the base type for all target types
type Target struct {
	Type string `json:"type"`
}

// Table represents a database table
type Table struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
}

// PostgresSource represents a PostgreSQL source configuration
type PostgresSource struct {
	Source
	Connection PostgresConnection `json:"connection"`
	Tables     []Table            `json:"tables"`
}

// SqlServerSource represents a SQL Server source configuration
type SqlServerSource struct {
	Source
	Connection SqlServerConnection `json:"connection"`
	Tables     []Table             `json:"tables"`
}

// MySqlSource represents a MySQL source configuration
type MySqlSource struct {
	Source
	Connection MySqlConnection `json:"connection"`
	Tables     []Table         `json:"tables"`
}

// OracleSource represents an Oracle source configuration
type OracleSource struct {
	Source
	Connection OracleConnection `json:"connection"`
	Tables     []Table          `json:"tables"`
}

// Db2IbmISource represents a DB2 for IBM i source configuration
type Db2IbmISource struct {
	Source
	Connection Db2IbmIConnection `json:"connection"`
	Tables     []Table           `json:"tables"`
}

// PostgresConnection represents a PostgreSQL connection configuration
type PostgresConnection struct {
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Database   string                  `json:"database"`
	Username   string                  `json:"username"`
	Password   PipelineSecretReference `json:"password"`
	Properties *map[string]string      `json:"jdbcProperties,omitempty"`
}

// SqlServerConnection represents a SQL Server connection configuration
type SqlServerConnection struct {
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Database   string                  `json:"database"`
	Username   string                  `json:"username"`
	Password   PipelineSecretReference `json:"password"`
	Properties *map[string]string      `json:"jdbcProperties,omitempty"`
}

// MySqlConnection represents a MySQL connection configuration
type MySqlConnection struct {
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Username   string                  `json:"username"`
	Password   PipelineSecretReference `json:"password"`
	Properties *map[string]string      `json:"jdbcProperties,omitempty"`
}

// OracleConnection represents an Oracle connection configuration
type OracleConnection struct {
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Database   string                  `json:"database"`
	Pdb        string                  `json:"pdb,omitempty"`
	Username   string                  `json:"username"`
	Password   PipelineSecretReference `json:"password"`
	Properties *map[string]string      `json:"jdbcProperties,omitempty"`
}

// Db2IbmIConnection represents a DB2 for IBM i connection configuration
type Db2IbmIConnection struct {
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Username   string                  `json:"username"`
	Password   PipelineSecretReference `json:"password"`
	Properties *map[string]string      `json:"jdbcProperties,omitempty"`
}

// PipelineSecretReference represents a secret reference in a streaming pipeline
type PipelineSecretReference struct {
	SecretType string `json:"secretType"`
	SecretName string `json:"secretName"`
	Key        string `json:"key,omitempty"`
}

// SnowflakeAuthentication represents Snowflake authentication configuration
type SnowflakeAuthentication struct {
	Type       string                   `json:"type"`
	PrivateKey PipelineSecretReference  `json:"privateKey"`
	Passphrase *PipelineSecretReference `json:"passphrase,omitempty"`
}

// SnowflakeConnection represents a Snowflake connection configuration
type SnowflakeConnection struct {
	AccountName    string                  `json:"accountName"`
	Username       string                  `json:"username"`
	Authentication SnowflakeAuthentication `json:"authentication"`
	Properties     *map[string]string      `json:"jdbcProperties,omitempty"`
}

// SnowflakeTargetModel represents a Snowflake target configuration
type SnowflakeTargetModel struct {
	Target
	Connection         SnowflakeConnection `json:"connection"`
	Role               string              `json:"role"`
	Warehouse          string              `json:"warehouse"`
	Database           string              `json:"database"`
	StageSchema        string              `json:"stageSchema"`
	StageName          string              `json:"stageName"`
	StagePrefix        string              `json:"stagePrefix,omitempty"`
	TableSchema        string              `json:"tableSchema"`
	TablePrefixType    string              `json:"tablePrefixType"`
	TransformationType string              `json:"transformationType"`
	TemporalMapping    string              `json:"temporalMapping,omitempty"`
}

// S3TargetModel represents an S3 target configuration
type S3TargetModel struct {
	Target
	Bucket         string `json:"bucket"`
	Prefix         string `json:"prefix"`
	DecimalMapping string `json:"decimalMapping,omitempty"`
}

// AbsTargetModel represents an Azure Blob Storage target configuration
type AbsTargetModel struct {
	Target
	Container      string                  `json:"container"`
	Prefix         string                  `json:"prefix"`
	AccountName    string                  `json:"accountName"`
	AccountKey     PipelineSecretReference `json:"accountKey"`
	DecimalMapping string                  `json:"decimalMapping,omitempty"`
}
