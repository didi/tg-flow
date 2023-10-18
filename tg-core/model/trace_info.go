package model

import (

)

type TraceInfo struct {
    TraceId string `thrift:"traceId,1,required" json:"traceId"`
    Caller string `thrift:"caller,2,required" json:"caller"`
    SpanId string `thrift:"spanId,3" json:"spanId,omitempty"`
    SrcMethod string `thrift:"srcMethod,4" json:"srcMethod,omitempty"`
    HintCode int64 `thrift:"hintCode,5" json:"hintCode,omitempty"`
    HintContent string `thrift:"hintContent,6" json:"hintContent,omitempty"`
}