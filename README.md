## Batch processing utilising a 'worker' pool.

Using Go builtin sync, channels and goroutines to prevent race conditions.

### Spec

- Process accepted Jobs in batches using a BatchProcessor
- Provide a way to configure the batching behaviour i.e. size and frequency
- Expose a shutdown method which returns after all previously accepted Jobs are processed
