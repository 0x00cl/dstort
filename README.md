# dstort

An utility to "organize" files randomly in a directory!

Pronounced like _distort_.

## Why

Well, thats interesting, you know why? Because ... 🏃

## Building
```bash
git clone https://github.com/0x00cl/dstort.git
cd dstort
go build
# Resulting binary is called dstort
./dstort -v
```

## Usage

```bash
$ ./dstort /home/user/directory
```

### Example
```shell
[~]$ tree test_directory 
test_directory
├── foo
│   └── bar
│       └── foobar.txt
└── hello.txt

3 directories, 2 files

[~/dstort]$ ./dstort ~/test_directory

# Result
[~]$ tree test_directory 
test_directory
├── bar
│   └── hello.txt
├── foo
└── foobar.txt

3 directories, 2 files
```
