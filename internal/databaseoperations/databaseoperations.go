// Package databaseoperations provides a system-aware routing layer for database operation data.
// This package delegates to registered providers instead of hard-coded switch statements.
package databaseoperations

import (
	"sort"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"

	hsdb "github.com/OperationsPAI/chaos-experiment/internal/hs/databaseoperations"
	mediadb "github.com/OperationsPAI/chaos-experiment/internal/media/databaseoperations"
	obdb "github.com/OperationsPAI/chaos-experiment/internal/ob/databaseoperations"
	oteldemodb "github.com/OperationsPAI/chaos-experiment/internal/oteldemo/databaseoperations"
	sndb "github.com/OperationsPAI/chaos-experiment/internal/sn/databaseoperations"
	sockshopdb "github.com/OperationsPAI/chaos-experiment/internal/sockshop/databaseoperations"
	teastoredb "github.com/OperationsPAI/chaos-experiment/internal/teastore/databaseoperations"
	tsdb "github.com/OperationsPAI/chaos-experiment/internal/ts/databaseoperations"
)

// DatabaseOperation represents a database operation from ClickHouse analysis.
type DatabaseOperation struct {
	ServiceName   string
	DBName        string
	DBTable       string
	Operation     string
	DBSystem      string
	ServerAddress string
	ServerPort    string
}

type staticDatabaseOperationProvider struct {
	operations map[string][]DatabaseOperation
}

func init() {
	registry := systemconfig.GetRegistry()
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemTrainTicket, newStaticDatabaseOperationProvider(convertTSOperationMap(tsdb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemOtelDemo, newStaticDatabaseOperationProvider(convertOtelDemoOperationMap(oteldemodb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemMediaMicroservices, newStaticDatabaseOperationProvider(convertMediaOperationMap(mediadb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemHotelReservation, newStaticDatabaseOperationProvider(convertHSOperationMap(hsdb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemSocialNetwork, newStaticDatabaseOperationProvider(convertSNOperationMap(sndb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemOnlineBoutique, newStaticDatabaseOperationProvider(convertOBOperationMap(obdb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemSockShop, newStaticDatabaseOperationProvider(convertSockShopOperationMap(sockshopdb.DatabaseOperations)))
	registry.RegisterDatabaseOperationProvider(systemconfig.SystemTeaStore, newStaticDatabaseOperationProvider(convertTeaStoreOperationMap(teastoredb.DatabaseOperations)))
}

func newStaticDatabaseOperationProvider(operations map[string][]DatabaseOperation) systemconfig.DatabaseOperationProvider {
	return &staticDatabaseOperationProvider{operations: operations}
}

func (p *staticDatabaseOperationProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.operations))
	for service := range p.operations {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *staticDatabaseOperationProvider) GetOperationsByService(serviceName string) []systemconfig.DatabaseOperationData {
	operations := p.operations[serviceName]
	result := make([]systemconfig.DatabaseOperationData, len(operations))
	for i, operation := range operations {
		result[i] = systemconfig.DatabaseOperationData{
			ServiceName:   operation.ServiceName,
			DBName:        operation.DBName,
			DBTable:       operation.DBTable,
			Operation:     operation.Operation,
			DBSystem:      operation.DBSystem,
			ServerAddress: operation.ServerAddress,
			ServerPort:    operation.ServerPort,
		}
	}
	return result
}

// GetOperationsByService returns all database operations for a service based on the current system.
func GetOperationsByService(serviceName string) []DatabaseOperation {
	data, err := systemconfig.GetMetadataStore().GetDatabaseOperations(string(systemconfig.GetCurrentSystem()), serviceName)
	if err == nil && len(data) > 0 {
		result := make([]DatabaseOperation, len(data))
		for i, operation := range data {
			result[i] = DatabaseOperation{
				ServiceName:   operation.ServiceName,
				DBName:        operation.DBName,
				DBTable:       operation.DBTable,
				Operation:     operation.Operation,
				DBSystem:      operation.DBSystem,
				ServerAddress: operation.ServerAddress,
				ServerPort:    operation.ServerPort,
			}
		}
		return result
	}
	return []DatabaseOperation{}
}

// GetAllDatabaseServices returns a list of all services that perform database operations.
func GetAllDatabaseServices() []string {
	names, err := systemconfig.GetMetadataStore().GetAllServiceNames(string(systemconfig.GetCurrentSystem()))
	if err == nil && len(names) > 0 {
		var services []string
		for _, service := range names {
			if len(GetOperationsByService(service)) > 0 {
				services = append(services, service)
			}
		}
		if len(services) > 0 {
			return services
		}
	}
	return []string{}
}

// GetOperationsByDatabase returns all operations for a specific database.
func GetOperationsByDatabase(dbName string) []DatabaseOperation {
	var results []DatabaseOperation
	for _, service := range GetAllDatabaseServices() {
		for _, operation := range GetOperationsByService(service) {
			if operation.DBName == dbName {
				results = append(results, operation)
			}
		}
	}
	return results
}

// GetOperationsByTable returns all operations for a specific table.
func GetOperationsByTable(dbTable string) []DatabaseOperation {
	var results []DatabaseOperation
	for _, service := range GetAllDatabaseServices() {
		for _, operation := range GetOperationsByService(service) {
			if operation.DBTable == dbTable {
				results = append(results, operation)
			}
		}
	}
	return results
}

func convertTSOperationMap(tsOps map[string][]tsdb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(tsOps))
	for service, operations := range tsOps {
		result[service] = convertTSOperations(operations)
	}
	return result
}

func convertOtelDemoOperationMap(otelOps map[string][]oteldemodb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(otelOps))
	for service, operations := range otelOps {
		result[service] = convertOtelDemoOperations(operations)
	}
	return result
}

func convertMediaOperationMap(mediaOps map[string][]mediadb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(mediaOps))
	for service, operations := range mediaOps {
		result[service] = convertMediaOperations(operations)
	}
	return result
}

func convertHSOperationMap(hsOps map[string][]hsdb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(hsOps))
	for service, operations := range hsOps {
		result[service] = convertHSOperations(operations)
	}
	return result
}

func convertSNOperationMap(snOps map[string][]sndb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(snOps))
	for service, operations := range snOps {
		result[service] = convertSNOperations(operations)
	}
	return result
}

func convertOBOperationMap(obOps map[string][]obdb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(obOps))
	for service, operations := range obOps {
		result[service] = convertOBOperations(operations)
	}
	return result
}

func convertSockShopOperationMap(sockshopOps map[string][]sockshopdb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(sockshopOps))
	for service, operations := range sockshopOps {
		result[service] = convertSockShopOperations(operations)
	}
	return result
}

func convertTeaStoreOperationMap(teastoreOps map[string][]teastoredb.DatabaseOperation) map[string][]DatabaseOperation {
	result := make(map[string][]DatabaseOperation, len(teastoreOps))
	for service, operations := range teastoreOps {
		result[service] = convertTeaStoreOperations(operations)
	}
	return result
}

// convertTSOperations converts ts-specific operations to the common type.
func convertTSOperations(tsOps []tsdb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(tsOps))
	for i, op := range tsOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertOtelDemoOperations converts otel-demo-specific operations to the common type.
func convertOtelDemoOperations(otelOps []oteldemodb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(otelOps))
	for i, op := range otelOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertMediaOperations converts media-specific operations to the common type.
func convertMediaOperations(mediaOps []mediadb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(mediaOps))
	for i, op := range mediaOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertHSOperations converts hs-specific operations to the common type.
func convertHSOperations(hsOps []hsdb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(hsOps))
	for i, op := range hsOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertSNOperations converts sn-specific operations to the common type.
func convertSNOperations(snOps []sndb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(snOps))
	for i, op := range snOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertOBOperations converts ob-specific operations to the common type.
func convertOBOperations(obOps []obdb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(obOps))
	for i, op := range obOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertSockShopOperations converts sockshop-specific operations to the common type.
func convertSockShopOperations(sockshopOps []sockshopdb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(sockshopOps))
	for i, op := range sockshopOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}

// convertTeaStoreOperations converts teastore-specific operations to the common type.
func convertTeaStoreOperations(teastoreOps []teastoredb.DatabaseOperation) []DatabaseOperation {
	result := make([]DatabaseOperation, len(teastoreOps))
	for i, op := range teastoreOps {
		result[i] = DatabaseOperation{
			ServiceName:   op.ServiceName,
			DBName:        op.DBName,
			DBTable:       op.DBTable,
			Operation:     op.Operation,
			DBSystem:      op.DBSystem,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
		}
	}
	return result
}
