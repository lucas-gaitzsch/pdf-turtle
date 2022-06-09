package models

import (
	"context"
	"errors"
	"io"
)

type Job struct {
	RequestCtx   context.Context
	RenderData   *RenderData
	CallbackChan chan io.Reader
}

func NewJob(requestCtx context.Context, renderData *RenderData) *Job {
	if renderData == nil {
		panic(errors.New("render data is nil"))
	}

	return &Job{
		RequestCtx:   requestCtx,
		RenderData:   renderData,
		CallbackChan: make(chan io.Reader, 1),
	}
}
