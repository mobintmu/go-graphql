package utils

import "database/sql"

// =============================================================================
// UTILITY FUNCTIONS: GraphQL Pointer → SQLC Nullable
// =============================================================================

// ToNullInt32 converts *int to sql.NullInt32
func ToNullInt32(val *int) sql.NullInt32 {
	if val == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*val), Valid: true}
}

// ToNullInt64 converts *int64 to sql.NullInt64
func ToNullInt64(val *int64) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *val, Valid: true}
}

// ToNullString converts *string to sql.NullString
func ToNullString(val *string) sql.NullString {
	if val == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *val, Valid: true}
}

// ToNullBool converts *bool to sql.NullBool
func ToNullBool(val *bool) sql.NullBool {
	if val == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *val, Valid: true}
}

// =============================================================================
// UTILITY FUNCTIONS: SQLC Nullable → GraphQL Pointer
// =============================================================================

// NullInt32ToPointer converts sql.NullInt32 to *int
func NullInt32ToPointer(val sql.NullInt32) *int {
	if !val.Valid {
		return nil
	}
	result := int(val.Int32)
	return &result
}

// NullInt64ToPointer converts sql.NullInt64 to *int64
func NullInt64ToPointer(val sql.NullInt64) *int64 {
	if !val.Valid {
		return nil
	}
	return &val.Int64
}

// NullStringToPointer converts sql.NullString to *string
func NullStringToPointer(val sql.NullString) *string {
	if !val.Valid || val.String == "" {
		return nil
	}
	return &val.String
}

// NullBoolToPointer converts sql.NullBool to *bool
func NullBoolToPointer(val sql.NullBool) *bool {
	if !val.Valid {
		return nil
	}
	return &val.Bool
}

// =============================================================================
// UTILITY FUNCTIONS: With Default Values
// =============================================================================

// GetBoolOrDefault returns the bool value or default if nil
func GetBoolOrDefault(val *bool, defaultVal bool) bool {
	if val == nil {
		return defaultVal
	}
	return *val
}

// GetStringOrDefault returns the string value or default if nil
func GetStringOrDefault(val *string, defaultVal string) string {
	if val == nil {
		return defaultVal
	}
	return *val
}

// GetInt64OrDefault returns the int64 value or default if nil
func GetInt64OrDefault(val *int64, defaultVal int64) int64 {
	if val == nil {
		return defaultVal
	}
	return *val
}
