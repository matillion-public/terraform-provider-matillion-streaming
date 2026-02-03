package provider

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-matillion-streaming/internal/client"
	"terraform-provider-matillion-streaming/internal/provider/models"
)

// ============================================================================
// Main Converters: Terraform Model → Client API
// ============================================================================

// mapSourceToClientSource converts the provider source model to the appropriate client API model
func mapSourceToClientSource(model *models.PipelineResourceModel) (interface{}, error) {
	if model == nil {
		return nil, errors.New("model is nil")
	}

	switch {
	case model.PostgresSource != nil:
		return client.PostgresSource{
			Source: client.Source{Type: "postgres"},
			Connection: client.PostgresConnection{
				Host:       model.PostgresSource.Connection.Host.ValueString(),
				Port:       int(model.PostgresSource.Connection.Port.ValueInt32()),
				Database:   model.PostgresSource.Connection.Database.ValueString(),
				Username:   model.PostgresSource.Connection.Username.ValueString(),
				Password:   buildSecretReference(model.PostgresSource.Connection.Password),
				Properties: extractPropertiesMap(model.PostgresSource.Connection.JdbcProperties),
			},
			Tables: extractTablesFromModel(model.PostgresSource.Tables),
		}, nil

	case model.SqlServerSource != nil:
		return client.SqlServerSource{
			Source: client.Source{Type: "sqlserver"},
			Connection: client.SqlServerConnection{
				Host:       model.SqlServerSource.Connection.Host.ValueString(),
				Port:       int(model.SqlServerSource.Connection.Port.ValueInt32()),
				Database:   model.SqlServerSource.Connection.Database.ValueString(),
				Username:   model.SqlServerSource.Connection.Username.ValueString(),
				Password:   buildSecretReference(model.SqlServerSource.Connection.Password),
				Properties: extractPropertiesMap(model.SqlServerSource.Connection.JdbcProperties),
			},
			Tables: extractTablesFromModel(model.SqlServerSource.Tables),
		}, nil

	case model.MySqlSource != nil:
		return client.MySqlSource{
			Source: client.Source{Type: "mysql"},
			Connection: client.MySqlConnection{
				Host:       model.MySqlSource.Connection.Host.ValueString(),
				Port:       int(model.MySqlSource.Connection.Port.ValueInt32()),
				Username:   model.MySqlSource.Connection.Username.ValueString(),
				Password:   buildSecretReference(model.MySqlSource.Connection.Password),
				Properties: extractPropertiesMap(model.MySqlSource.Connection.JdbcProperties),
			},
			Tables: extractTablesFromModel(model.MySqlSource.Tables),
		}, nil

	case model.OracleSource != nil:
		return client.OracleSource{
			Source: client.Source{Type: "oracle"},
			Connection: client.OracleConnection{
				Host:       model.OracleSource.Connection.Host.ValueString(),
				Port:       int(model.OracleSource.Connection.Port.ValueInt32()),
				Database:   model.OracleSource.Connection.Database.ValueString(),
				Pdb:        model.OracleSource.Connection.Pdb.ValueString(),
				Username:   model.OracleSource.Connection.Username.ValueString(),
				Password:   buildSecretReference(model.OracleSource.Connection.Password),
				Properties: extractPropertiesMap(model.OracleSource.Connection.JdbcProperties),
			},
			Tables: extractTablesFromModel(model.OracleSource.Tables),
		}, nil

	case model.Db2IbmISource != nil:
		return client.Db2IbmISource{
			Source: client.Source{Type: "db2ibmi"},
			Connection: client.Db2IbmIConnection{
				Host:       model.Db2IbmISource.Connection.Host.ValueString(),
				Port:       int(model.Db2IbmISource.Connection.Port.ValueInt32()),
				Username:   model.Db2IbmISource.Connection.Username.ValueString(),
				Password:   buildSecretReference(model.Db2IbmISource.Connection.Password),
				Properties: extractPropertiesMap(model.Db2IbmISource.Connection.JdbcProperties),
			},
			Tables: extractTablesFromModel(model.Db2IbmISource.Tables),
		}, nil
	}

	return nil, errors.New("no source type specified")
}

// mapTargetToClientTarget converts the provider target model to the appropriate client API model
// using concrete types to ensure all fields are properly serialized
func mapTargetToClientTarget(model *models.PipelineResourceModel) (interface{}, error) {
	if model == nil {
		return nil, errors.New("model is nil")
	}

	switch {
	case model.SnowflakeTarget != nil:
		return client.SnowflakeTargetModel{
			Target: client.Target{Type: "snowflake"},
			Connection: client.SnowflakeConnection{
				AccountName: model.SnowflakeTarget.Connection.AccountName.ValueString(),
				Username:    model.SnowflakeTarget.Connection.Username.ValueString(),
				Authentication: client.SnowflakeAuthentication{
					Type:       "key-pair",
					PrivateKey: buildSnowflakeSecretReference(model.SnowflakeTarget.Connection.Authentication.PrivateKey),
					Passphrase: buildSnowflakeSecretReferencePtr(model.SnowflakeTarget.Connection.Authentication.Passphrase),
				},
				Properties: extractPropertiesMap(model.SnowflakeTarget.Connection.JdbcProperties),
			},
			Role:               model.SnowflakeTarget.Role.ValueString(),
			Warehouse:          model.SnowflakeTarget.Warehouse.ValueString(),
			Database:           model.SnowflakeTarget.Database.ValueString(),
			StageSchema:        model.SnowflakeTarget.StageSchema.ValueString(),
			StageName:          model.SnowflakeTarget.StageName.ValueString(),
			StagePrefix:        model.SnowflakeTarget.StagePrefix.ValueString(),
			TableSchema:        model.SnowflakeTarget.TableSchema.ValueString(),
			TablePrefixType:    model.SnowflakeTarget.TablePrefixType.ValueString(),
			TransformationType: model.SnowflakeTarget.TransformationType.ValueString(),
			TemporalMapping:    model.SnowflakeTarget.TemporalMapping.ValueString(),
		}, nil

	case model.S3Target != nil:
		return client.S3TargetModel{
			Target:         client.Target{Type: "s3"},
			Bucket:         model.S3Target.Bucket.ValueString(),
			Prefix:         model.S3Target.Prefix.ValueString(),
			DecimalMapping: model.S3Target.DecimalMapping.ValueString(),
		}, nil

	case model.ABSTarget != nil:
		return client.AbsTargetModel{
			Target:         client.Target{Type: "abs"},
			Container:      model.ABSTarget.Container.ValueString(),
			Prefix:         model.ABSTarget.Prefix.ValueString(),
			AccountName:    model.ABSTarget.AccountName.ValueString(),
			AccountKey:     buildSecretReference(model.ABSTarget.AccountKey),
			DecimalMapping: model.ABSTarget.DecimalMapping.ValueString(),
		}, nil
	default:
		return nil, errors.New("no target type specified")
	}
}

// ============================================================================
// Helper Functions: Model to Client Conversions
// ============================================================================

// extractTablesFromModel converts Terraform table models to client table models
func extractTablesFromModel(tables []models.TableResourceModel) []client.Table {
	result := make([]client.Table, len(tables))
	for i, t := range tables {
		result[i] = client.Table{
			Schema: t.Schema.ValueString(),
			Table:  t.Table.ValueString(),
		}
	}
	return result
}

// extractPropertiesMap converts Terraform string map to regular string map pointer
func extractPropertiesMap(properties map[string]types.String) *map[string]string {
	if properties == nil {
		return nil
	}
	result := make(map[string]string)
	for k, v := range properties {
		result[k] = v.ValueString()
	}
	return &result
}

// buildSecretReference converts Terraform secret reference model to client model
func buildSecretReference(password models.SecretReferenceModel) client.PipelineSecretReference {
	return client.PipelineSecretReference{
		SecretType: password.Type.ValueString(),
		SecretName: password.Name.ValueString(),
		Key:        password.Key.ValueString(),
	}
}

// buildSnowflakeSecretReference converts Snowflake secret reference model to client model
// Snowflake uses a different secret reference structure without the Key field
func buildSnowflakeSecretReference(secret models.SnowflakeSecretReferenceModel) client.PipelineSecretReference {
	return client.PipelineSecretReference{
		SecretType: secret.SecretType.ValueString(),
		SecretName: secret.SecretName.ValueString(),
	}
}

// buildSnowflakeSecretReferencePtr converts Snowflake secret reference pointer to client model pointer
// Returns nil if the input is nil, otherwise converts and returns a pointer to the result
func buildSnowflakeSecretReferencePtr(secret *models.SnowflakeSecretReferenceModel) *client.PipelineSecretReference {
	if secret == nil {
		return nil
	}
	ref := buildSnowflakeSecretReference(*secret)
	return &ref
}
