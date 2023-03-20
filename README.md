This is a golang port of the perth-3 tide calculation algorithm: https://github.com/asbjorn-christensen/GridWetData/blob/master/fortran_sources/perth3.f

The tool works with DTU16 files from the danish technical university: ftp://ftp.space.dtu.dk/pub/DTU16/OCEAN_TIDE

`./cmd/ascii2dat/main.go` will read a fort.30 file (DTU-16) for constituents and DTU15MSS_2min.mss (DTU-15) for median sea surface height and outputs pre-computed values for all the constituents and mss to a bin file.

`./cmd/test/main.go` can read the created bin file and output the tides for a specific location/time 