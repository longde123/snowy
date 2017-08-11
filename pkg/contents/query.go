package contents

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/trussle/snowy/pkg/document"
	errs "github.com/trussle/snowy/pkg/http"
	"github.com/trussle/snowy/pkg/uuid"
)

const (
	defaultContentType = "application/json"
)

const (
	defaultKB = 1024
	defaultMB = 1024 * defaultKB

	defaultMaxContentLength = 5 * defaultMB
)

// SelectQueryParams defines all the dimensions of a query.
type SelectQueryParams struct {
	ResourceID uuid.UUID `json:"resource_id"`
}

// DecodeFrom populates a SelectQueryParams from a URL.
func (qp *SelectQueryParams) DecodeFrom(u *url.URL, rb queryBehavior) error {
	// Required depending on the query behavior
	if rb == queryRequired {
		var (
			err        error
			resourceID = u.Query().Get("resource_id")
		)
		if resourceID == "" {
			return errors.New("error reading 'resource_id' (required) query")
		}
		if qp.ResourceID, err = uuid.Parse(resourceID); err != nil {
			return errors.Wrap(err, "error parsing 'resource_id' (required) query")
		}
	}

	return nil
}

// SelectQueryResult contains statistics about the query.
type SelectQueryResult struct {
	Params   SelectQueryParams `json:"query"`
	Duration string            `json:"duration"`
	Content  document.Content  `json:"content"`
}

// EncodeTo encodes the SelectQueryResult to the HTTP response writer.
func (qr *SelectQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderResourceID, qr.Params.ResourceID.String())
	w.Header().Set(httpHeaderContentType, qr.Content.ContentType())
	w.Header().Set(httpHeaderContentLength, strconv.FormatInt(qr.Content.Size(), 10))

	bytes, err := qr.Content.Bytes()
	if err != nil {
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := w.Write(bytes); err != nil {
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// InsertQueryParams defines all the dimensions of a query.
type InsertQueryParams struct {
	ContentType   string
	ContentLength int64
}

// DecodeFrom populates a InsertQueryParams from a URL.
func (qp *InsertQueryParams) DecodeFrom(u *url.URL, h http.Header, rb queryBehavior) error {
	// Required depending on the query behavior
	if rb == queryRequired {
		// Get the content-type
		if qp.ContentType = h.Get("Content-Type"); qp.ContentType == "" {
			return errors.New("error reading 'content-type' (required) query")
		}

		// Get the content-length
		contentLength := h.Get("Content-Length")
		if contentLength == "" {
			return errors.New("error reading 'content-length' (required) query")
		}

		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return errors.New("error parsing 'content-length' (required) query")
		} else if size > defaultMaxContentLength {
			return errors.Errorf("error request body too large")
		} else if size < 1 {
			return errors.Errorf("error request body is empty")
		}

		qp.ContentLength = size
	}

	return nil
}

// InsertQueryResult contains statistics about the query.
type InsertQueryResult struct {
	Params   InsertQueryParams `json:"query"`
	Duration string            `json:"duration"`
	Content  document.Content  `json:"content"`
}

// EncodeTo encodes the InsertQueryResult to the HTTP response writer.
func (qr *InsertQueryResult) EncodeTo(w http.ResponseWriter) {
	w.Header().Set(httpHeaderDuration, qr.Duration)
	w.Header().Set(httpHeaderContentType, defaultContentType)

	if err := json.NewEncoder(w).Encode(qr.Content); err != nil {
		errs.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const (
	httpHeaderDuration      = "X-Duration"
	httpHeaderResourceID    = "X-ResourceID"
	httpHeaderContentType   = "Content-Type"
	httpHeaderContentLength = "Content-Length"
)

type queryBehavior int

const (
	queryRequired queryBehavior = iota
	queryOptional
)
