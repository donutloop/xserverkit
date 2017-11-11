package renderer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	// ContentType represents content type
	ContentType string = "Content-Type"
	// ContentJSON represents content type application/json
	ContentJSON string = "application/json"
	// ContentXML represents content type application/xml
	ContentXML string = "application/xml"
	// ContentYAML represents content type application/x-yaml
	ContentYAML string = "application/x-yaml"
	// ContentHTML represents content type text/html
	ContentHTML string = "text/html"
	// ContentText represents content type text/plain
	ContentText string = "text/plain"
	// ContentBinary represents content type application/octet-stream
	ContentBinary string = "application/octet-stream"

	// ContentDisposition describes contentDisposition
	ContentDisposition string = "Content-Disposition"
	// contentDispositionInline describes content disposition type
	contentDispositionInline string = "inline"
	// contentDispositionAttachment describes content disposition type
	contentDispositionAttachment string = "attachment"
	)

type (
	// Options describes an option type
	Options struct {
		// Charset represents the Response charset; default: utf-8
		Charset string
		// ContentJSON represents the Content-Type for JSON
		ContentJSON string
		// ContentXML represents the Content-Type for XML
		ContentXML string
		// ContentYAML represents the Content-Type for YAML
		ContentYAML string
		// ContentHTML represents the Content-Type for HTML
		ContentHTML string
		// ContentText represents the Content-Type for Text
		ContentText string
		// ContentBinary represents the Content-Type for octet-stream
		ContentBinary string

		// Debug set the debug mode. if debug is true then every time "VIEW" call parse the templates
		Debug bool
		// JSONIndent set JSON Indent in response; default false
		JSONIndent bool
		// XMLIndent set XML Indent in response; default false
		XMLIndent bool
	}

	// Render describes a renderer type
	Render struct {
		opts          Options
	}
)

// New return a new instance of a pointer to Render
func New(opts *Options) *Render {
	r := &Render{
		opts:      *opts,
	}

	// build options for the Render instance
	r.buildOptions()
	return r
}

// buildOptions builds the options and set deault values for options
func (r *Render) buildOptions() {
	r.opts.ContentJSON = ContentJSON
	r.opts.ContentXML = ContentXML
	r.opts.ContentYAML = ContentYAML
	r.opts.ContentHTML = ContentHTML
	r.opts.ContentText = ContentText
	r.opts.ContentBinary = ContentBinary
}


// NoContent serve success but no content response
func (r *Render) NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// String serve string content as text/plain response
func (r *Render) String(w http.ResponseWriter, status int, v string) error {
	w.Header().Set(ContentType, r.opts.ContentText)
	w.WriteHeader(status)
	_, err := w.Write([]byte(v))
	return err
}

// json converts the data as bytes using json encoder
func (r *Render) json(v interface{}) ([]byte, error) {
	var bs []byte
	var err error
	if r.opts.JSONIndent {
		bs, err = json.MarshalIndent(v, "", " ")
	} else {
		bs, err = json.Marshal(v)
	}
	if err != nil {
		return bs, err
	}
	return bs, nil
}

// JSON serve data as JSON as response
func (r *Render) JSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentJSON)
	w.WriteHeader(status)

	bs, err := r.json(v)
	if err != nil {
		return err
	}

	_, err = w.Write(bs)
	return err
}

// XML serve data as XML response
func (r *Render) XML(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentXML)
	w.WriteHeader(status)
	var bs []byte
	var err error

	if r.opts.XMLIndent {
		bs, err = xml.MarshalIndent(v, "", " ")
	} else {
		bs, err = xml.Marshal(v)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

// YAML serve data as YAML response
func (r *Render) YAML(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentYAML)
	w.WriteHeader(status)

	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

// Binary serve file as application/octet-stream response; you may add ContentDisposition by your own.
func (r *Render) Binary(w http.ResponseWriter, status int, reader io.Reader, filename string, inline bool) error {
	if inline {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionInline, filename))
	} else {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionAttachment, filename))
	}
	w.Header().Set(ContentType, r.opts.ContentBinary)
	w.WriteHeader(status)
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = w.Write(bs)
	return err
}

// File serve file as response from io.Reader
func (r *Render) File(w http.ResponseWriter, status int, reader io.Reader, filename string, inline bool) error {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	// set headers
	mime := http.DetectContentType(bs)
	if inline {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionInline, filename))
	} else {
		w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDispositionAttachment, filename))
	}
	w.Header().Set(ContentType, mime)
	w.WriteHeader(status)

	_, err = w.Write(bs)
	return err
}

// file serve file as response
func (r *Render) file(w http.ResponseWriter, status int, fpath, name, contentDisposition string) error {
	var bs []byte
	var err error
	bs, err = ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(bs)

	// filename, ext, mimes
	var fn, mime, ext string
	fn, err = filepath.Abs(fpath)
	ext = filepath.Ext(fpath)
	if name != "" {
		if !strings.HasSuffix(name, ext) {
			fn = name + ext
		}
	}

	mime = http.DetectContentType(bs)

	// set headers
	w.Header().Set(ContentType, mime)
	w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", contentDisposition, fn))
	w.WriteHeader(status)

	if _, err = buf.WriteTo(w); err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

// FileView serve file as response with content-disposition value inline
func (r *Render) FileView(w http.ResponseWriter, status int, fpath, name string) error {
	return r.file(w, status, fpath, name, contentDispositionInline)
}

// FileDownload serve file as response with content-disposition value attachment
func (r *Render) FileDownload(w http.ResponseWriter, status int, fpath, name string) error {
	return r.file(w, status, fpath, name, contentDispositionAttachment)
}
