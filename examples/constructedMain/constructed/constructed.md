# n01
- True negative
- No possible send on closed channel, because of happens before relation with other channel

# n02
- True negative
- No possible send to closed channel, because of happens before 
relation with wait group

# n03
- True negative
- No possible send to closed channel, because once.Do blocks until the subroutine has finished

# n04
- False negative
- Possible send to closed channel
- Is not detected, because critical sections are not reordered

# n05 
- True positive
- Correctly detected possible send and receive on closed unbuffered channel

# n06
- True positive
- Correctly detected possible send and receive on closed buffered channel

# n07
- True negative
- No possible send to closed channel, because receive must happen before close

# n08
- True positive
- Correctly detected possible receive on closed channel without send

# n09 
- True positive
- Correctly detected possible send and receive on closed unbuffered channel in select

# n10
- False negative
- Possible send to closed channel not detected, because false select statement was executed in run

# n11
- True negative
- No send to closed channel possible, because send and close are in once

# n12
- False negative 
- Potential send to closed channel not detected because false nonce was selected