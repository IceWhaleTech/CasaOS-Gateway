package route

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"time"

	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type StaticRoute struct {
	state *service.State
}

var startTime = time.Now()

func NewStaticRoute(state *service.State) *StaticRoute {
	return &StaticRoute{
		state: state,
	}
}

type CustomFS struct {
	base fs.FS
}

func NewCustomFS(prefix string) *CustomFS {
	return &CustomFS{
		base: fs.FS(os.DirFS(prefix)),
	}
}

func (c *CustomFS) Open(name string) (fs.File, error) {
	file, err := c.base.Open(name)
	if err != nil {
		return nil, err
	}
	return &CustomFile{
		File: file,
	}, nil
}

func (c *CustomFS) Stat(name string) (fs.FileInfo, error) {
	file, err := c.base.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &CustomFileInfo{
		FileInfo: info,
	}, nil
}

type CustomFile struct {
	fs.File
}

func (c *CustomFile) Stat() (fs.FileInfo, error) {
	info, err := c.File.Stat()
	if err != nil {
		return nil, err
	}
	return &CustomFileInfo{
		FileInfo: info,
	}, nil
}

func (c *CustomFile) Read(p []byte) (int, error) {
	if seeker, ok := c.File.(io.Reader); ok {
		return seeker.Read(p)
	}
	return 0, fmt.Errorf("file does not implement io.Reader")
}

func (c *CustomFile) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := c.File.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, fmt.Errorf("file does not implement io.Seeker")
}

type CustomFileInfo struct {
	fs.FileInfo
}

func (c *CustomFileInfo) ModTime() time.Time {
	return startTime
}

var indexRE = regexp.MustCompile(`/($|modules/[^\/]*/($|(index\.(html?|aspx?|cgi|do|jsp))|((default|index|home)\.php)))`)

func (s *StaticRoute) GetRoute() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if indexRE.MatchString(ctx.Request().URL.Path) {
				ctx.Response().Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate,proxy-revalidate, max-age=0")
			}
			return next(ctx)
		}
	})

	// sovle 304 cache problem by 'If-Modified-Since: Wed, 21 Oct 2015 07:28:00 GMT' from web browser
	e.StaticFS("/", NewCustomFS(s.state.GetWWWPath()))
	return e
}
