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

# Example

````
$ diff -u test.md <(./markdown-sql-format -filename test.md 2> /dev/null)
--- test.md	2024-01-19 09:48:08.396747750 +0100
+++ /dev/fd/63	2024-01-19 11:25:24.462083985 +0100
@@ -9,33 +9,33 @@
 Example 2:
 
 ```sql
-SeLeCt 2;
+SELECT 2;
 ```
 
 Example 3:
 
 ```sql
-mysql> SeLeCt 3;
+mysql> SELECT 3;
 ```
 
 Example 4:
 
 ```sql
-SeLeCt 1    Union aLL Select 4;
+SELECT 1 UNION ALL SELECT 4;
 ```
 
 Example 5:
 
 ```sql
-INSERT INTO t(c) VALUES (1);
-iNSERT INTO t(c) VALUES (2);
-InSERT INTO t(c) VALUES (3), (4), (5);
+INSERT INTO `t` (`c`) VALUES (1);
+INSERT INTO `t` (`c`) VALUES (2);
+INSERT INTO `t` (`c`) VALUES (3),(4),(5);
 ```
 
 Example 6:
 
 ```sql
-mysql> SelECT * FROM t;
+mysql> SELECT * FROM `t`;
 +----+---+
 | id | c |
 +----+---+
@@ -46,4 +46,4 @@
 | 5  | 5 |
 +----+---+
 5 rows in set (0.01 sec)
-```
\ No newline at end of file
+```
````
