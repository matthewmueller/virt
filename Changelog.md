# 0.0.12 / 2025-01-11

- support incrementally writing data to virtual file

# 0.0.11 / 2025-01-11

- fix and rename interface
- add missing error check

# 0.0.10 / 2025-01-11

- initial support for a writable file interface.
  - Note: The `Tree` and `List` filesystems don't yet support writable files (PR welcome!)
- add support for hashing files

# 0.0.9 / 2024-12-29

More helpers for transforming to and from `fs.File` and `*virt.File`

- Breaking:
  - Renamed `From` to `FromFile`
  - Renamed `Open` to `To`
- Added:
  - `From(fsys fs.FS, path string) (*virt.File, error)`

# 0.0.8 / 2024-12-12

- skip over permission and missing file errors when printing

# 0.0.7 / 2024-12-12

- add `fsys := virt.Merge(a, b, ...)` for merging `fs.FS` filesystems together
- add `fsys := virt.Mount("some/mounted/dir", nestedFs)` for mounting a nested filesystem

# 0.0.6 / 2024-11-16

- support syncing and writing different file modes

# 0.0.5 / 2024-11-02

- breaking: require path when using virt.From, virt.FromEntry and virt.MarshalJSON

# 0.0.4 / 2024-11-02

- don't skip . when writing, it might not exist

# 0.0.3 / 2024-03-23

- make exclude work on fs.FS

# 0.0.2 / 2024-03-20

- Switch `virt.File` to use `*virt.DirEntry` instead of `fs.DirEntry`.
- Add `file.Entry()` and `virt.FromEntry(de)`
- Fix `From(file)` not sorting directory entries.
- add test confirming duplicate dir entries aren't handled. expected to handle outside

# 0.0.1 / 2024-02-21

- initial commit
