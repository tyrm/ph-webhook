package models

import (
	"bytes"
	"fmt"
	"time"
)

type WebhookCall struct {
	ID            int

	Timestamp     time.Time
	Host          string
	Path          string
	Method        string
	ContentLength int
	From          string
	Body          string

	Headers  map[string]string
	Query    map[string]string
	PostForm map[string]string
	Files    map[string]WebhookFile
}

type WebhookFile struct {
	ID       int

	Filename string
	Size     int
	Contents bytes.Buffer

	Headers  map[string]string
}

const sqlCallsInsert = `
INSERT INTO "public"."webhook_calls" (timestamp, host, path, method, content_length, "from", body)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;`

const sqlCallsFilesInsert = `
INSERT INTO "public"."webhook_calls_files" (call_id, key, filename, size, contents)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;`

const sqlCallsFilesHeadersInsert = `
INSERT INTO "public"."webhook_calls_files_headers" (file_id, key, value)
VALUES ($1, $2, $3)
RETURNING id;`

const sqlCallsFormParamsInsert = `
INSERT INTO "public"."webhook_calls_form_params" (call_id, key, value)
VALUES ($1, $2, $3)
RETURNING id;`

const sqlCallsHeadersInsert = `
INSERT INTO "public"."webhook_calls_headers" (call_id, key, value)
VALUES ($1, $2, $3)
RETURNING id;`

const sqlCallsQueriesInsert = `
INSERT INTO "public"."webhook_calls_queries" (call_id, key, value)
VALUES ($1, $2, $3)
RETURNING id;`

func (c *WebhookCall) Save() error {
	if c.ID == 0 {
		id := 0
		err := DB.QueryRow(sqlCallsInsert, c.Timestamp, c.Host, c.Path, c.Method, c.ContentLength, c.From, c.Body).Scan(&id)
		if err != nil {
			logger.Errorf("Error inserting result: %s", err)
			return err
		}
		c.ID = id
		logger.Debugf("Got ID: %d", id)

		for k, v := range c.Headers {
			hid := 0
			err := DB.QueryRow(sqlCallsHeadersInsert, c.ID, k, v).Scan(&hid)
			if err != nil {
				logger.Errorf("Error inserting result: %s", err)
				return err
			}
		}
		for k, v := range c.Query {
			hid := 0
			err := DB.QueryRow(sqlCallsQueriesInsert, c.ID, k, v).Scan(&hid)
			if err != nil {
				logger.Errorf("Error inserting result: %s", err)
				return err
			}
		}
		for k, v := range c.PostForm {
			hid := 0
			err := DB.QueryRow(sqlCallsFormParamsInsert, c.ID, k, v).Scan(&hid)
			if err != nil {
				logger.Errorf("Error inserting result: %s", err)
				return err
			}
		}

		for k, file := range c.Files {
			fid := 0
			err := DB.QueryRow(sqlCallsFilesInsert, c.ID, k, file.Filename, file.Size, file.Contents.String()).Scan(&fid)
			if err != nil {
				logger.Errorf("Error inserting result: %s", err)
				return err
			}
			for k, v := range file.Headers {
				hid := 0
				err := DB.QueryRow(sqlCallsFilesHeadersInsert, fid, k, v).Scan(&hid)
				if err != nil {
					logger.Errorf("Error inserting result: %s", err)
					return err
				}
			}
		}

	}

	return nil
}

func (c *WebhookCall) GetString() (str string) {
	str = fmt.Sprintf("ID: %d\n", c.ID)
	str += fmt.Sprintf("Timestamp: %s\n", c.Timestamp)
	str += fmt.Sprintf("Host: %s\n", c.Host)
	str += fmt.Sprintf("Path: %s\n", c.Path)
	str += fmt.Sprintf("Method: %s\n", c.Method)
	str += fmt.Sprintf("Content-Length: %d\n", c.ContentLength)
	str += fmt.Sprintf("From: %s\n", c.From)

	for k, v := range c.Headers {
		str += fmt.Sprintf("Header[%s]: %s\n", k, v)
	}
	for k, v := range c.Query {
		str += fmt.Sprintf("Query[%s]: %s\n", k, v)
	}
	for k, v := range c.PostForm {
		str += fmt.Sprintf("PostForm[%s]: %s\n", k, v)
	}
	for k, v := range c.Files {
		str += fmt.Sprintf("File[%s]: Filename(%s) Size(%d)\n", k, v.Filename, v.Size)

		for sk, sv := range v.Headers {
			str += fmt.Sprintf("File[%s]Header[%s]: %s\n", k, sk, sv)
		}
		str += fmt.Sprintf("File[%s]Contents: %s\n", k, v.Contents.String())

	}

	if c.Body != "" {
		str += fmt.Sprintf("Body: %s\n", c.Body)
	}

	return
}
