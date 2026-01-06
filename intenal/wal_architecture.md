
# WAL Architecture (Authoritative Design)

This document defines the **Write-Ahead Log (WAL)** architecture for the system.
It is correctness-first and intentionally conservative.
Performance optimizations are explicitly out of scope.

If any behavior described here is violated, the system is **incorrect by definition**.



## 1. Core Principles

### 1.1 Single Source of Truth
- WAL is the **only authoritative state**
- All memory state is **derived**
- If it cannot be rebuilt from WAL, it must not exist

### 1.2 Append-Only Invariant
- WAL is append-only
- Existing bytes are never modified
- History is immutable



## 2. File Open Semantics

WAL file must be opened as:

O_WRONLY | O_APPEND | O_CREAT

Manual seeking is forbidden.



## 3. WAL Record Format

Each record is self-delimiting and self-validating.

[length:u32][crc32:u32][payload:N bytes]

- CRC32 computed over payload only
- Fixed endianness
- No nested framing



## 4. Write Path (Commit Protocol)

### 4.1 Write Algorithm

serialize payload  
compute CRC32(payload)  
write(length, crc32, payload)  
fsync(wal_fd)  
return success  

### 4.2 Commit Boundary

- fsync() return defines commit
- If fsync did not return → record never happened
- CRC does NOT define commit



## 5. fsync Semantics

### Guarantees
- File data durable
- File size durable
- Required inode metadata durable

### Does NOT guarantee
- Directory entries
- Atomic multi-record commits



## 6. Directory fsync Rules

Directory fsync required when:
- Creating WAL
- Renaming WAL
- Deleting WAL

Correct creation sequence:

open(O_CREAT)  
fsync(wal_fd)  
fsync(dir_fd)  



## 7. Recovery / Replay Algorithm

Prefix-only recovery.

offset = 0  
while true:  
  read length  
  if EOF: break  

  read crc + payload  
  if EOF: truncate(offset); break  

  if crc mismatch: truncate(offset); break  

  apply record  
  offset += record_size  

Everything after first invalid record is garbage.



## 8. In-Memory State Rebuild

- Memory starts empty
- WAL replay rebuilds all state
- Inflight jobs are treated as FAILED(TIMEOUT)



## 9. CRC32 Role and Limits

CRC32 answers:
- Are bytes internally consistent?

CRC32 does NOT answer:
- Was fsync completed?
- Was operation acknowledged?



## 10. Formal WAL Invariants

These must hold **always**.

1. WAL is append-only
2. WAL replay is prefix-based
3. fsync return defines commit
4. CRC validates integrity only
5. No committed record is ever lost
6. No uncommitted record is ever applied
7. In-memory state is fully reconstructible
8. WAL deletion = total data loss
9. System must refuse startup if WAL missing unexpectedly



## 11. Crash Timeline Matrix

| Crash Point | Disk State | Replay Outcome |
|------------|-----------|----------------|
| Before write | No record | No-op |
| During write | Partial bytes | CRC fail → truncate |
| After write, before fsync | Maybe bytes | Treated as uncommitted |
| During fsync | Undefined | CRC check + truncate |
| After fsync | Full record | Replayed |
| During replay | WAL intact | Replay restarts safely |



## 12. WAL Test Plan (Kill -9 Matrix)

### Mandatory Tests

1. kill -9 before write
2. kill -9 mid-write
3. kill -9 after write, before fsync
4. kill -9 during fsync
5. kill -9 immediately after fsync
6. power-loss simulation (if possible)
7. disk-full during write
8. corrupted tail bytes
9. random bit flip in payload
10. random bit flip in length
11. random bit flip in crc

### Expected Results

- System always recovers to a valid prefix
- No job partially applied
- No job silently lost
- WAL truncated deterministically



## 13. Explicitly Out of Scope

- Snapshots
- WAL truncation
- Compaction
- Replication
- Performance tuning
- Async fsync batching



## 14. Final Statement

This WAL design is:

- Crash-consistent
- Kill -9 honest
- Prefix-correct

This is the strongest guarantee achievable in user space.
