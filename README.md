<div align="center">
    <img src="./images/strangelet.png" width="400" height="400" />
</div>

# Strangelet
Terminal-based text editor written in golang whose code base is largely adapted from micro with minor infulence from Vim/Neovim.
This project's primary goals are to create a terminal-based text editor that is cross-platform and extensible with an efficient workflow.

## Specific code editor features intended include:

- [x] Easy to use and install
- [x] Portable binary
- [x] Common keybindings (Ctrl-s, Ctrl-c, Ctrl-v, Ctrl-z, …)
- [x] Customizable Keybinds
- [x] Customizable Theme
- [x] Integrated File Explorer
- [x] Vertical and Horizontal Splits
- [x] Tabs
- [ ] Mouse support
- [ ] Multiple cursors
- [ ] Cross-platform
- [ ] Goto Definition
- [ ] Goto File
- [ ] Git Gutter
- [ ] Code Folding
- [ ] Line Numbers
- [ ] Undo/Redo
- [x] Single Instance (via flag)
- [ ] Command Integration (similar to vim's `:<command>`)
- [ ] Copy and paste with the system clipboard
- [ ] Plugin Support
    - [ ] Lua
- [ ] Auto-completion
- [x] Syntax highlighting
and more

This editor is built as a learning process just as much as it is built out of desire for a more feature rich version of micro.

## Building from Source

If your system can run Go, you can build from source.

Make sure that you have Go version 1.19 or greater and Go modules are enabled.

```
git clone https://github.com/kyrasuum/strangelet
cd strangelet
make
sudo mv strangelet /usr/local/bin # optional
```

The binary will be placed in the current directory and can be moved to anywhere you like (for example /usr/local/bin).

## Contributing

If you have features/bugs please open an issue or PR.

I am open to accepting pull requests from anyone.
