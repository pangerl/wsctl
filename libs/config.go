// Package libs @Author lanpang
// @Date 2024/8/26 下午2:24:00
// @Desc
package libs

import (
	"github.com/jackc/pgx/v5"
	"github.com/nsf/termbox-go"
)

type DB struct {
	Ip       string `toml:"ip"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Sslmode  bool   `toml:"sslmode"`
}

type PGClienter struct {
	Conn map[string]*pgx.Conn
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Inventory map[string]Animation

type Animation struct {
	Metadata map[string]string
	Frames   [][]byte
}

var colors = []termbox.Attribute{
	// approx colors from original gif
	termbox.Attribute(210), // peach
	termbox.Attribute(222), // orange
	termbox.Attribute(120), // green
	termbox.Attribute(123), // cyan
	termbox.Attribute(111), // blue
	termbox.Attribute(134), // purple
	termbox.Attribute(177), // pink
	termbox.Attribute(207), // fuschia
	termbox.Attribute(206), // magenta
	termbox.Attribute(204), // red
}
