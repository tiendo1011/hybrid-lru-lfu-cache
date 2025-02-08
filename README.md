# GOAL: Design a high-performance in-memory cache that uses a hybrid LRU-LFU eviction strategy (like Redis does).

## Problem Statement
Build a fixed-size in-memory cache that:
1. Stores key-value pairs.
2. Supports Get(key) in O(1) time complexity.
3. Supports Put(key, value) in O(1) time complexity.
4. Evicts the least valuable item when the cache is full in O(1) time complexity,
   using this hybrid eviction strategy:
   - If two keys have the same frequency, evict the least recently used one.
   - Otherwise, evict the least frequently used item.

# Attacking the Problem
LRU-LFU cache is a cache where we can store a maximum of N elements. If the cache
is full and we want to add a new element, we must remove the least frequently
used element, if there are multiple items with the same frequency, remove the
least recently used element.
The cache should be able to get and set elements in O(1) time.

## Data Structures
- Have a keyMap to store the key-value pairs
- the value will be a pointer to a node in a doublely linked list, with the
following fields:
  - key (used to remove the key-pair from keyMap after evict it from freqMap)
  - value
  - freq
  - prev
  - next
  the prev, next is used to build a doublly linked list, this doubly linked list
  will reflect the order of most recently used items to the least recently used
  items with the same frequency, which is tracked by freqMap
- freqMap will have:
  - key as the frequency
  - value as a DLL struct, with the linked list of items with the same
  frequency, most recently used to least recently used, with pointer to the
  least recently used (so that we can evict it in O(1))
- we will have a minFreq to track the least frequently used count

## Algorithm
- get(key): gets the value at key. If no such key exists, return (0, false).
  - if the key not exists in keyMap, return (0, false)
  - if the key exists in keyMap:
    - get the node from the keyMap
    - increment the frequency of the node
    - remove the node from the freqMap old frequency list
    - add the node to the freqMap with the new frequency (old frequency + 1)
    - return the value of the node
- set(key, value): sets key to value. If there are already N items in the cache,
  and we are adding a new item, then it should also remove a item according to
  the eviction rule (lfu then lru).
  - if the key already exists:
    - update the value
    - do the same as get(key) when the key exists
  - if the key does not exist:
    - if the cache is not full:
      - create a new node with the key, value, and frequency 1
      - add the node to the keyMap
      - add the node to the freqMap with frequency 1
    - if the cache is full:
      - check freqMap[minFreq] to get the least frequently used node
      - remove the least frequently used node from the freqMap
      - remove the least frequently used node from the keyMap
      - do the same as set(key, value) when the cache is not full

- Should I use a dummy node?
  - Yes, because it removes the edge cases. Having a front and tail dummy nodes
    ensure that all nodes have prev/next node
