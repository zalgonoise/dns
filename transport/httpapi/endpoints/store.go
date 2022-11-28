package endpoints

import (
	"io"
	"net/http"

	"github.com/zalgonoise/dns/store"
	"github.com/zalgonoise/dns/transport/httpapi"
	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/attr"
)

func (e *endpoints) AddRecord(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:records:add",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidBody.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidBody.Error(), attr.String("error", err.Error()))
		return
	}

	record := &store.Record{}
	err = e.enc.Decode(b, record)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidJSON.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidJSON.Error(), attr.String("error", err.Error()))
		return
	}

	err = e.s.AddRecord(ctx, record)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	out, err := e.s.GetRecordByTypeAndDomain(ctx, record.Type, record.Name)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: true,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.StoreResponse{
		Success: true,
		Message: "added record successfully",
		Record:  out,
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("added record successfully")
}
func (e *endpoints) ListRecords(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:records:list",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)

	records, err := e.s.ListRecords(ctx)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	var out = make([]store.Record, len(records))
	for idx, record := range records {
		out[idx] = *record
	}

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.StoreResponse{
		Success: true,
		Message: "listing all records",
		Records: &out,
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("listed all records", attr.Int("len", len(out)))
}
func (e *endpoints) GetRecordByDomain(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:records:get",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidBody.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidBody.Error(), attr.String("error", err.Error()))
		return
	}

	record := &store.Record{}
	err = e.enc.Decode(b, record)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidJSON.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidJSON.Error(), attr.String("error", err.Error()))
		return
	}

	out, err := e.s.GetRecordByTypeAndDomain(ctx, record.Type, record.Name)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: true,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.StoreResponse{
		Success: true,
		Message: "fetched record for domain " + record.Name,
		Record:  out,
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("fetched record for domain", attr.String("domain", record.Name))
}

func (e *endpoints) GetRecordByAddress(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:records:list",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidBody.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidBody.Error(), attr.String("error", err.Error()))
		return
	}

	record := &store.Record{}
	err = e.enc.Decode(b, record)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidJSON.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidJSON.Error(), attr.String("error", err.Error()))
		return
	}

	records, err := e.s.GetRecordByAddress(ctx, record.Addr)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: true,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	var out = make([]store.Record, len(records))
	for idx, record := range records {
		out[idx] = *record
	}

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.StoreResponse{
		Success: true,
		Message: "listing all records for IP address " + record.Addr,
		Records: &out,
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("listed records for address", attr.String("address", record.Addr), attr.Int("len", len(out)))
}
func (e *endpoints) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:records:update",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidBody.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidBody.Error(), attr.String("error", err.Error()))
		return
	}

	rwt := &store.RecordWithTarget{}
	err = e.enc.Decode(b, rwt)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidJSON.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidJSON.Error(), attr.String("error", err.Error()))
		return
	}

	err = e.s.UpdateRecord(ctx, rwt.Target, &rwt.Record)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: true,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}
	out, err := e.s.GetRecordByTypeAndDomain(ctx, rwt.Record.Type, rwt.Record.Name)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: true,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.StoreResponse{
		Success: true,
		Message: "updated record successfully",
		Record:  out,
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("updated record successfully")
}
func (e *endpoints) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	ctx := e.newCtx("http:records:delete",
		attr.String("remote_addr", r.RemoteAddr),
		attr.String("user_agent", r.UserAgent()),
	)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidBody.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidBody.Error(), attr.String("error", err.Error()))
		return
	}

	record := &store.Record{}
	err = e.enc.Decode(b, record)
	if err != nil {
		w.WriteHeader(400)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: false,
			Message: httpapi.ErrInvalidJSON.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Warn(httpapi.ErrInvalidJSON.Error(), attr.String("error", err.Error()))
		return
	}

	err = e.s.DeleteRecord(ctx, record)
	if err != nil {
		w.WriteHeader(500)
		response, _ := e.enc.Encode(httpapi.StoreResponse{
			Success: true,
			Message: httpapi.ErrInternal.Error(),
			Error:   err.Error(),
		})
		_, _ = w.Write(response)
		logx.From(ctx).Error(httpapi.ErrInternal.Error(), attr.String("error", err.Error()))
		return
	}

	w.WriteHeader(200)
	response, _ := e.enc.Encode(httpapi.StoreResponse{
		Success: true,
		Message: "record deleted successfully",
	})
	_, _ = w.Write(response)
	logx.From(ctx).Debug("record deleted successfully")
}
