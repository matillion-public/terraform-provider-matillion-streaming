package provider

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-matillion-streaming/internal/client"
	"terraform-provider-matillion-streaming/internal/provider/models"
)

// ============================================================================
// Main Converters: Client API → Terraform Model
// ============================================================================

// mapSourceFromAPI maps the source from the API response to the state model
func mapSourceFromAPI(sourceData map[string]interface{}, state *models.PipelineResourceModel) error {
	sourceType, ok := sourceData["type"].(string)
	if !ok {
		return errors.New("could not read source type")
	}

	// Clear all source fields
	state.PostgresSource = nil
	state.SqlServerSource = nil
	state.MySqlSource = nil
	state.OracleSource = nil
	state.Db2IbmISource = nil

	switch sourceType {
	case "postgres":
		state.PostgresSource = &models.PostgresSourceModel{}
		mapDatabaseConnection(sourceData, &state.PostgresSource.Connection, true)
		if tables, ok := getTablesList(sourceData); ok {
			state.PostgresSource.Tables = tables
		}

	case "sqlserver":
		state.SqlServerSource = &models.SqlServerSourceModel{}
		mapDatabaseConnection(sourceData, &state.SqlServerSource.Connection, true)
		if tables, ok := getTablesList(sourceData); ok {
			state.SqlServerSource.Tables = tables
		}

	case "mysql":
		state.MySqlSource = &models.MysqlSourceModel{}
		mapDatabaseConnection(sourceData, &state.MySqlSource.Connection, false)
		if tables, ok := getTablesList(sourceData); ok {
			state.MySqlSource.Tables = tables
		}

	case "oracle":
		state.OracleSource = &models.OracleSourceModel{}
		mapDatabaseConnection(sourceData, &state.OracleSource.Connection, true)

		// Handle Oracle-specific field (pdb)
		connectionData, _ := sourceData["connection"].(map[string]interface{})
		if pdb, ok := getStringField(connectionData, "pdb"); ok {
			state.OracleSource.Connection.Pdb = types.StringValue(pdb)
		}

		if tables, ok := getTablesList(sourceData); ok {
			state.OracleSource.Tables = tables
		}

	case "db2ibmi":
		state.Db2IbmISource = &models.Db2IbmISourceModel{}
		mapDatabaseConnection(sourceData, &state.Db2IbmISource.Connection, false)
		if tables, ok := getTablesList(sourceData); ok {
			state.Db2IbmISource.Tables = tables
		}

	default:
		return errors.New("unknown source type: " + sourceType)
	}

	return nil
}

// mapTargetFromAPI maps the target from the API response to the state model
func mapTargetFromAPI(targetData map[string]interface{}, state *models.PipelineResourceModel) error {
	targetType, ok := targetData["type"].(string)
	if !ok {
		return errors.New("could not read target type")
	}

	// Clear all target fields
	state.SnowflakeTarget = nil
	state.S3Target = nil
	state.ABSTarget = nil

	switch targetType {
	case "snowflake":
		state.SnowflakeTarget = &models.SnowflakeTargetModel{}

		// Map direct fields
		directFieldMappings := map[string]*types.String{
			"role":        &state.SnowflakeTarget.Role,
			"warehouse":   &state.SnowflakeTarget.Warehouse,
			"database":    &state.SnowflakeTarget.Database,
			"stageSchema": &state.SnowflakeTarget.StageSchema,
			"stageName":   &state.SnowflakeTarget.StageName,
			"stagePrefix": &state.SnowflakeTarget.StagePrefix,
			"tableSchema": &state.SnowflakeTarget.TableSchema,
		}
		mapStringFields(targetData, directFieldMappings)

		// Handle case-sensitive fields with conversion
		if tablePrefixType, ok := getStringField(targetData, "tablePrefixType"); ok {
			state.SnowflakeTarget.TablePrefixType = types.StringValue(strings.ToLower(tablePrefixType))
		}
		if transformationType, ok := getStringField(targetData, "transformationType"); ok {
			state.SnowflakeTarget.TransformationType = types.StringValue(transformationTypeFromAPI(transformationType))
		}
		if temporalMapping, ok := getStringField(targetData, "temporalMapping"); ok {
			state.SnowflakeTarget.TemporalMapping = types.StringValue(strings.ToLower(temporalMapping))
		}

		// Map connection if present
		if connectionData, ok := targetData["connection"].(map[string]interface{}); ok {
			state.SnowflakeTarget.Connection = models.SnowflakeConnectionModel{}

			connectionFieldMappings := map[string]*types.String{
				"accountName": &state.SnowflakeTarget.Connection.AccountName,
				"username":    &state.SnowflakeTarget.Connection.Username,
			}
			mapStringFields(connectionData, connectionFieldMappings)

			// Map authentication if present
			if authData, ok := connectionData["authentication"].(map[string]interface{}); ok {
				state.SnowflakeTarget.Connection.Authentication = models.SnowflakeAuthenticationModel{}

				// Map private key
				if privateKeyData, ok := authData["privateKey"].(map[string]interface{}); ok {
					state.SnowflakeTarget.Connection.Authentication.PrivateKey = models.SnowflakeSecretReferenceModel{}
					if secretType, ok := getStringField(privateKeyData, "secretType"); ok {
						state.SnowflakeTarget.Connection.Authentication.PrivateKey.SecretType = types.StringValue(secretTypeFromAPI(secretType))
					}
					if secretName, ok := getStringField(privateKeyData, "secretName"); ok {
						state.SnowflakeTarget.Connection.Authentication.PrivateKey.SecretName = types.StringValue(secretName)
					}
				}

				// Map passphrase
				if passphraseData, ok := authData["passphrase"].(map[string]interface{}); ok {
					state.SnowflakeTarget.Connection.Authentication.Passphrase = &models.SnowflakeSecretReferenceModel{}
					if secretType, ok := getStringField(passphraseData, "secretType"); ok {
						state.SnowflakeTarget.Connection.Authentication.Passphrase.SecretType = types.StringValue(secretTypeFromAPI(secretType))
					}
					if secretName, ok := getStringField(passphraseData, "secretName"); ok {
						state.SnowflakeTarget.Connection.Authentication.Passphrase.SecretName = types.StringValue(secretName)
					}
				}
			}

			mapJdbcProperties(connectionData, &state.SnowflakeTarget.Connection.JdbcProperties)
		}

	case "s3":
		state.S3Target = &models.S3TargetModel{}

		// Map bucket (required)
		if bucket, ok := getStringField(targetData, "bucket"); ok {
			state.S3Target.Bucket = types.StringValue(bucket)
		}

		// Map prefix (optional) - set to null if empty
		if prefix, ok := getStringField(targetData, "prefix"); ok && prefix != "" {
			state.S3Target.Prefix = types.StringValue(prefix)
		} else {
			state.S3Target.Prefix = types.StringNull()
		}

		// Handle decimal mapping with case conversion
		if decimalMapping, ok := getStringField(targetData, "decimalMapping"); ok {
			state.S3Target.DecimalMapping = types.StringValue(strings.ToLower(decimalMapping))
		} else {
			state.S3Target.DecimalMapping = types.StringNull()
		}

	case "abs":
		state.ABSTarget = &models.AbsTargetModel{}

		// Map container and accountName (required fields)
		if container, ok := getStringField(targetData, "container"); ok {
			state.ABSTarget.Container = types.StringValue(container)
		}
		if accountName, ok := getStringField(targetData, "accountName"); ok {
			state.ABSTarget.AccountName = types.StringValue(accountName)
		}

		// Map prefix (optional) - set to null if empty
		if prefix, ok := getStringField(targetData, "prefix"); ok && prefix != "" {
			state.ABSTarget.Prefix = types.StringValue(prefix)
		} else {
			state.ABSTarget.Prefix = types.StringNull()
		}

		// Handle decimal mapping with case conversion
		if decimalMapping, ok := getStringField(targetData, "decimalMapping"); ok {
			state.ABSTarget.DecimalMapping = types.StringValue(strings.ToLower(decimalMapping))
		} else {
			state.ABSTarget.DecimalMapping = types.StringNull()
		}

		// Map account key secret reference
		if accountKeyData, ok := targetData["accountKey"].(map[string]interface{}); ok {
			state.ABSTarget.AccountKey = models.SecretReferenceModel{}
			if secretType, ok := getStringField(accountKeyData, "secretType"); ok {
				state.ABSTarget.AccountKey.Type = types.StringValue(secretTypeFromAPI(secretType))
			}
			if secretName, ok := getStringField(accountKeyData, "secretName"); ok {
				state.ABSTarget.AccountKey.Name = types.StringValue(secretName)
			}
			if key, ok := getStringField(accountKeyData, "key"); ok {
				state.ABSTarget.AccountKey.Key = types.StringValue(key)
			}
		}

	default:
		return errors.New("unknown target type: " + targetType)
	}

	return nil
}

// ============================================================================
// Database Connection Helpers
// ============================================================================

// mapDatabaseConnection is a generic helper to map database connection fields for sources with similar structures
// The hasDatabase parameter indicates whether the database field should be mapped (Postgres, SQL Server, Oracle have it; MySQL, DB2 don't)
func mapDatabaseConnection(sourceData map[string]interface{}, connection models.DatabaseConnection, hasDatabase bool) {
	connectionData, ok := sourceData["connection"].(map[string]interface{})
	if !ok {
		return
	}

	if host, ok := getStringField(connectionData, "host"); ok {
		connection.SetHost(types.StringValue(host))
	}
	if port, ok := connectionData["port"].(float64); ok {
		connection.SetPort(types.Int32Value(int32(port)))
	}
	if hasDatabase {
		if database, ok := getStringField(connectionData, "database"); ok {
			connection.SetDatabase(types.StringValue(database))
		}
	}
	if username, ok := getStringField(connectionData, "username"); ok {
		connection.SetUsername(types.StringValue(username))
	}
	mapPasswordSecretReference(connectionData, connection.GetPasswordRef())
	mapJdbcProperties(connectionData, connection.GetJdbcPropertiesRef())
}

// ============================================================================
// Secret & Authentication Helpers
// ============================================================================

// mapPasswordSecretReference maps a password secret reference from connection data
func mapPasswordSecretReference(connectionData map[string]interface{}, password *models.SecretReferenceModel) {
	if passwordData, ok := connectionData["password"].(map[string]interface{}); ok {
		*password = models.SecretReferenceModel{}
		if secretType, ok := getStringField(passwordData, "secretType"); ok {
			password.Type = types.StringValue(secretTypeFromAPI(secretType))
		}
		if secretName, ok := getStringField(passwordData, "secretName"); ok {
			password.Name = types.StringValue(secretName)
		}
		if key, ok := getStringField(passwordData, "key"); ok {
			password.Key = types.StringValue(key)
		}
	}
}

// ============================================================================
// Property Mapping Helpers
// ============================================================================

// mapJdbcProperties maps JDBC properties from connection data for database sources
func mapJdbcProperties(connectionData map[string]interface{}, properties *map[string]types.String) {
	if jdbcProps, ok := connectionData["jdbcProperties"].(map[string]interface{}); ok {
		*properties = make(map[string]types.String)
		for k, v := range jdbcProps {
			if strVal, ok := v.(string); ok {
				(*properties)[k] = types.StringValue(strVal)
			}
		}
	}
}

// ============================================================================
// Generic Field Mapping Helpers
// ============================================================================

// getStringField extracts a string field from a map
func getStringField(m map[string]interface{}, key string) (string, bool) {
	if val, ok := m[key].(string); ok {
		return val, true
	}
	return "", false
}

// getTablesList extracts a tables list from a map
func getTablesList(m map[string]interface{}) ([]models.TableResourceModel, bool) {
	raw, exists := m["tables"]
	if !exists {
		return nil, false
	}

	// Handle both []interface{} (from API) and []client.Table (from our models)
	var stateTables []models.TableResourceModel

	switch tablesData := raw.(type) {
	case []interface{}:
		// API response format: array of maps
		for _, tableInterface := range tablesData {
			if tableMap, ok := tableInterface.(map[string]interface{}); ok {
				table := models.TableResourceModel{}
				if tableSchema, ok := getStringField(tableMap, "schema"); ok {
					table.Schema = types.StringValue(tableSchema)
				}
				if tableName, ok := getStringField(tableMap, "table"); ok {
					table.Table = types.StringValue(tableName)
				}
				stateTables = append(stateTables, table)
			}
		}
		return stateTables, true
	case []client.Table:
		// Internal client format
		for _, t := range tablesData {
			stateTables = append(stateTables, models.TableResourceModel{Schema: types.StringValue(t.Schema), Table: types.StringValue(t.Table)})
		}
		return stateTables, true
	default:
		return nil, false
	}
}

// mapStringFields applies string field mappings from a source map to target fields
func mapStringFields(source map[string]interface{}, fieldMappings map[string]*types.String) {
	for key, field := range fieldMappings {
		if value, ok := getStringField(source, key); ok {
			*field = types.StringValue(value)
		}
	}
}

// ============================================================================
// Format Converters: API Format → Terraform Format
// ============================================================================

// secretTypeFromAPI converts API secret type values to lowercase terraform format
func secretTypeFromAPI(apiValue string) string {
	switch apiValue {
	case "AWS_SECRETS_MANAGER":
		return "aws_secrets_manager"
	case "AZURE_KEY_VAULT":
		return "azure_key_vault"
	case "GOOGLE_SECRET_MANAGER":
		return "google_secret_manager"
	default:
		return strings.ToLower(apiValue) // fallback
	}
}

// transformationTypeFromAPI converts API values to lowercase terraform format
func transformationTypeFromAPI(apiValue string) string {
	switch apiValue {
	case "COPY_TABLE":
		return "copy_table"
	case "COPY_TABLE_SOFT_DELETE":
		return "copy_table_soft_delete"
	case "CHANGE_LOG":
		return "change_log"
	default:
		return strings.ToLower(apiValue) // fallback
	}
}
