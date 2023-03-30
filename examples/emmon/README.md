# Electricity Meter Monitor (`emmon`)

`emmon` is an exemplary usage of gosml, to output electricity meter readings (obis code `1-1:1.8.0`) parsed from an sml file. The file can also be a serial port device file like `/dev/ttyUSB0`.

## Usage

```bash
Usage: ./emmon [FILE]...
  Reads FILE(s) and outputs found electricity meter readings (obis 1.8.0)
```

For example:

```bash
$ ./emmon output.bin
1-0:1.8.0*255 26564190.300000
1-0:1.8.0*255 26564190.500000
1-0:1.8.0*255 26564190.700000
1-0:1.8.0*255 26564190.900000
1-0:1.8.0*255 26564191.100000
1-0:1.8.0*255 26564191.300000
1-0:1.8.0*255 26564191.500000
1-0:1.8.0*255 26564191.700000
```
