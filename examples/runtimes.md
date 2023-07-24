| Program | #routines | #elems | log size | Without Tracer | With Tracer - Without Print | With Tracer - With Print | 
|-|-|-|-|-|-|-|
| gocrawl | 33 | 16199 | 1.2 MB | 2.7 &pm; 0.3 s   | 2.7 &pm; 0.4 s<br>(0.39 &pm; 0.18 %) | 2.9 &pm; 0.40 s<br>(8.63 &pm; 0.18) | 
| htcat   | 51 | 81254 | 6.5 MB | 3.7 &pm; 0.8 s   | 3.6 &pm; 0.7 s<br>(-0.4 &pm; 0.3 %)  | 4.0 &pm; 0.3 s<br>(9.44 &pm; 0.26 %)  |
| pgzip   | 1557 | 1139736 | 91 MB | 6.97 &pm; 0.04 s | 7.26 &pm; 0.06 s<br>(4.165 &pm; 0.010 %) | 20.48 &pm; 0.13 s<br>(193.913 &pm; 0.025 %) | 
| sorty   | 199 | 30001499 | 2.4 GB | 0.432 &pm; 0.010 s | 27.0 &pm; 1.0 s<br>(4140.2 &pm; 2.71 %)  | 410 &pm; 9<br>(95420 &pm; 30 %) |
| bbolt bench | 17 | 1247869 | 99.7 MB | 1.015 &pm; 0.027 s | 1.031 &pm; 0.013 s<br>(1.588 &pm; 0.012 %) | 15.23 &pm; 0.09<br>(1401.42 &pm; 0.09 %)|
| worst* 2000 | 2014 | 40022 | 3 MB | 2.269 &pm; 0.020 s | 2.337 &pm; 0.028 s<br>(3.012 &pm; 0.015 %) | 2.46 &pm; 0.048 s<br>(8.559 &pm; 0.023 %) |
| worst* 10000 | 10015 | 200026 | 15 MB | 11.320 &pm; 0.028 s | 12.48 &pm; 0.05 s<br>(10.247 &pm; 0.005 %) | 13.89 &pm; 0.06 s<br>(22.703 &pm; 0.006 %) |
| worst* 100000 | 100015 | 2000049 | 149 MB | 111.5 &pm; 0.9 s | 113.5 &pm; 0.5 s<br>(1.811 &pm; 0.009 %) | 118.0 &pm; 0.5 s<br>(5.863 &pm; 0.010 %)|

*: worst is an artificial program containing almost only concurrency objects (channel, select, mutex, ...). It thus represents a worst case scenario for the runtime increase.