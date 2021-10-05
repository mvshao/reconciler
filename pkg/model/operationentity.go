package model

import (
	"fmt"
	"time"

	"github.com/kyma-incubator/reconciler/pkg/db"
)

const tblOperation string = "scheduler_operations"

type OperationEntity struct {
	Priority      int64          `db:"notNull"`
	SchedulingID  string         `db:"notNull"`
	CorrelationID string         `db:"notNull"`
	Cluster       string         `db:"notNull"`
	ClusterConfig int64          `db:"notNull"`
	Component     string         `db:"notNull"`
	State         OperationState `db:"notNull"`
	Reason        string         `db:""`
	Created       time.Time      `db:"readOnly"`
	Updated       time.Time      `db:""`
}

func (o *OperationEntity) String() string {
	return fmt.Sprintf("OperationEntity [SchedulingID=%s,CorrelationID=%s,"+
		"Cluster=%s,ClusterConfig=%d,Component=%s,Prio=%d]",
		o.SchedulingID, o.CorrelationID, o.Cluster, o.ClusterConfig, o.Component, o.Priority)
}

func (*OperationEntity) New() db.DatabaseEntity {
	return &OperationEntity{}
}

func (o *OperationEntity) Marshaller() *db.EntityMarshaller {
	marshaller := db.NewEntityMarshaller(&o)
	marshaller.AddMarshaller("State", func(value interface{}) (interface{}, error) {
		return fmt.Sprintf("%s", value), nil
	})
	marshaller.AddUnmarshaller("State", func(value interface{}) (interface{}, error) {
		return NewOperationState(fmt.Sprintf("%s", value))
	})
	marshaller.AddUnmarshaller("Created", convertTimestampToTime)
	marshaller.AddUnmarshaller("Updated", convertTimestampToTime)
	return marshaller
}

func (*OperationEntity) Table() string {
	return tblOperation
}

func (o *OperationEntity) Equal(other db.DatabaseEntity) bool {
	if other == nil {
		return false
	}
	otherOpProp, ok := other.(*OperationEntity)
	if !ok {
		return false
	}
	return o.SchedulingID == otherOpProp.SchedulingID &&
		o.CorrelationID == otherOpProp.CorrelationID &&
		o.Component == otherOpProp.Component
}
