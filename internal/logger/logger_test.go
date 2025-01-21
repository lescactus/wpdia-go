package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var (
// 	newLogger = slog.New(slog.NewTextHandler(os.Stderr, nil))
// )

func TestNewLogger(t *testing.T) {
	type args struct {
		loglevel string
		format   string
	}
	tests := []struct {
		name       string
		args       args
		wantLogger *slog.Logger
		wantErr    bool
	}{
		{
			name: "Loglevel: empty / Format: empty",
			args: args{
				loglevel: "",
				format:   "",
			},
			wantErr: true,
		},
		{
			name: "Loglevel: invalid / Format: empty",
			args: args{
				loglevel: "invalid",
				format:   "",
			},
			wantErr: true,
		},
		{
			name: "Loglevel: empty / Format: invalid",
			args: args{
				loglevel: "",
				format:   "invalid",
			},
			wantErr: true,
		},
		{
			name: "Loglevel: info / Format: invalid",
			args: args{
				loglevel: "info",
				format:   "invalid",
			},
			wantLogger: &slog.Logger{},
			wantErr:    false,
		},
		{
			name: "Loglevel: info / Format: text",
			args: args{
				loglevel: "info",
				format:   "text",
			},
			wantLogger: &slog.Logger{},
			wantErr:    false,
		},
		{
			name: "Loglevel: info / Format: json",
			args: args{
				loglevel: "info",
				format:   "json",
			},
			wantLogger: &slog.Logger{},
			wantErr:    false,
		},
		{
			name: "Loglevel: debug / Format: text",
			args: args{
				loglevel: "debug",
				format:   "text",
			},
			wantLogger: &slog.Logger{},
			wantErr:    false,
		},
		{
			name: "Loglevel: warn / Format: text",
			args: args{
				loglevel: "warn",
				format:   "text",
			},
			wantLogger: &slog.Logger{},
			wantErr:    false,
		},
		{
			name: "Loglevel: error / Format: text",
			args: args{
				loglevel: "error",
				format:   "text",
			},
			wantLogger: &slog.Logger{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.loglevel, tt.args.format)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.ObjectsAreEqualValues(tt.wantLogger, got)
		})
	}
}
