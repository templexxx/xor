# XOR

XOR code engine in pure Go

more than 14GB/S per core

## Introduction:

1. ARCH: amd64, arm64
2. Go version:
3. Use SIMD for speeding up ( only supported in amd64 with avx2 or sse2)
4. ...

## Installation
To get the package use the standard:
```bash
go get github.com/templexxx/xor
```

## Usage

### API

**only two:**

1.
```
func Matrix(dst []byte, src [][]byte) (err error)
```
2.
```
func Bytes(dst, src1, src2 []byte) (err error)
```

## Performance

Performance depends mainly on:

1. SIMD extension
2. unit size of worker
3. hardware ( CPU RAM etc)

Example of performance on my MacBook 2014-mid(i5-4278U 2.6GHz 2 physical cores). The 16MB per shards.
```
speed = ( shards * size ) / cost
```
| data_shards    | shard_size |speed (MB/S) |
|----------|----|-----|
| 2       |1KB|13073.46  |
|2|1400B||
|2|16KB||
| 2       | 16MB|14016.86 |
| 5       |1KB| 14109.60 |
|5|1400B||
|5|16KB||
|5| 16MB||
