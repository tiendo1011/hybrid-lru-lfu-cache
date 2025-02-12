# Overview
a high-performance in-memory cache that uses a hybrid LRU-LFU eviction strategy (like Redis does).

# Properties
- A fixed-size in-memory cache that stores key-value pairs
- Supports Get(key string) in O(1) time complexity.
- Supports Set(key string, value int) in O(1) time complexity.
- Evicts the least valuable item when the cache is full in O(1) time complexity,
   using this hybrid eviction strategy:
   - If two keys have the same frequency, evict the least recently used one.
   - Otherwise, evict the least frequently used item.

# Example
```go
store := cache.New(2) # Create a new cache with a capacity of 2
store.Set("a", 1) # Set key "a" to value 1
store.Get("a") # Get key "a" => 1
```

# Implementation details
## Data Structures
- keyMap (a golang map):
  - this has the key as the key provided by the user
  - the value will be a pointer to a node in a doublely linked list, this
  doublely linked list is tracked by freqMap (see below)
- freqMap (a golang map):
  - this has the key as the frequency
  - the value will be a pointer to the head of a doublely linked list ,this
  doubly linked list will reflect the order of most recently used items to the
  least recently used items with the same frequency
- we will have a minFreq to track the least frequently used count

## Algorithm
- Get(key string):
  - if the key not exists in keyMap, return (0, false)
  - if the key exists in keyMap:
    - get the node from the keyMap
    - increment the frequency of the node
    - remove the node from the freqMap old frequency list
    - add the node to the freqMap with the new frequency (old frequency + 1)
    - return the value of the node
- Set(key string, value int):
  - if the key already exists:
    - update the value
    - do the same as Get(key) when the key exists
  - if the key does not exist:
    - if the cache is not full:
      - create a new node with the key, value, and frequency 1
      - add the node to the keyMap
      - add the node to the freqMap with frequency 1
    - if the cache is full:
      - check freqMap[minFreq] to get the least frequently used node
      - remove that node from the freqMap
      - remove that node from the keyMap
      - do the same as Set(key, value) when the cache is not full
