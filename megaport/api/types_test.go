package api

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProductAssociatedVxcResources_UnmarshalJSON(t *testing.T) {
	tc := []struct {
		in  string
		out ProductAssociatedVxcResources
	}{
		{
			`{"csp_connection":{"connectType":"AWS"}}`,
			ProductAssociatedVxcResources{CspConnection: []CspConnection{
				&ProductAssociatedVxcResourcesCspConnectionAws{ConnectType: VxcConnectTypeAws},
			}},
		},
		{
			`{"csp_connection":[{"connectType":"AWS"},{"connectType":"GOOGLE"}]}`,
			ProductAssociatedVxcResources{CspConnection: []CspConnection{
				&ProductAssociatedVxcResourcesCspConnectionAws{ConnectType: VxcConnectTypeAws},
				&ProductAssociatedVxcResourcesCspConnectionGcp{ConnectType: VxcConnectTypeGoogle},
			}},
		},
	}
	for i, test := range tc {
		v := ProductAssociatedVxcResources{}
		if err := json.Unmarshal([]byte(test.in), &v); err != nil {
			t.Errorf("TestProduct_UnmarshalJSON: unexpected error in test case %d: %v", i, err)
		}
		if diff := cmp.Diff(test.out, v); diff != "" {
			t.Errorf("TestProduct_UnmarshalJSON: unexpected result in test case %d:\n%s", i, diff)
		}
	}
	if err := json.Unmarshal([]byte(`{"csp_connection":{}}`), &ProductAssociatedVxcResources{}); err == nil {
		t.Errorf("TestProduct_UnmarshalJSON: expected an error but did not get one")
	}
	if err := json.Unmarshal([]byte(`{"csp_connection":"foo"}`), &ProductAssociatedVxcResources{}); err == nil {
		t.Errorf("TestProduct_UnmarshalJSON: expected an error but did not get one")
	}
}
