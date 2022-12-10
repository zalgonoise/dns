package endpoints

import (
	"net/http"

	"github.com/zalgonoise/dns/store"
)

func (e *endpoints) AddRecord(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.AddRecord")
	defer s.End()

	record, err := readBody[store.Record](ctx, r)
	if err != nil {
		res := NewResponse[store.Record](400, "failed to read record from body", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	err = e.s.AddRecord(ctx, record)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to add record", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	newRecord, err := e.s.GetRecordByTypeAndDomain(ctx, record.Type, record.Name)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to get new record", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "added record successfully", nil, newRecord)
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) ListRecords(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.ListRecords")
	defer s.End()

	records, err := e.s.ListRecords(ctx)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to list records", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "listed records successfully", nil, &records)
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) GetRecordByDomain(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.GetRecordByDomain")
	defer s.End()

	record, err := readBody[store.Record](ctx, r)
	if err != nil {
		res := NewResponse[store.Record](400, "failed to read record from body", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	output, err := e.s.GetRecordByTypeAndDomain(ctx, record.Type, record.Name)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to get record and domain", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "fetched record successfully", nil, output)
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) GetRecordByAddress(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.GetRecordByAddress")
	defer s.End()

	record, err := readBody[store.Record](ctx, r)
	if err != nil {
		res := NewResponse[store.Record](400, "failed to read record from body", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	records, err := e.s.GetRecordByAddress(ctx, record.Addr)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to get record by address", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "fetched record successfully", nil, &records)
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.UpdateRecord")
	defer s.End()

	rwt, err := readBody[store.RecordWithTarget](ctx, r)
	if err != nil {
		res := NewResponse[store.RecordWithTarget](400, "failed to read record from body", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	err = e.s.UpdateRecord(ctx, rwt.Target, &rwt.Record)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to update record", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	output, err := e.s.GetRecordByTypeAndDomain(ctx, rwt.Record.Type, rwt.Record.Name)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to get record and domain", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "updated record successfully", nil, output)
	res.WriteHTTP(ctx, w)
}

func (e *endpoints) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	ctx, s := e.newCtxAndSpan(r, "http.DeleteRecord")
	defer s.End()

	record, err := readBody[store.Record](ctx, r)
	if err != nil {
		res := NewResponse[store.Record](400, "failed to read record from body", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	err = e.s.DeleteRecord(ctx, record)
	if err != nil {
		res := NewResponse[store.Record](500, "failed to get record by address", err, nil)
		res.WriteHTTP(ctx, w)
		return
	}

	res := NewResponse(200, "deleted record successfully", nil, record)
	res.WriteHTTP(ctx, w)
}
