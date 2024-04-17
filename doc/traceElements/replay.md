# Replay
This element only exists in rewritten traces, not in recorded ones. It signalizes
the start and end of the part of the trace, that was detected as a potential 
bug during the analysis and then rewritten. Is the start signal reached during replay,
a message is printed, to inform the user that the interesting part of the replay
will now run. The end signal will also print a message and then disable the 
replay, so that the program will continue to run by itself.
The disabling is necessary, because the rewriting ignores the rentability of 
the code outside of the part of the trace, that contains the possible bug.
If the potential bug / situation is passed without crashing the program
it would most likely get stuck, because the run was altered by the rewrite.

## Trace element
In the routine, where the new routine is created, the following element is added.
```
X,[tpost],[se]
```
where `X` identifies the element as an replay control element.\
- [tpost] $\in \mathbb N$: This is the time. It is replaced by the float value of the global counter at the moment when it is supposed to be run
- [se]: `s` or `e` signals wether it is an start or end marker.
