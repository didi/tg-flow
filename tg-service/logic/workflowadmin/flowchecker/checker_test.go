package flowchecker

import (
	"reflect"
	"testing"
	"tg-service/idl"
)

func TestCheckAndConvertID(t *testing.T) {
	type args struct {
		flowChart  *idl.WorkflowChart
		workflowId int64
	}
	tests := []struct {
		name    string
		args    args
		want    *idl.WorkflowChart
		wantErr bool
	}{
		{
			name: "flow char is nil",
			args: args{
				flowChart:  nil,
				workflowId: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "workflowId illegal",
			args: args{
				flowChart:  &idl.WorkflowChart{ActionMap: map[string]*idl.Action{"action-1-1": {}}},
				workflowId: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "flow char actions emtpy",
			args: args{
				flowChart:  &idl.WorkflowChart{ActionMap: map[string]*idl.Action{}},
				workflowId: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "flow char actions id illegal",
			args: args{
				flowChart:  &idl.WorkflowChart{ActionMap: map[string]*idl.Action{"foo-123": {}}},
				workflowId: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repeat action uid",
			args: args{
				flowChart: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
					"action-103-1": {ActionId: "action-103-1"},
					"action-104-1": {ActionId: "action-104-1"},
					"action-103-3": {ActionId: "action-103-3"},
					"action-103-4": {ActionId: "action-103-4"},
					"action-103-5": {ActionId: "action-103-5"},
				}},
				workflowId: 100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "flow char",
			args: args{
				flowChart: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
					"action-103-1": {ActionId: "action-103-1"},
					"action-103-2": {ActionId: "action-103-2"},
					"action-103-3": {ActionId: "action-103-3"},
					"action-103-4": {ActionId: "action-103-4"},
					"action-103-5": {ActionId: "action-103-5"},
				}},
				workflowId: 100,
			},
			want: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
				"action-100-1": {ActionId: "action-100-1"},
				"action-100-2": {ActionId: "action-100-2"},
				"action-100-3": {ActionId: "action-100-3"},
				"action-100-4": {ActionId: "action-100-4"},
				"action-100-5": {ActionId: "action-100-5"},
			}},
			wantErr: false,
		},
		{
			name: "action next ids",
			args: args{
				flowChart: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
					"action-103-1": {ActionId: "action-103-1", NextActionIds: []string{"action-103-3"}},
					"action-103-2": {ActionId: "action-103-2", NextActionIds: []string{"action-103-4"}},
					"action-103-3": {ActionId: "action-103-3", NextActionIds: []string{"action-103-4", "action-103-5"}},
					"action-103-4": {ActionId: "action-103-4", NextActionIds: []string{"action-103-5"}},
					"action-103-5": {ActionId: "action-103-5"},
				}},
				workflowId: 100,
			},
			want: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
				"action-100-1": {ActionId: "action-100-1", NextActionIds: []string{"action-100-3"}},
				"action-100-2": {ActionId: "action-100-2", NextActionIds: []string{"action-100-4"}},
				"action-100-3": {ActionId: "action-100-3", NextActionIds: []string{"action-100-4", "action-100-5"}},
				"action-100-4": {ActionId: "action-100-4", NextActionIds: []string{"action-100-5"}},
				"action-100-5": {ActionId: "action-100-5"},
			}},
			wantErr: false,
		},
		{
			name: "action next ids duplicated",
			args: args{
				flowChart: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
					"action-103-1": {ActionId: "action-103-1", NextActionIds: []string{"action-103-3"}},
					"action-103-2": {ActionId: "action-103-2", NextActionIds: []string{"action-103-4"}},
					"action-103-3": {ActionId: "action-103-3", NextActionIds: []string{"action-103-4", "action-103-4"}},
					"action-103-4": {ActionId: "action-103-4", NextActionIds: []string{"action-103-5"}},
					"action-103-5": {ActionId: "action-103-5"},
				}},
				workflowId: 100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "action next ids not exits",
			args: args{
				flowChart: &idl.WorkflowChart{ActionMap: map[string]*idl.Action{
					"action-103-1": {ActionId: "action-103-1", NextActionIds: []string{"action-103-9"}},
					"action-103-2": {ActionId: "action-103-2", NextActionIds: []string{"action-103-4"}},
					"action-103-3": {ActionId: "action-103-3", NextActionIds: []string{"action-103-4", "action-103-5"}},
					"action-103-4": {ActionId: "action-103-4", NextActionIds: []string{"action-103-5"}},
					"action-103-5": {ActionId: "action-103-5"},
				}},
				workflowId: 100,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceActionIDsForFlow(tt.args.flowChart, tt.args.workflowId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceActionIDsForFlow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceActionIDsForFlow() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
