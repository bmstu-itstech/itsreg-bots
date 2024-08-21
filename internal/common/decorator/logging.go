package decorator

import (
	"context"
	"fmt"
	"log/slog"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *slog.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := d.logger.
		With("command", handlerType).
		With("command_body", fmt.Sprintf("%v", cmd))

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.With("error", err.Error()).
				Error("Failed to execute command")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger *slog.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := d.logger.
		With("query", generateActionName(cmd)).
		With("query_body", fmt.Sprintf("%v", cmd))

	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query executed successfully")
		} else {
			logger.With("error", err.Error()).
				Error("Failed to execute query")
		}
	}()

	return d.base.Handle(ctx, cmd)
}
