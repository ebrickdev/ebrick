package logger

import "testing"

func TestGetLevel(t *testing.T) {
	type args struct {
		levelStr string
	}
	tests := []struct {
		name    string
		args    args
		want    Level
		wantErr bool
	}{
		{
			name:    "Valid level DEBUG",
			args:    args{levelStr: "debug"},
			want:    DebugLevel,
			wantErr: false,
		},
		{
			name:    "Valid level INFO",
			args:    args{levelStr: "info"},
			want:    InfoLevel,
			wantErr: false,
		},
		{
			name:    "Valid level WARN",
			args:    args{levelStr: "warn"},
			want:    WarnLevel,
			wantErr: false,
		},
		{
			name:    "Valid level ERROR",
			args:    args{levelStr: "error"},
			want:    ErrorLevel,
			wantErr: false,
		},
		{
			name:    "Invalid level",
			args:    args{levelStr: "invalid"},
			want:    0, // Assuming 0 is the default value for invalid levels
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLevel(tt.args.levelStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
