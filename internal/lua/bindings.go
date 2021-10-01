package lua

import (
	"log"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"strangelet/internal/util"
)

func init() {
	L = lua.NewState()
	L.SetGlobal("import", luar.New(L, LuaImport))
}

// LuaImport is meant to be called from lua by a plugin and will import the given micro package
func LuaImport(pkg string) *lua.LTable {
	switch pkg {
	case "micro":
		return luaImportMicro()
	case "micro/shell":
		return luaImportMicroShell()
	case "micro/buffer":
		return luaImportMicroBuffer()
	case "micro/config":
		return luaImportMicroConfig()
	case "micro/util":
		return luaImportMicroUtil()
	default:
		return Import(pkg)
	}
}

func luaImportMicro() *lua.LTable {
	pkg := L.NewTable()

	// L.SetField(pkg, "TermMessage", luar.New(L, screen.TermMessage))
	// L.SetField(pkg, "TermError", luar.New(L, screen.TermError))
	// L.SetField(pkg, "InfoBar", luar.New(L, action.GetInfoBar))
	L.SetField(pkg, "Log", luar.New(L, log.Println))
	// L.SetField(pkg, "SetStatusInfoFn", luar.New(L, display.SetStatusInfoFnLua))
	// L.SetField(pkg, "CurPane", luar.New(L, func() action.Pane {
	// return action.MainTab().CurPane()
	// }))
	// L.SetField(pkg, "CurTab", luar.New(L, action.MainTab))
	// L.SetField(pkg, "Tabs", luar.New(L, func() *action.TabList {
	// return action.Tabs
	// }))
	L.SetField(pkg, "Lock", luar.New(L, Lock))

	return pkg
}

func luaImportMicroConfig() *lua.LTable {
	pkg := L.NewTable()

	// L.SetField(pkg, "MakeCommand", luar.New(L, action.MakeCommand))
	// L.SetField(pkg, "FileComplete", luar.New(L, buffer.FileComplete))
	// L.SetField(pkg, "HelpComplete", luar.New(L, action.HelpComplete))
	// L.SetField(pkg, "OptionComplete", luar.New(L, action.OptionComplete))
	// L.SetField(pkg, "OptionValueComplete", luar.New(L, action.OptionValueComplete))
	L.SetField(pkg, "NoComplete", luar.New(L, nil))
	// L.SetField(pkg, "TryBindKey", luar.New(L, action.TryBindKey))
	// L.SetField(pkg, "Reload", luar.New(L, action.ReloadConfig))
	// L.SetField(pkg, "AddRuntimeFileFromMemory", luar.New(L, config.PluginAddRuntimeFileFromMemory))
	// L.SetField(pkg, "AddRuntimeFilesFromDirectory", luar.New(L, config.PluginAddRuntimeFilesFromDirectory))
	// L.SetField(pkg, "AddRuntimeFile", luar.New(L, config.PluginAddRuntimeFile))
	// L.SetField(pkg, "ListRuntimeFiles", luar.New(L, config.PluginListRuntimeFiles))
	// L.SetField(pkg, "ReadRuntimeFile", luar.New(L, config.PluginReadRuntimeFile))
	// L.SetField(pkg, "NewRTFiletype", luar.New(L, config.NewRTFiletype))
	// L.SetField(pkg, "RTColorscheme", luar.New(L, config.RTColorscheme))
	// L.SetField(pkg, "RTSyntax", luar.New(L, config.RTSyntax))
	// L.SetField(pkg, "RTHelp", luar.New(L, config.RTHelp))
	// L.SetField(pkg, "RTPlugin", luar.New(L, config.RTPlugin))
	// L.SetField(pkg, "RegisterCommonOption", luar.New(L, config.RegisterCommonOptionPlug))
	// L.SetField(pkg, "RegisterGlobalOption", luar.New(L, config.RegisterGlobalOptionPlug))
	// L.SetField(pkg, "GetGlobalOption", luar.New(L, config.GetGlobalOption))
	// L.SetField(pkg, "SetGlobalOption", luar.New(L, action.SetGlobalOption))
	// L.SetField(pkg, "SetGlobalOptionNative", luar.New(L, action.SetGlobalOptionNative))
	// L.SetField(pkg, "ConfigDir", luar.New(L, config.ConfigDir))

	return pkg
}

func luaImportMicroShell() *lua.LTable {
	pkg := L.NewTable()

	// L.SetField(pkg, "ExecCommand", luar.New(L, shell.ExecCommand))
	// L.SetField(pkg, "RunCommand", luar.New(L, shell.RunCommand))
	// L.SetField(pkg, "RunBackgroundShell", luar.New(L, shell.RunBackgroundShell))
	// L.SetField(pkg, "RunInteractiveShell", luar.New(L, shell.RunInteractiveShell))
	// L.SetField(pkg, "JobStart", luar.New(L, shell.JobStart))
	// L.SetField(pkg, "JobSpawn", luar.New(L, shell.JobSpawn))
	// L.SetField(pkg, "JobStop", luar.New(L, shell.JobStop))
	// L.SetField(pkg, "JobSend", luar.New(L, shell.JobSend))
	// L.SetField(pkg, "RunTermEmulator", luar.New(L, action.RunTermEmulator))
	// L.SetField(pkg, "TermEmuSupported", luar.New(L, action.TermEmuSupported))

	return pkg
}

func luaImportMicroBuffer() *lua.LTable {
	pkg := L.NewTable()

	// L.SetField(pkg, "NewMessage", luar.New(L, buffer.NewMessage))
	// L.SetField(pkg, "NewMessageAtLine", luar.New(L, buffer.NewMessageAtLine))
	// L.SetField(pkg, "MTInfo", luar.New(L, buffer.MTInfo))
	// L.SetField(pkg, "MTWarning", luar.New(L, buffer.MTWarning))
	// L.SetField(pkg, "MTError", luar.New(L, buffer.MTError))
	// L.SetField(pkg, "Loc", luar.New(L, func(x, y int) buffer.Loc {
	// return buffer.Loc{x, y}
	// }))
	// L.SetField(pkg, "SLoc", luar.New(L, func(line, row int) display.SLoc {
	// return display.SLoc{line, row}
	// }))
	// L.SetField(pkg, "BTDefault", luar.New(L, buffer.BTDefault.Kind))
	// L.SetField(pkg, "BTHelp", luar.New(L, buffer.BTHelp.Kind))
	// L.SetField(pkg, "BTLog", luar.New(L, buffer.BTLog.Kind))
	// L.SetField(pkg, "BTScratch", luar.New(L, buffer.BTScratch.Kind))
	// L.SetField(pkg, "BTRaw", luar.New(L, buffer.BTRaw.Kind))
	// L.SetField(pkg, "BTInfo", luar.New(L, buffer.BTInfo.Kind))
	// L.SetField(pkg, "NewBuffer", luar.New(L, func(text, path string) *buffer.Buffer {
	// return buffer.NewBufferFromString(text, path, buffer.BTDefault)
	// }))
	// L.SetField(pkg, "NewBufferFromFile", luar.New(L, func(path string) (*buffer.Buffer, error) {
	// return buffer.NewBufferFromFile(path, buffer.BTDefault)
	// }))
	// L.SetField(pkg, "ByteOffset", luar.New(L, buffer.ByteOffset))
	// L.SetField(pkg, "Log", luar.New(L, buffer.WriteLog))
	// L.SetField(pkg, "LogBuf", luar.New(L, buffer.GetLogBuf))

	return pkg
}

func luaImportMicroUtil() *lua.LTable {
	pkg := L.NewTable()

	L.SetField(pkg, "RuneAt", luar.New(L, util.LuaRuneAt))
	L.SetField(pkg, "GetLeadingWhitespace", luar.New(L, util.LuaGetLeadingWhitespace))
	L.SetField(pkg, "IsWordChar", luar.New(L, util.LuaIsWordChar))
	L.SetField(pkg, "String", luar.New(L, util.String))
	L.SetField(pkg, "Unzip", luar.New(L, util.Unzip))
	L.SetField(pkg, "Version", luar.New(L, util.Version))
	L.SetField(pkg, "SemVersion", luar.New(L, util.SemVersion))
	L.SetField(pkg, "CharacterCountInString", luar.New(L, util.CharacterCountInString))
	L.SetField(pkg, "RuneStr", luar.New(L, func(r rune) string {
		return string(r)
	}))

	return pkg
}
