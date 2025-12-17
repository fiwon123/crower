package crowlog

import (
	"strconv"

	"go.uber.org/zap"
)

type LoggerInfo struct {
	logger *zap.Logger
}

type field struct {
	key  string
	data any
}

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

func NewField(key string, data any) field {
	return field{
		key:  key,
		data: data,
	}
}

func translateData(data ...any) []zap.Field {
	var zapFields []zap.Field
	for i, d := range data {
		zapFields = append(zapFields, zap.Any("data"+strconv.Itoa(i), d))
	}

	return zapFields
}

func translateDataWithKeys(rawFields []field) []zap.Field {
	var zapFields []zap.Field
	for _, field := range rawFields {
		zapFields = append(zapFields, zap.Any(field.key, field.data))
	}

	return zapFields
}

func (infoData LoggerInfo) Info(msg string, data ...any) {
	fields := translateData(data)
	infoData.logger.Info(msg, fields...)
}

func (infoData LoggerInfo) InfoWithKeys(msg string, rawFields ...field) {
	fields := translateDataWithKeys(rawFields)
	infoData.logger.Info(msg, fields...)
}

func (infoData LoggerInfo) Error(msg string, data ...any) {
	fields := translateData(data)
	infoData.logger.Error(msg, fields...)
}

func (infoData LoggerInfo) ErrorWithKeys(msg string, rawFields ...field) {
	fields := translateDataWithKeys(rawFields)
	infoData.logger.Error(msg, fields...)
}
