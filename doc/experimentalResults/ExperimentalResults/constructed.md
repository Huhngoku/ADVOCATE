# constructed

## Program

| Info                    | Value |
| ----------------------- | ----- |
| Number of go files      | 0     |
| Number of lines of code | 0     |

## Trace

| Info                                         | Value |
| -------------------------------------------- | ----- |
| Number of routines                           | 34    |
| Number of atomic variables                   | 16    |
| Number of atomic operations                  | 45    |
| Number of channels                           | 24    |
| Number of channel operations                 | 50    |
| Number of mutexes                            | 10    |
| Number of mutex operations                   | 30    |
| Number of once variables                     | 4     |
| Number of once operations                    | 7     |
| Number of selects                            | 0     |
| Number of select cases                       | 0     |
| Number of executed select channel operations | 0     |
| Number of executed select default cases      | 0     |
| Number of waitgroups                         | 2     |
| Number of waitgroup operations               | 6     |

## Runtime

| Info                                                 | Value |
| ---------------------------------------------------- | ----- |
| Runtime without modifications                        | 4.112 |
| Runtime with modified runtime                        | 4.114 |
| Runtime with modified runtime and trace creation     | 4.123 |
| Overhead of modified runtime [s]                     | 0.002 |
| Overhead of modified runtime [\%]                    | 0.049 |
| Overhead of modified runtime and trace creation [s]  | 0.011 |
| Overhead of modified runtime and trace creation [\%] | 0.268 |
| Runtime for analysis [s]                             | 0.009 |

## Found Results

==================== Summary ====================\
\
-------------------- Critical -------------------\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:131\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:123\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:145\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:140\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:182\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:188\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:183\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:198\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:315\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:303\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:362\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:352\
Possible send on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:438\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:430\
Found concurrent Recv on same channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:448\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:452\
-------------------- Warning --------------------\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:131\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:127\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:145\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:141\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:170\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:167\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:182\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:199\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:183\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:193\
Receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:341\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:335\
Possible receive on closed channel:\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:362\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/potentialBugs/constructed/potentialBugs.go:358\

=================================================
