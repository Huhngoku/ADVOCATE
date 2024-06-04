# Results

| Number | Explanation | Result | Replay (exit codes) | Type | OK |
| --- | --- | --- | --- | --- | --- |
| 1 | No send on closed because of other channel |  No found | No | TN | Y |
| 2 | No send on closed because of wait group    |  No found | No | TN | Y |
| 3 | No send on closed because of once          |  No found | No | TN | Y |
| 4 | Send on close not detected because of critical sections | No found | No | FN | Y |
| 5 | Possible Send and Recv on closed           | Found both | 30, 31 | TP | Y |
| 6 | Possible Send on closed                    | No found  | No | TP | N |
| 7 | No possible send/recv on closed                 | No found  | No | TN | Y |
| 8 | Actual recv on closed                     | Found | No | TP | Y |
| 9 | Send/recv on closed in select             | Found | 30, 30, 31, 31 | TP | Y |
| 10 | Send on closed in select with default | No Found | No | ? | ? |
| 11 | Send on closed in select with case    | No Found | No | TP | N |
| 12 | No send on close because of once | No Found | No | TN | Y |
| 13 | Send on close no detected because of once | No Found | No | FN | Y |
| 14 | Send on close detected in spite of once | Found | 30 | TP | Y |
| 15 | Send on close not detected because of tryLock | No Found | No | FN | Y |
| 16 | Send on close detected in spite of tryLock | Found | 30 | TP | Y |
| 17 | Send on close not detected because of tryLock | No Found | No | FN | Y |
| 18 | Send on close with wait group | No found | No | FN | N |
| 19 | Send on close not detected because of tryLock | No Found | No | FN | Y |
| 20 | Send on close in function | Found | 30 | TP | Y |
| 21 | Concurrent recv on channel | No found | No | FN | N |
| 22 | No concurrent recv on channel | No found | No | TN | Y |
| 23 | No concurrent recv on buffered channel | No found | No | TN | Y |
| 24 | No concurrent recv on channel | No found | No | FN | Y |
| 25 | No possible negative wait group counter | No found | No | TN | Y |
| 26 | No possible negative wait group counter | No found | No | TN | Y |
| 27 | Possible negative wait group counter | Found | 32 | TP | Y |
| 39 | Leak on unbuffered channel without possible partner | Found | No | TP | Y |
| 40 | Leak on unbuffered channel without possible partner | Found | No | TP | Y |
| 41 | No leak on unbuffered channel without possible partner | No Found | No | TN | Y |
| 42 | Leak on buffered channel without possible partner | Found | 20 | TP | Y |
| 43 | Leak on buffered channel without possible partner | Found | 20 | TP | Y |
| 44 | Leak on wait group | Found | No | TP | Y |
| 45 | Leak on mutex | Found | No | TP | Y |
| 46 | Leak on select | Found | No | TP | N -> finds leak on channel |
| 47 | Leak on unuffered channel with select as possible partner | Found | 20 | TP | Y |
| 48 | Select case without partner | Found | No | TP | Y |
| 49 | Leak on unbuffered channel without partner | Found | No | TP | Y |
| 50 | Leak on buffered channel with partner | Found | FAIL (12) | TP | N -> rewrite does not work |
| 51 | Leak on buffered channel with partner | Found | 20 | TP | Y |
| 52 | Leak on unbuffered channel without partner | Found | No | TP | Y |
| 53 | Leak on select with partner | Found | 20 | TP | N -> found channel with partner |
| 54 | Leak on select outwith partner | Found | 20 | TP | N -> found channel without partner |
| 55 | Leak on wait group | Found | No | TP | Y |
| 56 | Leak on conditional variable | Found | No | TP | Y |