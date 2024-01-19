# Building

```
go build
```

# Testing

```
go test -v
```

# Running

```
go build 
./markdown-sql-format -filename test.md
diff -u test.md <(./markdown-sql-format -filename test.md 2> /dev/null)
```
