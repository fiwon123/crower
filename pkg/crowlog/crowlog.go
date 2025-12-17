// Package wrapper using zap.
// Provides a easy way to print like fmt package.
package crowlog

import (
	"strconv"

	"go.uber.org/zap"
)

// The Crowlog wrapper to access logger provider.
type LoggerInfo struct {
	logger *zap.Logger
}

type field struct {
	key  string
	data any
}

// Create a new LoggerInfo pointer.
func New() *LoggerInfo {

	logger, _ := zap.NewProduction(
		zap.AddCaller(),      // enable caller info
		zap.AddCallerSkip(1), // skip your wrapper
	)
	defer logger.Sync()

	return &LoggerInfo{
		logger: logger,
	}
}

// Create a new raw field to represent logger provider field used in LoggerInfo struct.
func NewField(key string, data any) field {
	return field{
		key:  key,
		data: data,
	}
}

// Translate raw fields to logger provider field used in LoggerInfo struct.
func translateData(data ...any) []zap.Field {
	var zapFields []zap.Field
	for i, d := range data {
		zapFields = append(zapFields, zap.Any("data"+strconv.Itoa(i), d))
	}

	return zapFields
}

// Translate raw fields to logger provider field used in LoggerInfo struct.
func translateDataWithKeys(rawFields []field) []zap.Field {
	var zapFields []zap.Field
	for _, field := range rawFields {
		zapFields = append(zapFields, zap.Any(field.key, field.data))
	}

	return zapFields
}

// Print info message and data of any type.
// data can be string, int, slices, etc.
func (infoData LoggerInfo) Info(msg string, data ...any) {
	fields := translateData(data)
	infoData.logger.Info(msg, fields...)
}

// Print info message and data of any type.
// rawFields can be create using NewField method
// data can be string, int, slices, etc.
func (infoData LoggerInfo) InfoWithKeys(msg string, rawFields ...field) {
	fields := translateDataWithKeys(rawFields)
	infoData.logger.Info(msg, fields...)
}

// Print error message and data of any type.
// data can be string, int, slices, etc.
func (infoData LoggerInfo) Error(msg string, data ...any) {
	fields := translateData(data)
	infoData.logger.Error(msg, fields...)
}

// Print error message and data of any type.
// rawFields can be create using NewField method
// data can be string, int, slices, etc.
func (infoData LoggerInfo) ErrorWithKeys(msg string, rawFields ...field) {
	fields := translateDataWithKeys(rawFields)
	infoData.logger.Error(msg, fields...)
}
