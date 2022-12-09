package cc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	once       sync.Once
	handlers   []HandlerFunc
	index      int
	db         *sql.DB
	dialect    dialect
}

func newContext(w http.ResponseWriter, r *http.Request, db *sql.DB, dialect dialect) *Context {
	return &Context{
		Writer:  w,
		Req:     r,
		Method:  r.Method,
		Path:    r.URL.Path,
		index:   -1,
		db:      db,
		dialect: dialect,
	}
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...any) {
	c.once.Do(func() {
		c.SetHeader("Content-Type", "text/plain; charset=utf-8")
		c.SetStatus(code)
		c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	})
}

func (c *Context) HTML(code int, value string) {
	c.once.Do(func() {
		c.SetHeader("Content-Type", "text/html; charset=utf-8")
		c.SetStatus(code)
		c.Writer.Write([]byte(value))
	})
}

func (c *Context) JSON(code int, value any) {
	c.once.Do(func() {
		c.SetHeader("Content-Type", "application/json; charset=utf-8")
		c.SetStatus(code)
		encoder := json.NewEncoder(c.Writer)
		if err := encoder.Encode(value); err != nil {
			panic(err.Error())
		}
	})
}

func (c *Context) Data(code int, value []byte) {
	c.once.Do(func() {
		c.SetStatus(code)
		c.Writer.Write(value)
	})
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

func (c *Context) AbortWithStatus(code int) {
	c.StatusCode = code
	c.SetStatus(code)
	c.index = len(c.handlers)
}

func (c *Context) AbortWithData(code int, value []byte) {
	c.AbortWithStatus(code)
	c.Writer.Write(value)
}

func (c *Context) Next() {
	s := len(c.handlers)
	for c.index++; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) NewSession() *session {
	if c.db == nil {
		panic("database connection not found")
	}
	return newSession(c.db, c.dialect)
}

func (c *Context) GetIP() (string, error) {
	ip := c.Req.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	ip = c.Req.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}
	ip, _, err := net.SplitHostPort(c.Req.RemoteAddr)
	if err != nil {
		return "", err
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}
