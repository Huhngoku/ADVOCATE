| Program | #routines | #elems | Without Tracer | With Tracer - Without Print | With Tracer - With Print | 
|-|-|-|-|-|-|
| gocrawl | 31 | 786 | 2.3 &pm; 0.3 s   | 2.37 &pm; 0.10 s<br>(3.04 &pm; 0.14 %) | 2.38 &pm; 0.10 s<br>(3.49 &pm; 0.14) | 
| htcat   | 49 | 55 | 4.8 &pm; 0.9 s   | 4.8 &pm; 0.7 s<br>(0 &pm; 0.24 %)  | 4.8 &pm; 0.7 s<br>(0 &pm; 0.24 %)  |
| pgzip   | 1555 | 15027 | 7.74 &pm; 0.14 s | 7.86 &pm; 0.09 s<br>(1.550 &pm; 0.022 %) | 8.04 &pm; 0.09 s<br>(3.876 &pm; 0.022 %) | 
| sorty   | 769 | 1477 | 4.05 &pm; 0.04 s | 4.1 &pm; 0.5 s<br>(1.23 &pm; 0.12 %)  | 4.1 &pm; 0.5<br>(1.23 &pm; 0.12 %) |
| bbolt bench | 15 | 126 | 1.023 &pm; 0.027 s | 1.023 &pm; 0.008 s<br>(0.003 &pm; 0.027 %) | 1.023 &pm; 0.005<br>(0.004 &pm; 0.027 %)|
| worst* 2000 | 2013 | 10022 | 2.269 &pm; 0.020 s | 2.337 &pm; 0.028 s<br>(3.012 &pm; 0.015 %) | 2.46 &pm; 0.048 s<br>(8.559 &pm; 0.023 %) |
| worst* 10000 | 10013 | 50020 | 11.320 &pm; 0.028 s | 11.48 &pm; 0.05 s<br>(1.421 &pm; 0.006 %) | 11.89 &pm; 0.06 s<br>(5.055 &pm; 0.006 %) |
| worst* 100000 | 100013 | 500020 | 111.5 &pm; 0.9 s | 113.5 &pm; 0.5 s<br>(1.811 &pm; 0.009 %) | 118.0 &pm; 0.5 s<br>(5.863 &pm; 0.010 %)|

*: worst is an artificial program containing almost only concurrency objects (channel, select, mutex, ...). It thus represents a worst case scenario for the runtime increase.