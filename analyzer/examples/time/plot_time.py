import matplotlib.pyplot as plt

names = [
    "bbolt, 8.2 MB",
    "constructed 5.3 kB",
    "gocrawl 213 kB",
    "htcat 1.5 MB",
    "pgzip 7.7 MB",
    "readme 626 B",
    "serving3068 4.4 kB",
    "serving5865 1.1 kB",
    "sorty 197.8 MB"
]

size = [
    8200, 5.3, 213, 1500, 7700, 0.626, 4.4, 1.1, 197800
]

time = [
    1.0883, 0.0072, 0.0406, 0.3552, 146.3901, 0.0014, 0.0064, 0.0212, 177.3169
]

fig, ax = plt.subplots()
ax.plot(size, time, 'ro')

# for i, name in enumerate(names):
#     ax.annotate(name, (size[i] + 1000, time[i] + 5))

plt.show()