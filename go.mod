module strangelet

go 1.17

require (
	code.rocketnine.space/tslocum/cbind v0.1.5 // indirect
	github.com/Kyrasuum/cview v1.5.8-0.20210925024824-4ac3d4c57ab6
	github.com/blang/semver v3.5.1+incompatible
	github.com/dustin/go-humanize v1.0.0
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.4.1-0.20210828201608-73703f7ed490
	github.com/go-errors/errors v1.4.1
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14
	github.com/mattn/go-runewidth v0.0.14-0.20210830053702-dc8fe66265af
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sergi/go-diff v1.2.0
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	github.com/zyedidia/glob v0.0.0-20170209203856-dd4023a66dc3
	github.com/zyedidia/json5 v0.0.0-20200102012142-2da050b1a98d
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7
	gopkg.in/yaml.v2 v2.2.4
	layeh.com/gopher-luar v1.0.10
)

replace github.com/Kyrasuum/cview => ./external/cview
