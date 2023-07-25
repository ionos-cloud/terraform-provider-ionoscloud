package utils

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestWaitForResourceToBeReady(t *testing.T) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	type args struct {
		ctx context.Context
		d   *schema.ResourceData
		fn  ResourceReadyFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestReturnTrue",
			args: args{
				ctx: context.TODO(),
				d:   &schema.ResourceData{},
				fn: func(ctx context.Context, d *schema.ResourceData) (bool, error) {
					return true, nil
				},
			},
			wantErr: false,
		},
		{
			name: "TestTimeoutOnReturnFalse",
			args: args{
				ctx: timeoutCtx,
				d:   &schema.ResourceData{},
				fn: func(ctx context.Context, d *schema.ResourceData) (bool, error) {
					return false, nil
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WaitForResourceToBeReady(tt.args.ctx, tt.args.d, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("WaitForResourceToBeReady() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
