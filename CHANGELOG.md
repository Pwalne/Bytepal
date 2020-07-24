# Library Changes

## 0.1.7
  * Added Payload function to reader
  * Added ReadUMedium func to reader
  * Added ReadBigSmart func to reader

## 0.1.5
  * Added a Seek function for setting the index of the reader

## 0.1.4
  * Added a readByte method to read into a given byte array.
  * Renamed ReadInt to ReadUInt32

## 0.1.3
  * Reader now has  Read(size) func for reading a slice of the payload

## 0.1.2
  * Added bool to byte functions

## 0.1.1
  * Fixed issue with WriteBit logic
  * Converted tests to use testify/assert