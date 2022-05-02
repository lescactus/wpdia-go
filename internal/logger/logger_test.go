package logger

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	newLogger = logrus.New()
)

func TestNewLogger(t *testing.T) {
	type args struct {
		loglevel string
		format   string
	}
	tests := []struct {
		name    string
		args    args
		want    *logrus.Logger
		wantErr bool
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
			want: &logrus.Logger{
				Out:          os.Stdout,
				Formatter:    &logrus.TextFormatter{},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.InfoLevel,
				ExitFunc:     os.Exit,
				ReportCaller: false,
			},
			wantErr: false,
		},
		{
			name: "Loglevel: info / Format: text",
			args: args{
				loglevel: "info",
				format:   "text",
			},
			want: &logrus.Logger{
				Out:          os.Stdout,
				Formatter:    &logrus.TextFormatter{},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.InfoLevel,
				ExitFunc:     os.Exit,
				ReportCaller: false,
			},
			wantErr: false,
		},
		{
			name: "Loglevel: info / Format: json",
			args: args{
				loglevel: "info",
				format:   "json",
			},
			want: &logrus.Logger{
				Out:          os.Stdout,
				Formatter:    &logrus.JSONFormatter{},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.InfoLevel,
				ExitFunc:     os.Exit,
				ReportCaller: false,
			},
			wantErr: false,
		},
		{
			name: "Loglevel: debug / Format: text",
			args: args{
				loglevel: "debug",
				format:   "text",
			},
			want: &logrus.Logger{
				Out:          os.Stdout,
				Formatter:    &logrus.TextFormatter{},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.DebugLevel,
				ExitFunc:     os.Exit,
				ReportCaller: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLogger(tt.args.loglevel, tt.args.format)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.ObjectsAreEqualValues(tt.want, got)
		})
	}
}
