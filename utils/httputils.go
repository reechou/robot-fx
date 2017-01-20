package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/errors"
	"github.com/reechou/robot-fx/version"
	"golang.org/x/net/context"
)

const APIVersionKey = "api-version"

// APIFunc is an adapter to allow the use of ordinary functions as API server endpoints.
// Any function that has the appropriate signature can be register as a API endpoint (e.g. getVersion).
type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error

// HijackConnection interrupts the http response writer to get the
// underlying connection and operate with it.
func HijackConnection(w http.ResponseWriter) (io.ReadCloser, io.Writer, error) {
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		return nil, nil, err
	}
	// Flush the options to make sure the client sets the raw mode
	conn.Write([]byte{})
	return conn, conn, nil
}

// CloseStreams ensures that a list for http streams are properly closed.
func CloseStreams(streams ...interface{}) {
	for _, stream := range streams {
		if tcpc, ok := stream.(interface {
			CloseWrite() error
		}); ok {
			tcpc.CloseWrite()
		} else if closer, ok := stream.(io.Closer); ok {
			closer.Close()
		}
	}
}

// CheckForJSON makes sure that the request's Content-Type is application/json.
func CheckForJSON(r *http.Request) error {
	ct := r.Header.Get("Content-Type")

	if ct == "" {
		if r.Body == nil || r.ContentLength == 0 {
			return nil
		}
	}

	if MatchesContentType(ct, "application/json") {
		return nil
	}
	return fmt.Errorf("Content-Type specified (%s) must be 'application/json'", ct)
}

// ParseForm ensures the request form is parsed even with invalid content types.
func ParseForm(r *http.Request) error {
	if r == nil {
		return nil
	}
	if err := r.ParseForm(); err != nil && !strings.HasPrefix(err.Error(), "mime:") {
		return err
	}
	return nil
}

// ParseMultipartForm ensure the request form is parsed, even with invalid content types.
func ParseMultipartForm(r *http.Request) error {
	if err := r.ParseMultipartForm(4096); err != nil && !strings.HasPrefix(err.Error(), "mime:") {
		return err
	}
	return nil
}

// WriteError decodes a specific api server error and sends it in the response.
func WriteError(w http.ResponseWriter, err error) {
	if err == nil || w == nil {
		logrus.WithFields(logrus.Fields{"error": err, "writer": w}).Error("unexpected HTTP error handling")
		return
	}

	statusCode := http.StatusInternalServerError
	errMsg := err.Error()

	switch err.(type) {
	case errors.ErrorCode:
		daError, _ := err.(errors.ErrorCode)
		statusCode = daError.Descriptor().HTTPStatusCode
		errMsg = daError.Message()

	case errors.Error:

		daError, _ := err.(errors.Error)
		statusCode = daError.ErrorCode().Descriptor().HTTPStatusCode
		errMsg = daError.Message

	default:
		// FIXME: this is brittle and should not be necessary.
		errStr := strings.ToLower(err.Error())
		for keyword, status := range map[string]int{
			"not found":             http.StatusNotFound,
			"no such":               http.StatusNotFound,
			"bad parameter":         http.StatusBadRequest,
			"conflict":              http.StatusConflict,
			"impossible":            http.StatusNotAcceptable,
			"wrong login/password":  http.StatusUnauthorized,
			"hasn't been activated": http.StatusForbidden,
		} {
			if strings.Contains(errStr, keyword) {
				statusCode = status
				break
			}
		}
	}

	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, errMsg, statusCode)
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteXML(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	return xml.NewEncoder(w).Encode(v)
}

// VersionFromContext returns an API version from the context using APIVersionKey.
// It panics if the context value does not have version.Version type.
func VersionFromContext(ctx context.Context) (ver version.Version) {
	if ctx == nil {
		return
	}
	val := ctx.Value(APIVersionKey)
	if val == nil {
		return
	}
	return val.(version.Version)
}

// MatchesContentType validates the content type against the expected one
func MatchesContentType(contentType, expectedType string) bool {
	mimetype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		logrus.Errorf("Error parsing media type: %s error: %v", contentType, err)
	}
	return err == nil && mimetype == expectedType
}
