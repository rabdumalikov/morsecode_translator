

### Performance 
- naive implementation, doesn't utilize potential duplicate knowledge, thus doing more work
    - Step 1 (Text -> Morse) Time: 1m13,251s
    - Step 2 (Morse -> Text) Time: 1m0,784s
    - Step 3 (Diff Verify)   Time: 0m10,277s
    - Step 4 (Comparator)    Time: 0m42,723s


### word level
[2025-04-03 14:56:55] -------------------- Summary --------------------
[2025-04-03 14:56:55] Step 1 (Text -> Morse) Time: 1m52,692s
[2025-04-03 14:56:55] Step 2 (Morse -> Text) Time: 1m28,883s
[2025-04-03 14:56:55] Step 3 (Diff Verify)   Time: 0m14,056s
[2025-04-03 14:56:55] Step 4 (Comparator)    Time: 0m49,238s
[2025-04-03 14:56:55] -------------------------------------------------


# last run with chunk reading
[2025-04-05 09:06:47] -------------------- Summary --------------------
[2025-04-05 09:06:47] Step 1 (Text -> Morse) Time: 0m20,362s
[2025-04-05 09:06:47] Step 2 (Morse -> Text) Time: 0m30,398s
[2025-04-05 09:06:47] Step 3 (Diff Verify)   Time: 0m6,248s
[2025-04-05 09:06:47] Step 4 (Comparator)    Time: 2m42,303s
