# rubberRing

## Overview
This library provides implementation of a growing circular buffer.

There are two buffer:
- thread-unsafe RubberRing
- thread-safe SyncRubberRing

### Features

- no evacuations
- excess memory is released (buffer compression occurs)
- flexible configuration of allocated space and memory management

## Quick Start

### RubberRing Buffer

```go
rr := rubberring.NewRubberRing[int]( // create buffer
	rubberring.WithStartChankCount(2), // number of initially initiated chunks
	rubberring.WithStartChankSize(2), // size of initially initiated chunks
	rubberring.WithGrowStrategy(func(_ int) (int, int) { return 2, 2 }), // growth function
	rubberring.WithFreeChankBufferSize(2), // buffer size of passive chunks
)

rr.Push(1) // write to buffer
val, err := rr.Pull() // read from buffer
```

### SyncRubberRing Buffer

```go
rr := rubberring.NewSyncRubberRing[int]( // create buffer
	rubberring.WithStartChankCount(2), // number of initially initiated chunks
	rubberring.WithStartChankSize(2), // size of initially initiated chunks
	rubberring.WithGrowStrategy(func(_ int) (int, int) { return 2, 2 }), // growth function
	rubberring.WithFreeChankBufferSize(2), // buffer size of passive chunks
)

rr.Push(1) // write to buffer
val, err := rr.Pull(ctx) // read from buffer
```

## How it works

In fact, the buffer contains an array divided into chunks and a buffer of adjustable size for passive chunks.

Every time a chunk from the beginning of the buffer is released, it is transferred to the buffer for passive chunks. If the buffer is already full, the chunk is utilized (by the garbage collector). This is how the buffer is compressed.

Every time the last chunk is filled, the buffer tries to load the chunk from the buffer of passive chunks, if it is empty, new chunks are created.

![Schem](docs/scheme.drawio.svg)

## Description

### Constructor options

- `WithStartChankCount(int)` - the number of chunks created when initializing the buffer (default 4)
- `WithStartChankSize(int)` - the size of chunks created when initializing the buffer (default 256)
- `WithPassiveChankBufferSize(int)` - the size of the passive chunk buffer (default 3)
- `WithGrowStrategy(func(currentCapacity int) (newChankSize, newChankCount int))` - a function describing the size and number of chunks created when the buffer is full

By manipulating these parameters, you can customize the behavior of the buffer for different tasks.

### RubberRing Methods

- `Push(V)` - puts an element at the end of the buffer
- `Pull() (V, error)` - retrieves an element from the beginning of the buffer. If the buffer is empty, the `io.EOF` error will be returned
- `Size() int` - returns the current amount of data in the buffer
- `Capacity() int` - returns the current size of the buffer (including passive capacity)
- `Stat() RubberRingStat` - returns a detailed description of the buffer state
- `Elements() iter.Seq[V]` - returns an iterator for getting all elements of the buffer

### SyncRubberRing Methods

SyncRubberRing has the same methods as RubberRing, they work similarly (with an adjustment for thread safety) with the following exceptions
- `Pull(context.Context) (V, error)` - retrieves an element from the beginning of the buffer. If the buffer is empty - waits until at least one element appears there. If the context is closed - returns the error context.Canceled
- `Elements() iter.Seq[V]` - returns an iterator for streaming elements from the buffer. When the context is closed - the iterator will end.
