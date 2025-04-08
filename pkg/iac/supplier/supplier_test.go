package supplier

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/snyk/driftctl/enumeration/alerter"
	"github.com/snyk/driftctl/enumeration/terraform"
	"github.com/snyk/driftctl/pkg/resource"
	testresource "github.com/snyk/driftctl/test/resource"

	"github.com/snyk/driftctl/pkg/filter"
	"github.com/snyk/driftctl/pkg/iac/config"
	"github.com/snyk/driftctl/pkg/iac/terraform/state/backend"
	"github.com/snyk/driftctl/pkg/output"
)

func TestGetIACSupplier(t *testing.T) {
	type args struct {
		config  []config.SupplierConfig
		options *backend.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "test unknown supplier",
			args: args{
				config: []config.SupplierConfig{
					{
						Key: "foobar",
					},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr: fmt.Errorf("Unsupported supplier 'foobar'"),
		},
		{
			name: "test unknown supplier in multiples states",
			args: args{
				config: []config.SupplierConfig{
					{
						Key: "foobar",
					},
					{
						Key:     "tfstate",
						Backend: "",
						Path:    "terraform.tfstate",
					},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr: fmt.Errorf("Unsupported supplier 'foobar'"),
		},
		{
			name: "test valid tfstate://terraform.tfstate",
			args: args{
				config: []config.SupplierConfig{
					{Key: "tfstate", Backend: "", Path: "terraform.tfstate"},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr: nil,
		},
		{
			name: "test valid multiples states",
			args: args{
				config: []config.SupplierConfig{
					{Key: "tfstate", Backend: "", Path: "terraform.tfstate"},
					{Key: "tfstate", Backend: "s3", Path: "terraform.tfstate"},
					{Key: "tfstate", Backend: "", Path: "terraform2.tfstate"},
				},
				options: &backend.Options{
					Headers: map[string]string{},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := &output.MockProgress{}
			progress.On("Start").Return().Times(1)

			repo := testresource.InitFakeSchemaRepository(terraform.AWS, "5.94.1")
			factory := resource.NewDriftctlResourceFactory(repo)
			alerter := alerter.NewAlerter()

			testFilter := &filter.MockFilter{}

			_, err := GetIACSupplier(tt.args.config, terraform.NewProviderLibrary(), tt.args.options, progress, alerter, factory, testFilter)

			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetIACSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetSupportedSchemes(t *testing.T) {

	want := []string{
		"tfstate://",
		"tfstate+s3://",
		"tfstate+http://",
		"tfstate+https://",
		"tfstate+tfcloud://",
		"tfstate+gs://",
		"tfstate+azurerm://",
	}

	if got := GetSupportedSchemes(); !reflect.DeepEqual(got, want) {
		t.Errorf("GetSupportedSchemes() = %v, want %v", got, want)
	}
}
